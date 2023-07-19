package config

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
	"sync"
	"time"
)

var (
	config *ConsulConfig
	once   sync.Once
)

type ConsulConfig struct {
	Addr           string
	WatchFrequency time.Duration

	lock    sync.RWMutex
	config  map[string]interface{}
	watcher map[string]func()
}

func InitConsulConfig(addr, path, consulType string, watchFrequency time.Duration) {
	err := viper.AddRemoteProvider("consul", addr, path)
	if err != nil {
		logrus.Fatalf("viper.AddRemoteProvider error: %v", err)
	}

	viper.SetConfigType(consulType)
	err = viper.ReadRemoteConfig()
	if err != nil {
		logrus.Fatalf("viper.ReadRemoteConfig error: %v", err)
	}

	once.Do(func() {
		config = &ConsulConfig{
			Addr:           addr,
			WatchFrequency: watchFrequency,

			lock:    sync.RWMutex{},
			config:  make(map[string]interface{}),
			watcher: make(map[string]func()),
		}

		go func() {
			for range time.NewTicker(watchFrequency).C {
				err := viper.WatchRemoteConfig()
				if err != nil {
					logrus.Errorf("watch remote config err: %+v", err)
				}
				doWatch()
			}
		}()
	})
}

func Get[T any](key string) T {
	return viper.Get(key).(T)
}

func GetOrDefault[T any](key string, defaultValue T) T {
	viper.SetDefault(key, defaultValue)
	return viper.Get(key).(T)
}

func Watch(key string, callback func()) {
	config.lock.Lock()
	defer config.lock.Unlock()
	config.config[key] = viper.Get(key)
	config.watcher[key] = callback
}

func UnWatch(key string) {
	config.lock.Lock()
	defer config.lock.Unlock()
	delete(config.watcher, key)
	delete(config.config, key)
}

func doWatch() {
	config.lock.Lock()
	defer config.lock.Unlock()
	for key, callback := range config.watcher {
		ov := config.config[key]
		nv := viper.Get(key)
		if ov != nv {
			config.config[key] = nv
			callback()
		}
	}
}

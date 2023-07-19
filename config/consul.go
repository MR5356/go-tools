package config

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
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

func NewConsulConfig(addr string, watchFrequency time.Duration) *ConsulConfig {
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
				config.doWatch()
			}
		}()
	})

	return config
}

func (c *ConsulConfig) Get(key string) interface{} {
	return viper.Get(key)
}

func (c *ConsulConfig) GetOrDefault(key string, defaultValue any) any {
	viper.SetDefault(key, defaultValue)
	return viper.Get(key)
}

func (c *ConsulConfig) Watch(key string, callback func()) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.config[key] = viper.Get(key)
	c.watcher[key] = callback
}

func (c *ConsulConfig) UnWatch(key string) {
	c.lock.Lock()
	defer c.lock.Unlock()
	delete(c.watcher, key)
	delete(c.config, key)
}

func (c *ConsulConfig) doWatch() {
	c.lock.Lock()
	defer c.lock.Unlock()
	for key, callback := range c.watcher {
		ov := c.config[key]
		nv := viper.Get(key)
		if ov != nv {
			c.config[key] = nv
			callback()
		}
	}
}

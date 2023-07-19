package config

type Config interface {
	Get(key string) interface{}
	GetOrDefault(key string, defaultValue any) any
	Watch(key string, callback func())
	UnWatch(key string)
}

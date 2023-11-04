package config

import "sync"

var (
	mp sync.Map
)

func setConfig(config *Config) {
	mp.Store("config", config)
}

func DefaultConfig() *Config {
	v, ok := mp.Load("config")
	if !ok {
		return nil
	}
	return v.(*Config)
}

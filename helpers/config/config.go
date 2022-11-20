package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

// Config
type Config struct {
	Admin struct {
		MaxLiveInstances int `yaml:"max_live_instances"`
		InstanceInterval int `yaml:"instance_interval"`
		MaxTotalBots     int `yaml:"max_total_bots"`
	} `yaml:"admin"`
	Bot struct {
		Address struct {
			Router  string `yaml:"router"`
			Factory string `yaml:"factory"`
			Weth    string `yaml:"weth"`
			Token   string `yaml:"token"`
		} `yaml:"address"`
		Config struct {
			ChainId      int `yaml:"chainid"`
			SwapInterval int `yaml:"swap_interval"`
			MaxSwaps     int `yaml:"max_swaps"`
		}
	}
}

// ParseConfig
func ParseConfig(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var conf Config
	if err := yaml.Unmarshal(data, &conf); err != nil {
		return nil, err
	}

	return &conf, nil
}

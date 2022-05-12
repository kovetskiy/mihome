package main

import (
	"github.com/kovetskiy/ko"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Listen string `yaml:"listen" required:"true" env:"LISTEN" default:":5000"`
	Mi     struct {
		Extractor string `yaml:"extractor" required:"true" env:"EXTRACTOR" default:""`
		Username  string `yaml:"username" required:"true" env:"MI_USERNAME" default:""`
		Password  string `yaml:"password" required:"true" env:"MI_PASSWORD" default:""`
		Server    string `yaml:"server"   required:"true" env:"MI_SERVER"   default:""`
	} `yaml:"mi"     required:"true" env:""       default:""`
}

func LoadConfig(path string) (*Config, error) {
	config := &Config{}
	err := ko.Load(path, config, yaml.Unmarshal)
	if err != nil {
		return nil, err
	}

	return config, nil
}

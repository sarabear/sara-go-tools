package main

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	AccessKeyId     string   `yaml:"accessKeyId"`
	AccessKeySecret string   `yaml:"accessKeySecret"`
	RegionId        string   `yaml:"regionId"`
	Domains         []Domain `yaml:"domains"`
}

type Domain struct {
	Suffix string `yaml:"suffix"`
	Domain string `yaml:"domain"`
	RR     string `yaml:"rr"`
}

func loadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

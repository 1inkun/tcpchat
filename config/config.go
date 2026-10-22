package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server ServerConfig `yaml:"Server"`
}

type ServerConfig struct {
	PORT       string `yaml:"Port"`
	SAFEMODE   bool   `yaml:"Safemode"`
	TIMEWINDOW int    `yaml:"TimeWindow"`
	MSGLIMIT   int    `yaml:"MsgLimit"`
}

func InitConfig() *Config {
	var cfg Config
	data, err := os.ReadFile("config.yaml")
	if err != nil {
		// 读取配置出错,使用默认配置
		log.Printf("读取配置文件出错,使用默认配置:%s", err.Error())
		cfg := defaultSetting(cfg)
		return &cfg
	}
	if err = yaml.Unmarshal(data, &cfg); err != nil {
		// 读取配置出错,使用默认配置
		log.Printf("读取配置文件出错,使用默认配置:%s", err.Error())
		cfg := defaultSetting(cfg)
		return &cfg
	}
	return &cfg
}

func defaultSetting(cfg Config) Config {
	cfg.Server.PORT = ":8080"
	cfg.Server.SAFEMODE = true
	cfg.Server.TIMEWINDOW = 2
	cfg.Server.MSGLIMIT = 8
	return cfg
}

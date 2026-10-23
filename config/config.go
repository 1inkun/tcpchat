package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server ServerConfig `yaml:"Server"`
	Redis  RedisConfig  `yaml:"Redis"`
}

type ServerConfig struct {
	PORT       string `yaml:"Port"`
	SAFEMODE   bool   `yaml:"Safemode"`
	TIMEWINDOW int    `yaml:"TimeWindow"`
	MSGLIMIT   int    `yaml:"MsgLimit"`
	TIMEOUT    int    `yaml:"Timeout"`
	BANTIME    int    `yaml:"BanTime"`
}

type RedisConfig struct {
	Addr     string `yaml:"Addr"`
	Password string `yaml:"Password"`
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
	cfg.Server.TIMEOUT = 2
	cfg.Redis.Addr = "localhost:6379"
	cfg.Redis.Password = "a123456"
	cfg.Server.BANTIME = 3
	return cfg
}

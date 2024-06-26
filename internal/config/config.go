package config

import (
	"flag"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Database DatabaseConfig `yaml:"database"`
}

type DatabaseConfig struct {
	Host         string `yaml:"host"`
	MasterPort   int    `yaml:"master_port"`
	SyncPort     int    `yaml:"sync_port"`
	AsyncPort    int    `yaml:"async_port"`
	User         string `yaml:"user"`
	Password     string `yaml:"password"`
	DatabaseName string `yaml:"dbname"`
}

var configPath string

func init() {
	flag.StringVar(&configPath, "config", "", "Specifies the path to the config file.")
	flag.Parse()
}

func LoadConfig() *Config {
	if configPath == "" {
		panic("path to a config file not specified")
	}

	var config Config
	err := cleanenv.ReadConfig(configPath, &config)

	if err != nil {
		panic(err)
	}

	return &config
}

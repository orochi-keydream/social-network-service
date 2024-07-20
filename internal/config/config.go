package config

import (
	"flag"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	KafkaBrokers    []string              `yaml:"kafka_brokers"`
	Producers       ProducerConfigs       `yaml:"producers"`
	Consumers       ConsumerConfigs       `yaml:"consumers"`
	Database        DatabaseConfig        `yaml:"database"`
	ShardedDatabase ShardedDatabaseConfig `yaml:"sharded_database"`
	Redis           RedisConfig           `yaml:"redis"`
}

type ProducerConfigs struct {
	Posts ProducerConfig `yaml:"posts"`
	Feed  ProducerConfig `yaml:"feed"`
}

type ProducerConfig struct {
	Topic string `yaml:"topic"`
}

type ConsumerConfigs struct {
	Posts ConsumerConfig `yaml:"posts"`
	Feed  ConsumerConfig `yaml:"feed"`
}

type ConsumerConfig struct {
	Topic string `yaml:"topic"`
}

type DatabaseConfig struct {
	Host         string `yaml:"host"`
	MasterPort   int    `yaml:"master_port"`
	SyncPort     int    `yaml:"sync_port"`
	AsyncPort    int    `yaml:"async_port"`
	DatabaseName string `yaml:"dbname"`
	User         string `yaml:"user"`
	Password     string `yaml:"password"`
}

type ShardedDatabaseConfig struct {
	Host         string `yaml:"host"`
	Port         int    `yaml:"port"`
	DatabaseName string `yaml:"dbname"`
	User         string `yaml:"user"`
	Password     string `yaml:"password"`
}

type RedisConfig struct {
	ConnectionString string `yaml:"connection_string"`
}

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config", "", "Specifies the path to the config file.")
	flag.Parse()
}

func LoadConfig() Config {
	if configPath == "" {
		panic("path to a config file not specified")
	}

	var config Config
	err := cleanenv.ReadConfig(configPath, &config)

	if err != nil {
		panic(err)
	}

	return config
}

package config

import (
	"flag"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	KafkaBrokers []string        `yaml:"kafka_brokers"`
	Producers    ProducerConfigs `yaml:"producers"`
	Consumers    ConsumerConfigs `yaml:"consumers"`
	Database     DatabaseConfig  `yaml:"database"`
	Redis        RedisConfig     `yaml:"redis"`
	Tarantool    TarantoolConfig `yaml:"tarantool"`

	UseTarantool bool
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
	User         string `yaml:"user"`
	Password     string `yaml:"password"`
	DatabaseName string `yaml:"dbname"`
}

type RedisConfig struct {
	ConnectionString string `yaml:"connection_string"`
}

type TarantoolConfig struct {
	ConnectionString string `yaml:"connection_string"`
}

var (
	configPath   string
	useTarantool bool
)

func init() {
	flag.StringVar(&configPath, "config", "", "Specifies the path to the config file.")
	flag.BoolVar(&useTarantool, "use-tarantool", false, "Specifies if Tarantool DB should be used.")
	flag.Parse()
}

func LoadConfig() Config {
	if configPath == "" {
		panic("path to a config file not specified")
	}

	var config Config
	err := cleanenv.ReadConfig(configPath, &config)

	config.UseTarantool = useTarantool

	if err != nil {
		panic(err)
	}

	return config
}

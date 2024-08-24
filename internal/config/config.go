package config

import (
	"flag"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Service     ServiceConfig  `yaml:"service"`
	Kafka       KafkaConfig    `yaml:"kafka"`
	Database    DatabaseConfig `yaml:"database"`
	Redis       RedisConfig    `yaml:"redis"`
	GrpcClients GrpcClients    `yaml:"grpc_clients"`
}

type ServiceConfig struct {
	HttpPort      int `yaml:"http_port"`
	MetricsPort   int `yaml:"metrics_port"`
	WebSocketPort int `yaml:"websocket_port"`
}

type KafkaConfig struct {
	Brokers   []string        `yaml:"brokers"`
	Producers ProducerConfigs `yaml:"producers"`
	Consumers ConsumerConfigs `yaml:"consumers"`
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

type GrpcClients struct {
	DialogueServiceAddr string `yaml:"dialogue_service_addr"`
	CounterServiceAddr  string `yaml:"counter_service_addr"`
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

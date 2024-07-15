package config

import (
	"os"
)

type Config struct {
	PostgresUser     string
	PostgresPassword string
	PostgresDB       string
	PostgresHost     string
	PostgresPort     string

	KafkaBrokerHost string
	KafkaBrokerPort string
	KafkaInTopic    string
	KafkaOutTopic   string
}

func LoadConfig() (*Config, error) {
	config := &Config{
		PostgresUser:     getEnv("POSTGRES_USER", ""),
		PostgresPassword: getEnv("POSTGRES_PASSWORD", ""),
		PostgresDB:       getEnv("POSTGRES_DB", ""),
		PostgresHost:     getEnv("POSTGRES_HOST", ""),
		PostgresPort:     getEnv("POSTGRES_PORT", ""),

		KafkaBrokerHost: getEnv("KAFKA_BROKER_HOST", ""),
		KafkaBrokerPort: getEnv("KAFKA_BROKER_PORT", ""),
		KafkaInTopic:    getEnv("KAFKA_IN_TOPIC", ""),
		KafkaOutTopic:   getEnv("KAFKA_OUT_TOPIC", ""),
	}
	return config, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

package config

import (
	"os"
	"strconv"
)

type GlobalConfig struct {
	App *AppConfig
	DB  *DBConfig
}

func LoadGlobalConfig() *GlobalConfig {
	return &GlobalConfig{
		App: LoadAppConfig(),
		DB:  LoadDBConfig(),
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if envValue := os.Getenv(key); envValue != "" {
		return envValue
	}

	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if envValue := os.Getenv(key); envValue != "" {
		if intEnvValue, err := strconv.Atoi(envValue); err == nil {
			return intEnvValue
		}
	}

	return defaultValue
}

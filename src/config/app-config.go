package config

type AppConfig struct {
	Port int
}

func LoadAppConfig() *AppConfig {
	config := &AppConfig{}

	config.Port = getEnvAsInt("PORT", 8080)

	return config
}

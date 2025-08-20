package config

type AppConfig struct {
	Port int
	// CookieSecretKey string
}

func LoadAppConfig() *AppConfig {
	config := &AppConfig{}

	config.Port = getEnvAsInt("PORT", 8080)
	// config.CookieSecretKey = getEnvOrDefault("COOKIE_SECRET_KEY", "secret-cookie")

	return config
}

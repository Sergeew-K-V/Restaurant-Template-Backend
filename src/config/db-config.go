package config

type DBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

func LoadDBConfig() *DBConfig {
	config := &DBConfig{}

	config.DBName = getEnvOrDefault("DB_NAME", "go-postgres")
	config.Host = getEnvOrDefault("DB_HOST", "localhost")
	config.Password = getEnvOrDefault("DB_PASS", "admin")
	config.Port = getEnvAsInt("DB_PORT", 8888)
	config.User = getEnvOrDefault("DB_USER", "admin")
	config.SSLMode = getEnvOrDefault("DB_SSLMODE", "disabled")

	return config
}

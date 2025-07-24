package config

func GetDevelopmentConfig() *Config {
	cfg := &Config{}

	cfg.Database.Host = "localhost"
	cfg.Database.Port = "5421"
	cfg.Database.User = "admin"
	cfg.Database.Password = "admin"
	cfg.Database.DBName = "subscri-docker"

	return cfg
}

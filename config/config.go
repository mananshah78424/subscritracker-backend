package config

import (
	"os"
	"strings"
)

type Config struct {
	Database   DatabaseConfig
	GoogleAuth GoogleAuthConfig
	Frontend   FrontendConfig
}

type FrontendConfig struct {
	URL string
}

type GoogleAuthConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  bool
}

func GetConfig() *Config {
	env := strings.ToLower(os.Getenv("NODE_ENV"))
	// If the environment is not set, default to development
	if env == "" {
		env = "development"
	}

	switch env {
	case "development":
		return GetDevelopmentConfig()
	default:

		cfg := &Config{}

		cfg.Database.Host = os.Getenv("DB_HOST")
		cfg.Database.Port = os.Getenv("DB_PORT")
		cfg.Database.User = os.Getenv("DB_USER")
		cfg.Database.Password = os.Getenv("DB_PASSWORD")
		cfg.Database.DBName = os.Getenv("DB_NAME")
		cfg.Database.SSLMode = os.Getenv("DB_SSL_MODE") == "true"

		// Frontend configuration
		cfg.Frontend.URL = os.Getenv("FRONTEND_URL")
		if cfg.Frontend.URL == "" {
			cfg.Frontend.URL = "http://localhost:3000" // Default fallback
		}

		return cfg
	}
}

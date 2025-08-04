package config

import (
	"os"
)

func GetDevelopmentConfig() *Config {
	cfg := &Config{}

	cfg.Database.Host = "localhost"
	cfg.Database.Port = "5421"
	cfg.Database.User = "admin"
	cfg.Database.Password = "admin"
	cfg.Database.DBName = "subscri-docker"

	cfg.GoogleAuth.ClientID = os.Getenv("GOOGLE_CLIENT_ID")
	cfg.GoogleAuth.ClientSecret = os.Getenv("GOOGLE_CLIENT_SECRET")
	cfg.GoogleAuth.RedirectURL = "http://localhost:8080/auth/google/callback"

	// Frontend configuration
	cfg.Frontend.URL = os.Getenv("FRONTEND_URL")
	if cfg.Frontend.URL == "" {
		cfg.Frontend.URL = "http://localhost:3000" // Default for development
	}

	return cfg
}

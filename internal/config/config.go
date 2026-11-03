// Package config loads runtime configuration from the environment.
package config

import (
	"errors"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port          string
	DatabaseURI   string
	AdminEmail    string
	AdminPassword string
}

// Load reads configuration from environment variables, falling back to a
// .env file in the working directory when one is present.
func Load() (Config, error) {
	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		log.Printf("load .env: %v", err)
	}

	cfg := Config{
		Port:          getenv("PORT", "3000"),
		DatabaseURI:   os.Getenv("DB_URI"),
		AdminEmail:    getenv("ADMIN_EMAIL", "root"),
		AdminPassword: getenv("ADMIN_PASSWORD", "root"),
	}

	if cfg.DatabaseURI == "" {
		return Config{}, errors.New("DB_URI is not set")
	}

	return cfg, nil
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	log.Printf("%s is not set, using %q", key, fallback)
	return fallback
}

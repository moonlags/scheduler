package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	port  string
	dbUri string
}

func newConfig() (Config, error) {
	if err := godotenv.Load(); err != nil {
		return Config{}, err
	}

	cfg := Config{
		port:  os.Getenv("PORT"),
		dbUri: os.Getenv("DB_URI"),
	}

	if cfg.port == "" {
		log.Printf("PORT variable is not set, setting to 3000\n")
		cfg.port = "3000"
	}
	if cfg.dbUri == "" {
		return Config{}, fmt.Errorf("DB_URI variable is not set")
	}

	return cfg, nil
}

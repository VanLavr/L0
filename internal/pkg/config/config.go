package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Environment string
	NatsUrl     string
	ClusterID   string
	SubName     string
	SubjName    string
}

func New() *Config {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("cannot load env variables due to:", err)
	}

	return &Config{
		Environment: os.Getenv("ENV"),
		NatsUrl:     os.Getenv("NATS_URL"),
		ClusterID:   os.Getenv("CLUSTER"),
		SubName:     os.Getenv("SUB_NAME"),
		SubjName:    os.Getenv("SUBJ"),
	}
}

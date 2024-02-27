package config

import (
	"log"
	"log/slog"
	"os"
	"strconv"

	er "github.com/VanLavr/L0/internal/pkg/err"
	"github.com/joho/godotenv"
)

type Config struct {
	Environment        string
	NatsUrl            string
	ClusterID          string
	SubName            string
	SubjName           string
	PostgresConnection string
	Eviction           int // seconds
	Ttl                int // seconds
}

func New() *Config {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("cannot load env variables due to:", err)
	}

	ev := os.Getenv("CACHE_EVICTION")
	eviction, err := strconv.Atoi(ev)
	if err != nil {
		slog.Error(er.ErrInvalidEnvironment.Error())
		os.Exit(1)
	}

	ttl := os.Getenv("TTL")
	t, err := strconv.Atoi(ttl)
	if err != nil {
		slog.Error(er.ErrInvalidEnvironment.Error())
		os.Exit(1)
	}

	return &Config{
		Environment:        os.Getenv("ENV"),
		NatsUrl:            os.Getenv("NATS_URL"),
		ClusterID:          os.Getenv("CLUSTER"),
		SubName:            os.Getenv("SUB_NAME"),
		SubjName:           os.Getenv("SUBJ"),
		PostgresConnection: os.Getenv("POSTGRES"),
		Eviction:           eviction,
		Ttl:                t,
	}
}

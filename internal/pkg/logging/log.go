package logging

import (
	"log"
	"log/slog"
	"os"

	"github.com/VanLavr/L0/internal/pkg/config"
	"github.com/VanLavr/L0/internal/pkg/err"
)

type Logger struct {
	l *slog.Logger
}

func New(cfg *config.Config) *Logger {
	if cfg.Environment == "dev" {
		log := new(Logger)
		log.l = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			AddSource: true,
			Level:     slog.LevelDebug,
		}))

		return log
	} else if cfg.Environment == "demo" {
		log := new(Logger)
		log.l = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			AddSource: false,
			Level:     slog.LevelInfo,
		}))

		return log
	} else {
		log.Fatal(err.ErrInvalidEnvironment)
		return nil
	}
}

func (l *Logger) SetAsDefault() {
	slog.SetDefault(l.l)
}

package repo

import (
	"context"
	"log/slog"
	"os"

	"github.com/VanLavr/L0/internal/pkg/config"
	"github.com/VanLavr/L0/internal/service"
	"github.com/VanLavr/L0/model"
	"github.com/jackc/pgx/v4"
)

type postgres struct {
	db   *pgx.Conn
	conn string
}

func NewPostgres(cfg config.Config) service.Repository {
	return &postgres{conn: cfg.PostgresConnection}
}

func (p *postgres) GenerateTrackNumber() string {
	row, err := p.db.Query(context.Background(), "select * from generate_unique_id()")
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	defer row.Close()

	var trackNumber string
	row.Next()
	if err := row.Scan(&trackNumber); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	return trackNumber
}

func (p *postgres) SaveOrder(*model.Order) string {
	panic("not implemented")
}

func (p *postgres) GetOrder(string) (*model.Order, error) {
	panic("not implemented")
}

func (p *postgres) Connect(string) error {
	db, err := pgx.Connect(context.Background(), p.conn)
	if err != nil {
		return err
	}

	p.db = db
	return nil
}

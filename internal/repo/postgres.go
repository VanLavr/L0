package repo

import (
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
	panic("not implemented")
}

func (p *postgres) SaveOrder(*model.Order) string {
	panic("not implemented")
}

func (p *postgres) GetOrder(string) (*model.Order, error) {
	panic("not implemented")
}

func (p *postgres) Connect(string) error {
	panic("not implemented")
}

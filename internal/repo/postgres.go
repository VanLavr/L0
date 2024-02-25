package repo

import (
	"context"
	"fmt"
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

func (p *postgres) SaveOrder(order *model.Order) {
	delID, err := p.saveDelivery(&order.D)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	payID, err := p.savePayment(&order.P)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}

func (p *postgres) savePayment(pay *model.Payment) (int, error) {
	var payID int
	if err := p.db.QueryRow(context.Background(), fmt.Sprintf(`
	insert into payment (t_action, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee) 
	values (%s, %s, %s, %s, %f, %d, %s, %f, %f, %f) 
	returning delivery_id`, pay.Transaction, pay.Request_id, pay.Currency, pay.Provider, pay.Amount, pay.Payment_dt, pay.Bank, pay.Delivery_cost, pay.Goods_total, pay.Custom_fee)).Scan(&payID); err != nil {
		return 0, err
	}

	return payID, nil
}

func (p *postgres) saveDelivery(del *model.Delivery) (int, error) {
	var delID int
	if err := p.db.QueryRow(context.Background(), fmt.Sprintf(`
	insert into delivery (name, phone, zip, city, address, region, email) 
	values (%s, %s, %s, %s, %s, %s) 
	returning delivery_id`, del.Name, del.Phone, del.Zip, del.City, del.Address, del.Email)).Scan(&delID); err != nil {
		return 0, err
	}

	return delID, nil
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

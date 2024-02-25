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

func NewPostgres(cfg *config.Config) service.Repository {
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
	delID, err := p.saveDelivery(order.D)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	if err := p.savePayment(order.P); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	p.saveItems(order.Items...)

	if err = p.saveOrder(order, delID); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	p.saveOrderItems(order.Items, order.Order_uid)
}

func (p *postgres) saveOrderItems(items []model.Item, orderUID string) {
	for _, item := range items {
		if err := p.saveOrderItem(item, orderUID); err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}
	}
}

func (p *postgres) saveOrderItem(item model.Item, orderUID string) error {
	if _, err := p.db.Exec(context.Background(), fmt.Sprintf("insert into items_to_orders (order_uid, chrt_id) values (%s, %d)", orderUID, item.Chrt_id)); err != nil {
		return err
	}

	return nil
}

func (p *postgres) saveOrder(order *model.Order, deliveryID int) error {
	if _, err := p.db.Exec(context.Background(), fmt.Sprintf(`
	insert into orders (order_uid, track_number, entr, delivery_id, t_action, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard)
	values (%s, %s, %s, %d, %s, %s, %s, %s, %s, %s, %d, %s, %s)`,
		order.Order_uid, order.Track_number, order.Entry, deliveryID, order.P.Transaction, order.Locale, order.Internal_signature, order.Customer_id, order.Delivery_service, order.Shardkey, order.Sm_id, order.Date_created, order.Oof_shard)); err != nil {
		return err
	}

	return nil
}

func (p *postgres) savePayment(pay model.Payment) error {
	if _, err := p.db.Exec(context.Background(), fmt.Sprintf(`
	insert into payment (t_action, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee) 
	values (%s, %s, %s, %s, %f, %d, %s, %f, %f, %f)`, pay.Transaction, pay.Request_id, pay.Currency, pay.Provider, pay.Amount, pay.Payment_dt, pay.Bank, pay.Delivery_cost, pay.Goods_total, pay.Custom_fee)); err != nil {
		return err
	}

	return nil
}

func (p *postgres) saveDelivery(del model.Delivery) (int, error) {
	var delID int
	if err := p.db.QueryRow(context.Background(), fmt.Sprintf("insert into delivery (name, phone, zip, city, address, region, email) values (\"%s\", \"%s\", \"%s\", \"%s\", \"%s\", \"%s\") returning delivery_id", del.Name, del.Phone, del.Zip, del.City, del.Address, del.Email)).Scan(&delID); err != nil {
		return 0, err
	}

	return delID, nil
}

func (p *postgres) saveItems(items ...model.Item) {
	for _, item := range items {
		if err := p.saveItem(item); err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}
	}
}

func (p *postgres) saveItem(item model.Item) error {
	if _, err := p.db.Exec(context.Background(), fmt.Sprintf(`
	insert into items (track_number, price, rid, name, sale, size, total_price, nm_id, brand, status)
	values (%s, %f, %s, %s, %f, %s, %f, %d, %s, %d)`, item.Track_number, item.Price, item.Rid, item.Name, item.Sale, item.Size, item.Total_Price, item.Nm_id, item.Brand, item.Status)); err != nil {
		return err
	}

	return nil
}

func (p *postgres) GetOrder(string) (*model.Order, error) {
	panic("not implemented")
}

func (p *postgres) Connect() error {
	db, err := pgx.Connect(context.Background(), p.conn)
	if err != nil {
		return err
	}

	p.db = db
	return nil
}

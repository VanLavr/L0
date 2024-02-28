package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"time"

	"github.com/VanLavr/L0/internal/delivery/nats"
	"github.com/VanLavr/L0/internal/pkg/config"
	"github.com/VanLavr/L0/internal/pkg/logging"
	"github.com/VanLavr/L0/internal/repo"
	"github.com/VanLavr/L0/internal/service"
	"github.com/VanLavr/L0/model"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	cfg := config.New()
	db := repo.NewPostgres(cfg)
	db.Connect()

	c := repo.NewCache(time.Duration(cfg.Ttl), time.Duration(cfg.Eviction))

	sv := service.New(db, c, cfg)
	h := nats.New(sv)

	logger := logging.New(cfg)
	logger.SetAsDefault()

	h.Connect(cfg)
	h.Subscribe(cfg)
	slog.Info("listening channel in nats...")

	order, err := sv.GetOrder("46588f63467744299add4d14fdf27f96")
	if err != nil {
		slog.Debug(err.Error())
		return
	}
	OrderPrinter(*order)

	<-ctx.Done()
	slog.Info("shutting down")
	h.Unsubscribe()
	h.CloseConnection()
}

func OrderPrinter(o model.Order) {
	fmt.Println("------------------")
	fmt.Println(o.Order_uid)
	fmt.Println(o.Track_number)
	fmt.Println(o.Entry)
	fmt.Println(o.D.Delivery_id)
	fmt.Println(o.D.Name)
	fmt.Println(o.D.Phone)
	fmt.Println(o.D.Zip)
	fmt.Println(o.D.City)
	fmt.Println(o.D.Address)
	fmt.Println(o.D.Region)
	fmt.Println(o.D.Email)
	fmt.Println(o.P.Transaction)
	fmt.Println(o.P.Request_id)
	fmt.Println(o.P.Currency)
	fmt.Println(o.P.Provider)
	fmt.Println(o.P.Amount)
	fmt.Println(o.P.Payment_dt)
	fmt.Println(o.P.Bank)
	fmt.Println(o.P.Delivery_cost)
	fmt.Println(o.P.Goods_total)
	fmt.Println(o.P.Custom_fee)
	for _, i := range o.Items {
		fmt.Println(i)
	}
	fmt.Println(o.Locale)
	fmt.Println(o.Internal_signature)
	fmt.Println(o.Customer_id)
	fmt.Println(o.Delivery_service)
	fmt.Println(o.Shardkey)
	fmt.Println(o.Sm_id)
	fmt.Println(o.Date_created)
	fmt.Println(o.Oof_shard)
}

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/VanLavr/L0/model"

	stan "github.com/nats-io/stan.go"
)

const Alph = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"

func main() {
	// Connect to NATS Streaming
	sc, err := stan.Connect("test-cluster", "publisher")
	if err != nil {
		log.Fatal(err)
	}
	defer sc.Close()

	// Publish a message
	for i := 0; i < 1; i++ {
		/*
			if i%10 == 0 {
				err = sc.Publish("model-channel", []byte("invalid message"))
				if err != nil {
					log.Fatal(err)
				}
				continue
			} else if i%15 == 0 {
				err = sc.Publish("model-channel", GenerateIvalidJSON())
				if err != nil {
					log.Fatal(err)
				}
			}
		*/

		model := GenerateModel()
		err = sc.Publish("model-channel", CastTojson(model))
		if err != nil {
			log.Fatal(err)
		}
		time.Sleep(time.Millisecond * 120)
		fmt.Println(model)
	}
}

func GenerateModel() model.Order {
	numberOfItems := rand.Intn(10)
	var items []model.Item
	for i := 0; i < numberOfItems; i++ {
		items = append(items, GenerateRandomItem())
	}

	return model.Order{
		Track_number: GenerateRandomString(),
		Entry:        GenerateRandomString(),
		D: model.Delivery{
			Name:    GenerateRandomString(),
			Phone:   GenerateRandomString(),
			Zip:     GenerateRandomString(),
			City:    GenerateRandomString(),
			Address: GenerateRandomString(),
			Region:  GenerateRandomString(),
			Email:   GenerateRandomString(),
		},
		P: model.Payment{
			Transaction:   GenerateRandomString(),
			Request_id:    GenerateRandomString(),
			Currency:      GenerateRandomString(),
			Provider:      GenerateRandomString(),
			Amount:        rand.Float32(),
			Payment_dt:    rand.Intn(100000),
			Bank:          GenerateRandomString(),
			Delivery_cost: rand.Float32(),
			Goods_total:   rand.Float32(),
			Custom_fee:    rand.Float32(),
		},
		Items:              items,
		Locale:             GenerateRandomString(),
		Internal_signature: GenerateRandomString(),
		Customer_id:        GenerateRandomString(),
		Delivery_service:   GenerateRandomString(),
		Shardkey:           GenerateRandomString(),
		Sm_id:              rand.Intn(100000),
		Date_created:       GenerateRandomString(),
		Oof_shard:          GenerateRandomString(),
	}
}

func GenerateRandomItem() model.Item {
	return model.Item{
		Track_number: GenerateRandomString(),
		Price:        rand.Float32(),
		Rid:          GenerateRandomString(),
		Name:         GenerateRandomString(),
		Sale:         rand.Float32(),
		Size:         GenerateRandomString(),
		Total_Price:  rand.Float32(),
		Nm_id:        rand.Intn(100000),
		Brand:        GenerateRandomString(),
		Status:       rand.Intn(100000),
	}
}

func GenerateRandomString() string {
	length := 10
	bytes := make([]byte, length)
	for i := 0; i < length; i++ {
		bytes[i] = Alph[rand.Intn(len(Alph))]
	}
	return string(bytes)
}

func CastTojson(order model.Order) []byte {
	b, err := json.Marshal(order)
	if err != nil {
		log.Fatal(err)
	}
	return b
}

func GenerateIvalidJSON() []byte {
	type A struct {
		Id   int    `json:"id"`
		Name string `json:"name"`
	}

	a := A{Id: rand.Int(), Name: GenerateRandomString()}

	b, err := json.Marshal(a)
	if err != nil {
		log.Fatal(err)
	}

	return b
}

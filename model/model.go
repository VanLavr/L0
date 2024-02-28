package model

type Order struct {
	Order_uid          string   `json:"order_uid"`
	Track_number       string   `json:"track_number" validate:"required"`
	Entry              string   `json:"entry" validate:"required"`
	D                  Delivery `json:"delivery" validate:"required"`
	P                  Payment  `json:"payment" validate:"required"`
	Items              []Item   `json:"items" validate:"required"`
	Locale             string   `json:"locale" validate:"required"`
	Internal_signature string   `json:"internal_signature" validate:"required"`
	Customer_id        string   `json:"customer_id" validate:"required"`
	Delivery_service   string   `json:"delivery_service" validate:"required"`
	Shardkey           string   `json:"shardkey" validate:"required"`
	Sm_id              int      `json:"sm_id" validate:"required"`
	Date_created       string   `json:"date_created" validate:"required"`
	Oof_shard          string   `json:"oof_shard" validate:"required"`
}

type Delivery struct {
	Delivery_id int    `json:"delivery_id"`
	Name        string `json:"name" validate:"required"`
	Phone       string `json:"phone" validate:"required"`
	Zip         string `json:"zip" validate:"required"`
	City        string `json:"city" validate:"required"`
	Address     string `json:"address" validate:"required"`
	Region      string `json:"region" validate:"required"`
	Email       string `json:"email" validate:"required"`
}

type Payment struct {
	Transaction   string  `json:"transaction" validate:"required"`
	Request_id    string  `json:"request_id" validate:"required"`
	Currency      string  `json:"currency" validate:"required"`
	Provider      string  `json:"provider" validate:"required"`
	Amount        float32 `json:"amount" validate:"required"`
	Payment_dt    int     `json:"payment_dt" validate:"required"`
	Bank          string  `json:"bank" validate:"required"`
	Delivery_cost float32 `json:"delivery_cost" validate:"required"`
	Goods_total   float32 `json:"goods_total" validate:"required"`
	Custom_fee    float32 `json:"custom_fee" validate:"required"`
}

type Item struct {
	Chrt_id      int     `json:"chrt_id" validate:"required"` // set as postgres to serial (do not send this data)
	Track_number string  `json:"track_number" validate:"required"`
	Price        float32 `json:"price" validate:"required"`
	Rid          string  `json:"rid" validate:"required"`
	Name         string  `json:"name" validate:"required"`
	Sale         float32 `json:"sale" validate:"required"`
	Size         string  `json:"size" validate:"required"`
	Total_Price  float32 `json:"total_price" validate:"required"`
	Nm_id        int     `json:"nm_id" validate:"required"`
	Brand        string  `json:"brand" validate:"required"`
	Status       int     `json:"status" validate:"required"`
}

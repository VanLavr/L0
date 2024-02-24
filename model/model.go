package model

type Order struct {
	Track_number       string   `json:"track_number"`
	Entry              string   `json:"entry"`
	D                  Delivery `json:"delivery"`
	P                  Payment  `json:"payment"`
	Items              []Item   `json:"items"`
	Locale             string   `json:"locale"`
	Internal_signature string   `json:"internal_signature"`
	Customer_id        string   `json:"customer_id"`
	Delivery_service   string   `json:"delivery_service"`
	Shardkey           string   `json:"shardkey"`
	Sm_id              int      `json:"sm_id"`
	Date_created       string   `json:"date_created"`
	Oof_shard          string   `json:"oof_shard"`
}

type Delivery struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Zip     string `json:"zip"`
	City    string `json:"city"`
	Address string `json:"address"`
	Region  string `json:"region"`
	Email   string `json:"email"`
}

type Payment struct {
	Transaction   string  `json:"transaction"`
	Request_id    string  `json:"request_id"`
	Currency      string  `json:"currency"`
	Provider      string  `json:"provider"`
	Amount        float32 `json:"amount"`
	Payment_dt    int     `json:"payment_dt"`
	Bank          string  `json:"bank"`
	Delivery_cost float32 `json:"delivery_cost"`
	Goods_total   float32 `json:"goods_total"`
	Custom_fee    float32 `json:"custom_fee"`
}

type Item struct {
	Chrt_id      int     `json:"chrt_id"`
	Track_number string  `json:"track_number"`
	Price        float32 `json:"price"`
	Rid          string  `json:"rid"`
	Name         string  `json:"name"`
	Sale         float32 `json:"sale"`
	Size         string  `json:"size"`
	Total_Price  float32 `json:"total_price"`
	Nm_id        int     `json:"nm_id"`
	Brand        string  `json:"brand"`
	Status       int     `json:"status"`
}

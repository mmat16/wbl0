package repository

import (
	"time"

	"gorm.io/gorm"
)

type DB struct {
	*gorm.DB
}

type Order struct {
	OrderUID          string    `json:"order_uid"          gorm:"unique"`
	TrackNumber       string    `json:"track_number"       gorm:"primaryKey"`
	Entry             string    `json:"entry"`
	Delivery          Delivery  `json:"delivery"           gorm:"embedded"`
	Payment           Payment   `json:"payment"            gorm:"embedded"`
	Items             []Item    `json:"items"              gorm:"foreignKey:track_number"`
	Locale            string    `json:"locale"`
	InternalSignature string    `json:"internal_signature"`
	CustomerID        string    `json:"customer_id"`
	DeliveryService   string    `json:"delivery_service"`
	ShardKey          string    `json:"shard_key"`
	SmID              int       `json:"sm_id"`
	DateCreated       time.Time `json:"date_created"`
	OofShard          string    `json:"oof_shard"`
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
	Transaction  string `json:"transaction"`
	RequestID    string `json:"request_id"`
	Currency     string `json:"currency"`
	Provider     string `json:"provider"`
	Amount       int    `json:"amount"`
	PaymentDt    int64  `json:"payment_dt"`
	Bank         string `json:"bank"`
	DeliveryCost int    `json:"delivery_cost"`
	GoodsTotal   int    `json:"goods_total"`
	CustomFee    int    `json:"custom_fee"`
}

type Item struct {
	gorm.Model  `       json:"-"            gorm:"primaryKey"`
	ChrtID      int    `json:"chrt_id"`
	TrackNumber string `json:"track_number" gorm:"index"`
	Price       int    `json:"price"`
	RID         string `json:"rid"`
	Name        string `json:"name"`
	Sale        int    `json:"sale"`
	Size        string `json:"size"`
	TotalPrice  int    `json:"total_price"`
	NmID        int    `json:"nm_id"`
	Brand       string `json:"brand"`
	Status      int    `json:"status"`
}

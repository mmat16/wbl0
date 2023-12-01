package repository

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect(dsn string) (*DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&Order{}, &Item{})
	if err != nil {
		return nil, err
	}

	return &DB{db}, nil
}

func (db *DB) Insert(order *Order) {
	db.Create(order)
}

func (db *DB) FindAll() []Order {
	var res []Order
	db.Find(&res)
	for i := 0; i < len(res); i++ {
		db.Where("track_number = ?", res[i].TrackNumber).Find(&res[i].Items)
	}
	return res
}

func (Item) TableName() string {
	return "items"
}

func (Order) TableName() string {
	return "orders"
}

func (db *DB) FindByID(uid string) Order {
	var res Order
	db.Where("order_uid = ?", uid).Find(&res)
	db.Where("track_number = ?", res.TrackNumber).Find(&res.Items)
	return res
}

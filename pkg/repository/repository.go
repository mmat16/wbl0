package repository

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Connects to the postgres database and returns pointer to DB and error
// if occured
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

// Inserts given Order to the postgres database
func (db *DB) Insert(order *Order) {
	db.Create(order)
}

// Returns slice of Orders presented in postgres database
func (db *DB) FindAll() []Order {
	var res []Order
	db.Find(&res)
	for i := 0; i < len(res); i++ {
		db.Where("track_number = ?", res[i].TrackNumber).Find(&res[i].Items)
	}
	return res
}

// Retrieves and returns the order from the postgres database by given order_uid
func (db *DB) FindByID(uid string) Order {
	var res Order
	db.Where("order_uid = ?", uid).Find(&res)
	db.Where("track_number = ?", res.TrackNumber).Find(&res.Items)
	return res
}

// Gorm related method
func (Item) TableName() string {
	return "items"
}

// Gorm related method
func (Order) TableName() string {
	return "orders"
}

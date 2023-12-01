package main

import (
	"log"
	"time"

	"github.com/nats-io/stan.go"
)

var msg = []byte(`{
  "order_uid": "b563feb7b2b84b6test",
  "track_number": "WBILMTESTTRACK",
  "entry": "WBIL",
  "delivery": {
    "name": "Test Testov",
    "phone": "+9720000000",
    "zip": "2639809",
    "city": "Kiryat Mozkin",
    "address": "Ploshad Mira 15",
    "region": "Kraiot",
    "email": "test@gmail.com"
  },
  "payment": {
    "transaction": "b563feb7b2b84b6test",
    "request_id": "",
    "currency": "USD",
    "provider": "wbpay",
    "amount": 1817,
    "payment_dt": 1637907727,
    "bank": "alpha",
    "delivery_cost": 1500,
    "goods_total": 317,
    "custom_fee": 0
  },
  "items": [
    {
      "chrt_id": 9934930,
      "track_number": "WBILMTESTTRACK",
      "price": 453,
      "rid": "ab4219087a764ae0btest",
      "name": "Mascaras",
      "sale": 30,
      "size": "0",
      "total_price": 317,
      "nm_id": 2389212,
      "brand": "Vivienne Sabo",
      "status": 202
    }
  ],
  "locale": "en",
  "internal_signature": "",
  "customer_id": "test",
  "delivery_service": "meest",
  "shardkey": "9",
  "sm_id": 99,
  "date_created": "2021-11-26T06:22:19Z",
  "oof_shard": "1"
}`)

var msg4 = []byte(`{
  "order_uid": "b563feb7b2b85b9test",
  "track_number": "WBILMTESTTRACK4",
  "entry": "WBIL",
  "delivery": {
    "name": "Test Testov",
    "phone": "+9720000000",
    "zip": "2639809",
    "city": "Kiryat Mozkin",
    "address": "Ploshad Mira 15",
    "region": "Kraiot",
    "email": "test@gmail.com"
  },
  "payment": {
    "transaction": "b563feb7b2b85b9test",
    "request_id": "",
    "currency": "USD",
    "provider": "wbpay",
    "amount": 2700,
    "payment_dt": 1637907727,
    "bank": "bankofamerica",
    "delivery_cost": 1500,
    "goods_total": 1200,
    "custom_fee": 0
  },
  "items": [
    {
      "chrt_id": 9934934,
      "track_number": "WBILMTESTTRACK4",
      "price": 1200,
      "rid": "ab4219087a764ae4btest",
      "name": "Tracking Shoes",
      "sale": 0,
      "size": "45",
      "total_price": 1200,
      "nm_id": 2389212,
      "brand": "Salomon",
      "status": 202
    }
  ],
  "locale": "en",
  "internal_signature": "",
  "customer_id": "test",
  "delivery_service": "meest",
  "shardkey": "9",
  "sm_id": 99,
  "date_created": "2021-11-26T06:22:19Z",
  "oof_shard": "1"
}`)

var msg3 = []byte(`{
  "order_uid": "b563feb7b2b84b9test",
  "track_number": "WBILMTESTTRACK3",
  "entry": "WBIL",
  "delivery": {
    "name": "Test Testov",
    "phone": "+9720000000",
    "zip": "2639809",
    "city": "Kiryat Mozkin",
    "address": "Ploshad Mira 15",
    "region": "Kraiot",
    "email": "test@gmail.com"
  },
  "payment": {
    "transaction": "b563feb7b2b84b9test",
    "request_id": "",
    "currency": "USD",
    "provider": "wbpay",
    "amount": 3500,
    "payment_dt": 1637907727,
    "bank": "sber",
    "delivery_cost": 1500,
    "goods_total": 2000,
    "custom_fee": 0
  },
  "items": [
    {
      "chrt_id": 9934933,
      "track_number": "WBILMTESTTRACK3",
      "price": 2000,
      "rid": "ab4219087a764ae3btest",
      "name": "Windstopper trousers",
      "sale": 0,
      "size": "46",
      "total_price": 2000,
      "nm_id": 2389212,
      "brand": "The North Face",
      "status": 202
    }
  ],
  "locale": "en",
  "internal_signature": "",
  "customer_id": "test",
  "delivery_service": "meest",
  "shardkey": "9",
  "sm_id": 99,
  "date_created": "2021-11-26T06:22:19Z",
  "oof_shard": "1"
}`)

var msg2 = []byte(`{
  "order_uid": "b563feb7b2b84b8test",
  "track_number": "WBILMTESTTRACK2",
  "entry": "WBIL",
  "delivery": {
    "name": "Test Testov",
    "phone": "+9720000000",
    "zip": "2639809",
    "city": "Kiryat Mozkin",
    "address": "Ploshad Mira 15",
    "region": "Kraiot",
    "email": "test@gmail.com"
  },
  "payment": {
    "transaction": "b563feb7b2b84b8test",
    "request_id": "",
    "currency": "USD",
    "provider": "wbpay",
    "amount": 2800,
    "payment_dt": 1637907727,
    "bank": "alpha",
    "delivery_cost": 1500,
    "goods_total": 1300,
    "custom_fee": 0
  },
  "items": [
    {
      "chrt_id": 9934932,
      "track_number": "WBILMTESTTRACK2",
      "price": 1300,
      "rid": "ab4219087a764ae2btest",
      "name": "Sun Glasses",
      "sale": 0,
      "size": "0",
      "total_price": 1300,
      "nm_id": 2389212,
      "brand": "Ocean",
      "status": 202
    }
  ],
  "locale": "en",
  "internal_signature": "",
  "customer_id": "test",
  "delivery_service": "meest",
  "shardkey": "9",
  "sm_id": 99,
  "date_created": "2021-11-26T06:22:19Z",
  "oof_shard": "1"
}`)

var msg1 = []byte(`{
  "order_uid": "b563feb7b2b84b7test",
  "track_number": "WBILMTESTTRACK1",
  "entry": "WBIL",
  "delivery": {
    "name": "Test Testov",
    "phone": "+9720000000",
    "zip": "2639809",
    "city": "Kiryat Mozkin",
    "address": "Ploshad Mira 15",
    "region": "Kraiot",
    "email": "test@gmail.com"
  },
  "payment": {
    "transaction": "b563feb7b2b84b7test",
    "request_id": "",
    "currency": "USD",
    "provider": "wbpay",
    "amount": 3800,
    "payment_dt": 1637907727,
    "bank": "tink",
    "delivery_cost": 1500,
    "goods_total": 2300,
    "custom_fee": 0
  },
  "items": [
    {
      "chrt_id": 9934931,
      "track_number": "WBILMTESTTRACK1",
      "price": 800,
      "rid": "ab4219087a764ae1btest",
      "name": "Hat",
      "sale": 0,
      "size": "0",
      "total_price": 800,
      "nm_id": 2389212,
      "brand": "Arcteryx",
      "status": 202
    },
    {
      "chrt_id": 9934931,
      "track_number": "WBILMTESTTRACK1",
      "price": 1500,
      "rid": "ab4219087a764ae1btest",
      "name": "Softshell",
      "sale": 0,
      "size": "46",
      "total_price": 1500,
      "nm_id": 23892114,
      "brand": "Arcteryx",
      "status": 202
    }
  ],
  "locale": "en",
  "internal_signature": "",
  "customer_id": "test",
  "delivery_service": "meest",
  "shardkey": "9",
  "sm_id": 99,
  "date_created": "2021-11-26T06:22:19Z",
  "oof_shard": "1"
}`)

func main() {
	msgs := [][]byte{msg, msg1, msg2, msg3, msg4}
	sc, err := stan.Connect("wb", "pub")
	if err != nil {
		log.Fatal(err)
	}
	defer sc.Close()

	for i := 0; i < 5; i++ {
		err = sc.Publish("order", msgs[i])
		if err != nil {
			log.Println(err)
			return
		}
		time.Sleep(2 * time.Second)
		log.Println("sent")
	}
}

package reciever

import (
	"encoding/json"
	"log"
	"os"
	"os/signal"

	"github.com/nats-io/stan.go"

	"demoService/pkg/cache"
	"demoService/pkg/repository"
)

type Receiver struct {
	db     *repository.DB
	cacher *cache.Cache
	sc     stan.Conn
}

func New(db *repository.DB, cacher *cache.Cache, sc stan.Conn) *Receiver {
	return &Receiver{
		db:     db,
		cacher: cacher,
		sc:     sc,
	}
}

func (r *Receiver) UpdateCache() {
	orders := r.db.FindAll()
	for _, order := range orders {
		r.cacher.Set(order.OrderUID, order, cache.NoExpire)
	}
}

func (r *Receiver) Receive() {
	sub, err := r.sc.QueueSubscribe("order", "", func(msg *stan.Msg) {
		log.Println("received message!")
		var order repository.Order
		err := json.Unmarshal(msg.Data, &order)
		if err != nil {
			log.Println("error unmarshalling message. possibly invalid data structure: ", err)
			return
		}
		r.db.Insert(&order)
		r.cacher.Set(order.OrderUID, order, cache.NoExpire)
	})
	if err != nil {
		sub.Unsubscribe()
		r.sc.Close()
		log.Fatal(err)
	}
	signalChan := make(chan os.Signal, 1)
	cleanupDone := make(chan bool)
	signal.Notify(signalChan, os.Interrupt)
	go func() {
		for range signalChan {
			log.Println("Received signal interrupt, attempting graceful shutdown...")
			sub.Unsubscribe()
			cleanupDone <- true
		}
	}()
	<-cleanupDone
	log.Println("Receiver stopped gracefully")
}

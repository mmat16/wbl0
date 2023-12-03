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

// Receiver instance contains pointers to connected database and cacher
// and NATS-Streaming connection
type Receiver struct {
	db     *repository.DB
	cacher *cache.Cache
	sc     stan.Conn
}

// Creates and returns pointer to new Receiver instance
func New(db *repository.DB, cacher *cache.Cache, sc stan.Conn) *Receiver {
	return &Receiver{
		db:     db,
		cacher: cacher,
		sc:     sc,
	}
}

// Finds all entries in the database and writes them to the cache
func (r *Receiver) UpdateCache() {
	orders := r.db.FindAll()
	for _, order := range orders {
		r.cacher.Set(order.OrderUID, order, cache.NoExpire)
	}
}

// Subscibes to the NATS-Streaming subject named "order" and waits for new
// messages. All new messages would be written in database and cache
func (r *Receiver) Receive() {
	sub, err := r.sc.QueueSubscribe("order", "", func(msg *stan.Msg) {
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

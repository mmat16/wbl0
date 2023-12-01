package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nats-io/stan.go"

	"demoService/pkg/cache"
	"demoService/pkg/reciever"
	"demoService/pkg/repository"
)

func main() {
	sc, err := stan.Connect("wb", "sub")
	if err != nil {
		log.Fatal(err)
	}
	defer sc.Close()

	dsn := "host=localhost port=9920 user=v_util password=admin dbname=postgres sslmode=disable TimeZone=Europe/Moscow"
	db, err := repository.Connect(dsn)
	if err != nil {
		sc.Close()
		log.Fatal(err)
	}

	cacher := cache.New(cache.NoExpire, cache.NoCleanup)

	reciever := reciever.New(db, cacher, sc)
	reciever.UpdateCache()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		reciever.Receive()
	}()

	r := gin.Default()
	r.GET("/get_order", func(ctx *gin.Context) {
		order, ok := cacher.Get(ctx.Query("order_uid"))
		if !ok {
			ctx.IndentedJSON(
				404,
				gin.H{
					"error": "order not found",
				},
			)
			return
		}
		ctx.IndentedJSON(200, order)
	})

	// r.Run("localhost:8080")
	server := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Println(err)
			return
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = server.Shutdown(ctx)
	if err != nil {
		panic(err)
	}
	wg.Wait()
}

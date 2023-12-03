package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nats-io/stan.go"

	"demoService/pkg/cache"
	"demoService/pkg/reciever"
	"demoService/pkg/repository"
)

func Home(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "index.tmpl", nil)
}

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

	go func() {
		reciever.Receive()
	}()

	r := gin.Default()
	r.LoadHTMLGlob("./pkg/templates/*")

	r.GET("/", Home)

	r.POST("/proccess", func(ctx *gin.Context) {
		inputText, ok := ctx.GetPostForm("id")
		if !ok {
			ctx.HTML(404, "notfound.tmpl", nil)
			return
		}
		order, ok := cacher.Get(inputText)
		if !ok {
			ctx.HTML(404, "notfound.tmpl", nil)
			return
		}
		ctx.HTML(200, "result.tmpl", order)
	})

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
}

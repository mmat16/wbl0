package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/nats-io/stan.go"

	"demoService/pkg/cache"
	"demoService/pkg/reciever"
	"demoService/pkg/repository"
)

var cacher *cache.Cache

func main() {
	sc, err := connectToNATS()
	if err != nil {
		log.Fatal(err)
	}
	defer sc.Close()

	db, err := connectToDatabase()
	if err != nil {
		sc.Close()
		log.Fatal(err)
	}

	cacher = initCache()

	rec := initReceiver(db, cacher, sc)
	rec.UpdateCache()
	go func() { rec.Receive() }()

	startHTTPServer()
}

// connectToNATS connects to the NATS Streaming Server.
// returns stan.Conn and error
func connectToNATS() (stan.Conn, error) {
	sc, err := stan.Connect("wb", "sub")
	return sc, err
}

// connectToDatabase loads .env file, forms an dsn string to connect to the
// PostgreSQL database.
// returns pointer to repository.DB and error
func connectToDatabase() (*repository.DB, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	passw := os.Getenv("DB_PASSW")
	name := os.Getenv("DB_NAME")
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=Europe/Moscow",
		host,
		port,
		user,
		passw,
		name,
	)
	return repository.Connect(dsn)
}

// initCache initializes the cache.
func initCache() *cache.Cache {
	return cache.New(cache.NoExpire, cache.NoCleanup)
}

// initReceiver initializes the message receiver.
func initReceiver(db *repository.DB, cacher *cache.Cache, sc stan.Conn) *reciever.Receiver {
	rec := reciever.New(db, cacher, sc)
	return rec
}

// startHTTPServer starts the HTTP server.
func startHTTPServer() {
	r := setupRoutes()

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

	mustWaitForShutdown(server)
}

// setupRoutes configures the HTTP routes.
// returns pointer to gin.Engine
func setupRoutes() *gin.Engine {
	r := gin.Default()
	r.LoadHTMLGlob("./pkg/templates/*")

	r.GET("/", home)

	r.POST("/proccess", handleProcessRequest)

	return r
}

// Home handles GET request and writes HTML to the *gin.Context
func home(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "index.tmpl", nil)
}

// handleProcessRequest handles the POST request for processing.
// which searches the given in request order_uid in cache and
// writes it as HTML to *gin.Context
func handleProcessRequest(ctx *gin.Context) {
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
}

// waitForShutdown waits for an OS interrupt signal and performs a graceful shutdown.
// may cause a panic if some error occurs during server shutdown
func mustWaitForShutdown(server *http.Server) {
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := server.Shutdown(ctx)
	if err != nil {
		panic(err)
	}
}

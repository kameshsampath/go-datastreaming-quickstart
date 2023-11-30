package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/kameshsampath/examples/handlers"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	var seeds string
	flag.StringVar(&seeds, "brokers", "127.0.0.1:9092", "comma separated list of broker addresses")
	flag.Parse()

	brokers := handlers.Seeds(strings.Split(seeds, ","))

	e := echo.New()
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())
	//middleware to set context variables that will be available to all contexts
	e.Use(brokers.SeedsContext())

	// API
	e.GET("/", handlers.Consume)
	e.POST("/", handlers.Produce)

	go func() {
		if err := e.Start(":8080"); err != nil && err == http.ErrServerClosed {
			log.Fatalf("shutting down server, %s", err)
		}
	}()

	//Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
}

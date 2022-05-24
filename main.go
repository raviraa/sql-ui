package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"sql-ui/config"
	"sql-ui/routes"
	"sql-ui/services"
)

func main() {
	// Start a new container
	c := services.NewContainer()
	defer func() {
		if err := c.Shutdown(); err != nil {
			log.Fatal(err)
		}
	}()

  addr := "127.0.0.1:8080"
	// Build the router
	routes.BuildRouter(c)
	srv := http.Server{
		Addr: addr,
		Handler:      c.Web,
		ReadTimeout:  config.HttpTimeout,
		WriteTimeout: config.HttpTimeout,
	}
  log.Println("Listening on server: ", addr)

	// Start the server
	go func() {
		// if err := c.Web.Run(":8080"); err != http.ErrServerClosed {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	quit := make(chan os.Signal, 1)

	signal.Notify(quit, os.Interrupt)
	signal.Notify(quit, syscall.SIGTERM)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
}

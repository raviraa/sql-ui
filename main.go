package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/raviraa/sql-ui/config"
	"github.com/raviraa/sql-ui/routes"
	"github.com/raviraa/sql-ui/services"
)

//addr := "127.0.0.1:9292"
var addr = flag.String("a", "127.0.0.1:9292", "address to listen on ")

func main() {
	flag.Parse()
	// Start a new container
	c := services.NewContainer()
	defer func() {
		if err := c.Shutdown(); err != nil {
			log.Fatal(err)
		}
	}()

	// Build the router
	routes.BuildRouter(c)
	srv := http.Server{
		Addr:         *addr,
		Handler:      c.Web,
		ReadTimeout:  config.HttpTimeout,
		WriteTimeout: config.HttpTimeout,
	}
	log.Println("Listening on server: ", *addr)

	// Start the server
	go func() {
		// if err := c.Web.Run(":8080"); err != http.ErrServerClosed {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	if c.Config.OpenInWebBrowser {
		time.Sleep(time.Second * 3)
		openBrowser("http://" + *addr)
	}

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

func openBrowser(url string) {
	fmt.Println("Opening url: ", url)
	var err error
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Println("Unable to open in browser ", err, "\nPlease open url in browser: ", url)
	}
}

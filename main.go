package main

import (
	"log"
	"net/http"

	"sql-ui/routes"
	"sql-ui/services"
)

func main() {
	log.SetFlags(log.Lshortfile)
	// Start a new container
	c := services.NewContainer()
	defer func() {
		if err := c.Shutdown(); err != nil {
			log.Fatal(err)
		}
	}()

	// Build the router
	routes.BuildRouter(c)

	// Start the server
	// go func() {
		// srv := http.Server{
		//     Addr: ":8080",
		// Addr:         fmt.Sprintf("%s:%d", c.Config.HTTP.Hostname, c.Config.HTTP.Port),
		// Handler:      c.Web,
		// ReadTimeout:  c.Config.HTTP.ReadTimeout,
		// WriteTimeout: c.Config.HTTP.WriteTimeout,
		// IdleTimeout:  c.Config.HTTP.IdleTimeout,
		// }

		if err := c.Web.Run(":8080"); err != http.ErrServerClosed {
			log.Fatal(err)
		}
	// }()

	// Start the scheduler service to queue periodic tasks
	// go func() {
	// 	if err := c.Tasks.StartScheduler(); err != nil {
	// 		c.Web.Logger.Fatalf("scheduler shutdown: %v", err)
	// 	}
	// }()
	//
	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	// quit := make(chan os.Signal, 1)
	// signal.Notify(quit, os.Interrupt)
	// signal.Notify(quit, os.Kill)
	// <-quit
	// ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	// defer cancel()
	// if err := c.Web.Shutdown(ctx); err != nil {
	// 	c.Web.Logger.Fatal(err)
	// }
}

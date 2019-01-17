package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/xujiajun/gorouter"
)

func main() {
	logger := log.New(os.Stdout, "[gorouter] ", log.Ldate|log.Ltime)

	mux := gorouter.New()
	mux.GET("/", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second)
		w.Write([]byte("hello world"))
		logger.Println("Handle request success")
	})

	srv := &http.Server{
		Addr:    ":8181",
		Handler: mux,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			logger.Printf("listen: %s\n", err)
		}
	}()

	graceful(srv, logger, 3*time.Second)
}

// reference: https://gist.github.com/peterhellberg/38117e546c217960747aacf689af3dc2
func graceful(hs *http.Server, logger *log.Logger, timeout time.Duration) {
	stop := make(chan os.Signal, 1)
	// Handle SIGINT and SIGTERM.
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	logger.Printf("Shutdown with timeout: %s\n", timeout)
	// Stop the service gracefully.
	if err := hs.Shutdown(ctx); err != nil {
		logger.Printf("Error: %v\n", err)
	} else {
		logger.Println("Server stopped")
	}
}

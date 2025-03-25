package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"mini-kafka/handler"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/slow", handler.DelayedResponse(7*time.Second))
	mux.HandleFunc("/veryslow", handler.DelayedResponse(20*time.Second))

	srv := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	// Run the server in a separate goroutine so that it doesn't block
	go func() {
		log.Println("server starting on :8080")
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	stopChan := make(chan os.Signal, 1)                    // Create a channel to receive OS signals
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM) // Notify the channel for the interrupt signal

	<-stopChan // This will block the main goroutine until a signal is received

	// Begin shutdown process
	log.Println("shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second) // give the server 10 seconds to gracefully shutdown
	defer cancel()                                                           // ensure that the context is canceled to release resources
	if err := srv.Shutdown(ctx); err != nil {                                // Shutdown the server by providing the context
		log.Fatalf("shutdown error: %s\n", err)
	}
	log.Println("server gracefully stopped")
}

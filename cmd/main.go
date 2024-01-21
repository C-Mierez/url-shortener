package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/c-mierez/url-shortener/internal/handlers"
	"github.com/go-chi/chi/v5"
)

func main() {

	// Create a new router
	r := chi.NewRouter()

	// Create a new logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	// Create a global context to wrap all others
	serverCtx, serverCtxCancel := context.WithCancel(context.Background())

	// Signal handling
	killSig := make(chan os.Signal, 1)
	signal.Notify(killSig, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	// Create a new server
	svr := http.Server{
		Addr:    "127.0.0.1:8080",
		Handler: r, // Chi router
	}

	// Graceful shutdown
	go func() {
		sig := <-killSig // Block until a signal is received

		logger.Info("Received kil signal. Shutting down...", slog.String("signal", sig.String()))

		shutdownCtx, shutdownCtxCancel := context.WithTimeout(serverCtx, 5*time.Second)
		defer shutdownCtxCancel()

		go func() {
			<-shutdownCtx.Done() // Block until the shutdown context is cancelled

			if shutdownCtx.Err() == context.DeadlineExceeded {
				log.Fatal("Graceful shutdown timed out")
			}
		}()

		if err := svr.Shutdown(shutdownCtx); err != nil {
			log.Fatal(err)
		}

		serverCtxCancel()
		logger.Info("Shutting down server...")
	}()

	// Start the server in a goroutine
	go func() {
		if err := svr.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	// Routes
	r.Get("/health", handlers.NewHealthCheckHandler().ServeHTTP)

	<-serverCtx.Done() // Block until the server context is cancelled
}

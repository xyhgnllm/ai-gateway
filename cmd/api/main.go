package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"ai-gateway/internal/config"
	"ai-gateway/internal/database"
	"ai-gateway/internal/db"
	"ai-gateway/internal/router"
)

func main() {
	cfg := config.Load()

	if err := cfg.Validate(); err != nil {
		log.Fatal(err)
	}

	pool := db.New(cfg.DBURL)
	defer pool.Close()

	queries := database.New(pool)

	r := router.New(
		queries,
		pool,
		cfg.JWTSecret,
		cfg.OpenAIBaseURL,
		cfg.OpenAIAPIKey,
	)

	server := &http.Server{
		Addr:    ":" + cfg.AppPort,
		Handler: r,
	}

	go func() {
		log.Printf("server running on :%s", cfg.AppPort)

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit
	log.Println("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}

	log.Println("server stopped")
}

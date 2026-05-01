package main

import (
	"ai-gateway/internal/config"
	"ai-gateway/internal/database"
	"ai-gateway/internal/db"
	"ai-gateway/internal/router"
	"log"
	"net/http"
)

func main() {
	cfg := config.Load()

	pool := db.New(cfg.DBURL)
	defer pool.Close()

	queries := database.New(pool)
	r := router.New(queries, cfg.JWTSercet)

	log.Printf("server running on :%s", cfg.AppPort)
	err := http.ListenAndServe(":"+cfg.AppPort, r)
	if err != nil {
		log.Fatal(err)
	}
}

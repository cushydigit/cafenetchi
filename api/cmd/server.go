package main

import (
	"cafenetchi-api/internal/config"
	"cafenetchi-api/internal/redis"
	"context"
	"fmt"
	"log"
	"net/http"
)

func main() {
	// global context
	ctx := context.Background()
	// load configuration
	cfg := config.Load()
	// Initialize Database
	db := config.InitDB(cfg)

	// Initialize Cache (Redis)
	redis.Init(ctx, fmt.Sprintf("%s:%s", cfg.RdsHost, cfg.RdsPort), cfg.RdsPass, 0)
	defer redis.Close()

	// Auto Migrate (for early development)
	if err := config.AutoMigrate(db); err != nil {
		log.Printf("migration failed: %v", err)
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf("%s", ":"+cfg.ServerPort),
		Handler: Routes(*cfg),
	}

	log.Printf("server is up and listen at %s", cfg.ServerPort)

	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("failed to listen or serve ... %v", err)
	}
}

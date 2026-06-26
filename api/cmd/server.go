package main

import (
	"cafenetchi-api/internal/config"
	"cafenetchi-api/internal/redis"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	// load configuration
	cfg := config.Load()
	// Initialize Database
	db := config.InitDB(cfg)

	// Initialize Cache (Redis)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	redis.Init(
		ctx,
		fmt.Sprintf("%s:%s", cfg.RdsHost, cfg.RdsPort),
		cfg.RdsPass, 0,
	)
	defer redis.Close()

	// Auto Migrate (for early development)
	if err := config.AutoMigrate(db); err != nil {
		log.Printf("migration failed: %v", err)
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf("%s", ":"+cfg.ServerPort),
		Handler: Routes(*cfg),
	}

	go func() {
		log.Printf("server is up and listening at %s", cfg.ServerPort)
		if err := srv.ListenAndServe(); err != nil {
			log.Fatalf("server failed: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("server shutdown failed: %v", err)
	}

	log.Println("server stopped")

}

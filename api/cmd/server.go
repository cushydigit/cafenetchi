package main

import (
	"cafenetchi-api/internal/config"
	"cafenetchi-api/internal/redis"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// load configuration
	cfg := config.Load()

	// Initialize Database
	if err := config.InitPool(context.Background(), cfg.DB.DSN()); err != nil {
		log.Fatal(err)
	}
	defer config.GetPool(context.Background()).Close()

	// Initialize Cache (Redis)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	redis.Init(
		ctx,
		cfg.Redis.RedisAddr(),
		cfg.Redis.Pass,
		0,
	)
	defer redis.Close()

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: Routes(*cfg),
	}

	go func() {
		log.Printf("✅ Server started on :%s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				log.Fatalf("server failed: %v", err)
			}
		}
	}()

	quit, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	<-quit.Done()

	log.Println("shutting down...")

	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("server shutdown failed: %v", err)
	}

	log.Println("Server stopped")

}

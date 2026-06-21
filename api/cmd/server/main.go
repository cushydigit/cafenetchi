package main

import (
	"cafenetchi-api/internal/config"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	chi_middleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func main() {
	// load configuration
	cfg := config.Load()

	// Initialize Database
	db := config.InitDB(cfg)

	// Auto Migrate (for early developement)
	if err := config.AutoMigrate(db); err != nil {
		log.Printf("migration failed: %v", err)
	}

	// Initilize Router
	r := chi.NewRouter()

	// Global Router
	r.Use(chi_middleware.Logger)
	r.Use(chi_middleware.Recoverer)
	// TODO: check is it working?
	r.Use(chi_middleware.Timeout(60 * time.Second)) // is it working?
	r.Use(cors.Handler(cors.Options{
		// TODO: change in production
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Health Check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("cafenetchi api server up and healthy"))
	})

	log.Printf("server is up and listen at %s", cfg.ServerPort)
	if err := http.ListenAndServe(":"+cfg.ServerPort, r); err != nil {
		log.Fatalf("failed to listen or serve ... %v", err)
	}
}

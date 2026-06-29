package main

import (
	"cafenetchi-api/internal/config"
	"cafenetchi-api/internal/logger"

	db "cafenetchi-api/internal/db/generated"
	"cafenetchi-api/internal/handler"
	"cafenetchi-api/internal/service"
	"context"
	"time"

	"github.com/go-chi/chi/v5"
	chi_middleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

// Routes initializes and configures a chi router with the given configuration.
//
// Parameters:
// - cfg: The configuration for the router.
//
// Returns:
// - r: The configured chi router.
func Routes(cfg config.Config) chi.Router {
	// Initialize Router
	r := chi.NewRouter()

	// Global Router
	r.Use(chi_middleware.Logger)
	r.Use(chi_middleware.Recoverer)
	// TODO: check is it working?
	r.Use(chi_middleware.Timeout(60 * time.Second)) // is it working?
	r.Use(cors.Handler(cors.Options{
		// TODO: change in production add domain
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	pool := config.GetPool(context.Background())
	queries := db.New(pool)

	otpSvc := service.NewInRedisOTP()
	smsSvc := service.NewKavenegar(cfg.KavenegarAPIKey, cfg.KavenegarSender)

	l := logger.New()

	s := service.NewAuth(queries, otpSvc, smsSvc, cfg.JWTSecret, l)

	h := handler.NewAuth(s, l)

	// routes
	r.Post("/otp", h.SendOTP)
	r.Post("/verify", h.VerifyOTP)

	return r
}

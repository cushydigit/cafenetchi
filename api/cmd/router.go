package main

import (
	"cafenetchi-api/internal/handler"
	"cafenetchi-api/internal/limiter"
	"cafenetchi-api/internal/logger"
	"cafenetchi-api/internal/middleware"
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
func Routes(
	auth *handler.Auth,
	rateLimiter limiter.Limiter,
	appLogger *logger.Logger,
) chi.Router {
	// Initialize Router
	r := chi.NewRouter()

	// global middlewares
	r.Use(
		chi_middleware.RequestID,
		chi_middleware.Logger,
		chi_middleware.Recoverer,
		chi_middleware.Timeout(60*time.Second),
		middleware.Auth("secret"),
		cors.Handler(cors.Options{
			// TODO: change in production add domain
			AllowedOrigins:   []string{"*"},
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
			AllowCredentials: true,
			MaxAge:           300,
		}),
	)

	// routes
	r.Route("/auth", func(r chi.Router) {
		r.With(middleware.RateLimit(rateLimiter)).Post("/otp", auth.SendOTP)

	})
	r.Post("/otp", auth.SendOTP)
	r.Post("/verify", auth.VerifyOTP)

	return r
}

package main

import (
	"cafenetchi-api/internal/config"
	db "cafenetchi-api/internal/db/generated"
	"cafenetchi-api/internal/handler"
	"cafenetchi-api/internal/logger"
	"cafenetchi-api/internal/otp"
	"cafenetchi-api/internal/redis"
	"cafenetchi-api/internal/repository"
	"cafenetchi-api/internal/service"
	"cafenetchi-api/internal/sms"
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// application context
	appCtx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	// create logger
	appLogger := logger.New()

	// load configuration
	cfg := config.Load(appLogger)

	// create lifecycle for connecting to db and redis
	initCtx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	// database
	dbPool, err := config.NewPool(initCtx, cfg.DB.DSN())
	if err != nil {
		appLogger.Error(
			"database initialization failed",
			"error",
			err,
		)
		return
	}
	defer dbPool.Close()
	queries := db.New(dbPool)

	// redis
	rds, err := redis.New(initCtx, cfg.Redis.RedisAddr(), cfg.Redis.Pass, 0)
	if err != nil {
		appLogger.Error(
			"redis initialization failed",
			"error",
			err,
		)
		return
	}
	defer rds.Close(appCtx)

	// Repositories
	userRepo := repository.NewUser(queries)

	// OTP store
	otpStore := redis.NewOTPStore(rds, 2*time.Minute)
	limitStore := redis.NewRedisLimiter(rds, "otp", 5, time.Minute)

	// services
	otpSvc := otp.New(otpStore)

	smsSvc := sms.NewKavenegar(
		cfg.KavenegarAPIKey,
		cfg.KavenegarSender,
	)

	authService := service.NewAuth(
		userRepo,
		otpSvc,
		smsSvc,
		cfg.JWTSecret,
		appLogger,
	)

	userService := service.NewUser(
		userRepo,
		appLogger,
	)

	// Handlers
	authHandler := handler.NewAuth(
		authService,
	)

	userHandler := handler.NewUser(
		userService,
		appLogger,
	)

	// router
	router := Routes(
		authHandler,
		userHandler,
		limitStore,
		cfg.JWTSecret,
	)

	// server
	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,

		ReadTimeout:       15 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	serverErr := make(chan error, 1)

	go func() {
		appLogger.Info(
			"server started",
			"port",
			cfg.Port,
		)
		serverErr <- srv.ListenAndServe()
	}()

	select {
	case err := <-serverErr:
		appLogger.Error(
			"server crashed",
			"error",
			err,
		)
	case <-appCtx.Done():
		appLogger.Info(
			"shutdown signal received",
		)
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		appLogger.Error(
			"graceful shutdown failed",
			"error",
			err,
		)
		return
	}

	appLogger.Info(
		"server stopped",
	)
}

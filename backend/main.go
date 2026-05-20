package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mflargooo/internal/api/handler"
	"github.com/mflargooo/internal/config"
	"github.com/mflargooo/internal/middleware"
	"github.com/mflargooo/internal/service"
	"github.com/mflargooo/internal/store"
	"github.com/mflargooo/internal/workers"
	"github.com/redis/go-redis/v9"
)

func main() {
	cfg := config.Load()

	// -- Postgres

	db, err := pgxpool.New(context.Background(), cfg.Postgres.FormatDSN())
	if err != nil {
		log.Fatalf("open postgres: %v", err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.Ping(ctx); err != nil {
		log.Fatalf("ping postgres: %s %v", cfg.Postgres.FormatDSN(), err)
	}

	log.Println("[POSTGRES] Connected")

	// -- Redis
	opt, err := redis.ParseURL(cfg.Redis.FormatDSN())
	if err != nil {
		panic(err)
	}
	rdb := redis.NewClient(opt)
	defer rdb.Close()

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatalf("ping redis: %v", err)
	}
	log.Println("[REDIS] Connected")

	// -- Setup
	cameraStore := store.NewCameraStore(db)
	stateStore := store.NewCameraStateStore(rdb)
	ffmpegManager := workers.NewFFmpegManager(cfg.MediaMTX)
	cameraSvc := service.NewCameraService(cameraStore, stateStore, ffmpegManager)
	cameraHandler := handler.NewCameraHandler(cameraSvc)

	// -- Route
	mux := http.NewServeMux()
	cameraHandler.RegisterRoutes(mux)

	server := &http.Server{
		Addr:         ":8080",
		Handler:      middleware.Cors(mux),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("Listening at %s\n", ":8080")
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("listen: %v", err)
	}
}

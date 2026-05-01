package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/fasthttp/router"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"github.com/valyala/fasthttp"

	"github.com/voznikaetnepriyazn/Good-service/internal/config"
	handlers "github.com/voznikaetnepriyazn/Good-service/internal/http-server/handlers/good"
	"github.com/voznikaetnepriyazn/Good-service/internal/lib/logger/sl"
	"github.com/voznikaetnepriyazn/Good-service/internal/storage"
	"github.com/voznikaetnepriyazn/Good-service/internal/storage/postgresql"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	if err := godotenv.Load(); err != nil {
		slog.Error(".env file not found", sl.Err(err))
	}

	cfg := config.MustLoad()

	logger := setUpLogger(cfg.Env)

	if logger == nil {
		logger = slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}))
		slog.SetDefault(logger)
	} else {
		slog.SetDefault(logger)
	}

	logger.Info("starting customer service", slog.String("env", cfg.Env))
	logger.Debug("debug messages are enabled")

	db, err := postgresql.New(cfg.DB.DSN())
	if err != nil {
		slog.Error("failed to connect to database", sl.Err(err))
		os.Exit(1)
	}
	defer db.Close()

	goodService := storage.GoodService(db)

	r := registerRouter(logger, goodService)

	app := &fasthttp.Server{
		Name:    "Customer Service v1.0",
		Handler: r.Handler,
	}

	addr := cfg.HTTPServer.Address
	if addr == "" {
		addr = ":8081"
	}

	logger.Info("starting server", slog.String("address", addr))
	if err := app.ListenAndServe(addr); err != nil {
		logger.Error("failed to start server", slog.Any("error", err))
		os.Exit(1)
	}
}

func registerRouter(log *slog.Logger, service storage.GoodService) *router.Router {
	r := router.New()

	r.GET("/health", func(ctx *fasthttp.RequestCtx) {
		ctx.SetContentType("application/json")
		ctx.SetStatusCode(fasthttp.StatusOK)
		ctx.Write([]byte(`{"status":"ok"}`))
	})

	r.GET("/", func(ctx *fasthttp.RequestCtx) {
		ctx.SetContentType("text/plain; charset=utf-8")
		fmt.Fprintf(ctx, "Welcome to Good Service v1.0")
	})

	api := r.Group("/api/v1")

	api.POST("/good", handlers.NewAdd(log, service))
	api.GET("/good/:id", handlers.NewGetById(log, service))
	api.GET("/good", handlers.NewGetAll(log, service))
	api.PUT("/good/:id/fullName", handlers.NewUpdate(log, service))
	api.DELETE("/good/:id", handlers.NewDelete(log, service))

	api.GET("/good/:id/brand", handlers.NewGetListOfGoodsByBrand(log, service))
	api.GET("/good/:id/type", handlers.NewGetListOfGoodsByType(log, service))
	api.GET("/good/:id/avaliable", handlers.NewIsAvaliableForOrder(log, service))
	api.GET("/good/:id/rest", handlers.NewRestOfGood(log, service))

	return r
}

func setUpLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}

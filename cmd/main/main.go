package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"golang.org/x/sync/errgroup"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"
	"user-segmentation/internal/api/http"
	"user-segmentation/internal/config"
	"user-segmentation/internal/logger"
	"user-segmentation/internal/repo/history"
	"user-segmentation/internal/repo/segments"
	"user-segmentation/internal/service"
)

func main() {
	cfg := config.MustLoad()

	eg, ctx := errgroup.WithContext(context.Background())
	log := logger.Create(cfg.Env)
	log.Info("starting app", slog.String("env", cfg.Env))

	conn, err := pgx.Connect(ctx, cfg.DbConn)
	for i := 0; i < 5 && err != nil; i++ {
		time.Sleep(time.Second * 3)
		log.Info("reconnect to PostgreSQL", slog.Int("attempt", i+1))
		conn, err = pgx.Connect(ctx, cfg.DbConn)
	}
	if err != nil {
		log.Error("cannot connect to database", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer func() { _ = conn.Close(ctx) }()
	svc := service.New(segments.New(conn), history.New(conn))

	srv := http.New(log, cfg.HTTPAddr, cfg.Env, svc)
	sigQuit := make(chan os.Signal, 1)
	signal.Ignore(syscall.SIGHUP, syscall.SIGPIPE)
	signal.Notify(sigQuit, syscall.SIGINT, syscall.SIGTERM)

	eg.Go(func() error {
		select {
		case s := <-sigQuit:
			return fmt.Errorf("captured signal: %v", s)
		case <-ctx.Done():
			return nil
		}
	})

	eg.Go(func() (err error) {
		return srv.Listen(ctx)
	})
	if err := eg.Wait(); err != nil {
		log.Error("caught error for graceful shutdown", slog.String("error", err.Error()))
	}
	log.Info("server has been shutdown successfully")
}

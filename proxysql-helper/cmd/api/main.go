package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"proxysql-helper/internal/app"
	"proxysql-helper/internal/config"
	"proxysql-helper/internal/pool"
	"proxysql-helper/internal/repository"
	"proxysql-helper/internal/router"
)

func main() {
	cfg := config.Load()
	var strat pool.RouterStrategy
	switch cfg.RouterStrategy {
	case "lowest_latency":
		strat = router.LowestLatency{}
	case "random":
		strat = router.Random{}
	default:
		strat = &router.RoundRobin{}
	}

	rp, err := pool.NewRouterPool(cfg.ProxySQLEndpoints, cfg.User, cfg.Password, cfg.DBName, cfg.MaxOpenConns, cfg.MaxIdleConns, strat)
	if err != nil {
		log.Fatalf("router pool: %v", err)
	}

	w := repository.NewMySQLWriter(rp, cfg.DBName)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	if err := w.InitSchema(ctx); err != nil {
		cancel()
		log.Fatalf("init schema: %v", err)
	}
	cancel()

	api := app.NewServer(w)
	server := &http.Server{Addr: cfg.HTTPAddr, Handler: api.Routes()}

	go func() {
		log.Printf("listening on %s | strategy=%s | endpoints=%s", cfg.HTTPAddr, strat.Name(), os.Getenv("PROXYSQL_ENDPOINTS"))
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("http server: %v", err)
		}
	}()

	// graceful shutdown
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
	ctx2, cancel2 := context.WithTimeout(context.Background(), 5*time.Second)
	_ = server.Shutdown(ctx2)
	cancel2()
}

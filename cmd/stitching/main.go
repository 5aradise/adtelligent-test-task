package main

import (
	"flag"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/5aradise/adtelligent-test-task/config"
	stitchingHandler "github.com/5aradise/adtelligent-test-task/internal/controllers/http/stitching"
	stitchingService "github.com/5aradise/adtelligent-test-task/internal/services/stitching"
	"github.com/5aradise/adtelligent-test-task/pkg/httpserver"
	"github.com/5aradise/adtelligent-test-task/pkg/logger"
	"github.com/5aradise/adtelligent-test-task/pkg/middleware"
)

var configPath = *flag.String("config", "./config.yaml", "Path to config file")

func main() {
	flag.Parse()

	cfg, err := config.New(configPath)
	if err != nil {
		log.Fatal("can't load config: ", err)
	}

	l := logger.New(os.Stdout, cfg.Env)

	l.Info("init stitching service")
	ss := stitchingService.New(cfg.Stitching.AuctionUrl, l, cfg.Stitching.RequestTimeout)
	l.Info("init stitching handler")
	sh := stitchingHandler.New(ss, l)

	router := http.NewServeMux()

	sh.Init(router)

	server := httpserver.New(
		middleware.Use(router,
			middleware.Recoverer(l),
			middleware.Cors(l),
			middleware.RequestID(l),
			middleware.Logger(l),
		),
		httpserver.Port(cfg.Server.StitchingPort),
		httpserver.ReadTimeout(cfg.Server.Timeout),
		httpserver.IdleTimeout(cfg.Server.IdleTimeout),
		httpserver.ErrorLog(slog.NewLogLogger(l.With(slog.String("source", "stitching-server")).Handler(), slog.LevelError)),
	)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	l.Info("starting server", slog.String("address", server.Addr()))
	go server.Run()

	select {
	case s := <-interrupt:
		l.Error("signal interrupt", slog.String("error", s.String()))
	case err := <-server.Notify():
		l.Error("server notify", logger.Err(err))
	}

	err = server.Shutdown()
	if err != nil {
		l.Error("can't shutdown server", logger.Err(err))
	}
}

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
	auctionHandler "github.com/5aradise/adtelligent-test-task/internal/controllers/http/auction"
	auctionService "github.com/5aradise/adtelligent-test-task/internal/services/auction"
	"github.com/5aradise/adtelligent-test-task/internal/storage"
	"github.com/5aradise/adtelligent-test-task/pkg/db/postgresql"
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

	conn, err := postgresql.New(cfg.DB.URL)
	if err != nil {
		l.Error("can't open sql", logger.Err(err))
		os.Exit(1)
	}
	defer conn.Close()

	l.Info("init app storage")
	s := storage.New(conn, l, cfg.DB.BigRequestTimeout, cfg.DB.UpdateCacheDelay)

	l.Info("init auction service")
	as := auctionService.New(s, l)
	l.Info("init auction handler")
	ah := auctionHandler.New(as, l)

	router := http.NewServeMux()
	ah.Init(router)

	server := httpserver.New(
		middleware.Use(router,
			middleware.Recoverer(l),
			middleware.Cors(l),
			middleware.RequestID(l),
			middleware.Logger(l),
		),
		httpserver.Port(cfg.Server.AuctionPort),
		httpserver.ReadTimeout(cfg.Server.Timeout),
		httpserver.IdleTimeout(cfg.Server.IdleTimeout),
		httpserver.ErrorLog(slog.NewLogLogger(l.With(slog.String("source", "auction-server")).Handler(), slog.LevelError)),
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

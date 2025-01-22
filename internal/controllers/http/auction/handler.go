package auction

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/5aradise/adtelligent-test-task/internal/models"
)

type auctionService interface {
	GetProfitCreative(sourceId int, maxDuration time.Duration) models.Creative
}

type handler struct {
	s auctionService
	l *slog.Logger
}

func New(s auctionService, logger *slog.Logger) *handler {
	return &handler{s, logger}
}

// /auction?sourceID={source_id}&maxDuration={max_duration}
func (h *handler) Init(router *http.ServeMux) {
	router.HandleFunc("GET /auction", h.handleAuction)
}

func (h *handler) handleAuction(w http.ResponseWriter, r *http.Request) {
}

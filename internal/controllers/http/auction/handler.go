package auction

import (
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/5aradise/adtelligent-test-task/internal/models"
	"github.com/5aradise/adtelligent-test-task/pkg/api"
	"github.com/5aradise/adtelligent-test-task/pkg/middleware"
	"github.com/5aradise/adtelligent-test-task/pkg/util"
)

type auctionService interface {
	GetProfitCreative(sourceId int, maxDuration time.Duration) (models.Creative, error)
}

type handler struct {
	as auctionService
	l  *slog.Logger
}

func New(as auctionService, logger *slog.Logger) *handler {
	return &handler{as, logger}
}

// /auction?sourceID={source_id}&maxDuration={max_duration}
func (h *handler) Init(router *http.ServeMux) {
	router.HandleFunc("GET /auction", h.handleAuction)
}

func (h *handler) handleAuction(w http.ResponseWriter, r *http.Request) {
	const op = "handler.handleAuction"
	l := h.l.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetRequestID(r)),
	)

	query := r.URL.Query()

	sourceIDStr := query.Get("sourceID")
	if sourceIDStr == "" {
		err := api.WriteError(w, http.StatusBadRequest, "sourceID query is empty")
		if err != nil {
			l.Warn("cannot make response", util.SlErr(err))
		}
		return
	}
	maxDurationStr := query.Get("maxDuration")
	if maxDurationStr == "" {
		err := api.WriteError(w, http.StatusBadRequest, "maxDuration query is empty")
		if err != nil {
			l.Warn("cannot make response", util.SlErr(err))
		}
		return
	}

	sourceID, err := strconv.Atoi(sourceIDStr)
	if err != nil {
		err := api.WriteError(w, http.StatusBadRequest, "sourceID query must be a number")
		if err != nil {
			l.Warn("cannot make response", util.SlErr(err))
		}
		return
	}
	maxDuration, err := time.ParseDuration(maxDurationStr)
	if err != nil {
		err := api.WriteError(w, http.StatusBadRequest, "maxDuration query must be valid duration")
		if err != nil {
			l.Warn("cannot make response", util.SlErr(err))
		}
		return
	}

	creative, err := h.as.GetProfitCreative(sourceID, maxDuration)
	if err != nil {
		err := api.WriteError(w, http.StatusBadRequest, err.Error())
		if err != nil {
			l.Warn("cannot make response", util.SlErr(err))
		}
		return
	}

	err = api.WriteJSON(w, http.StatusOK, creative)
	if err != nil {
		l.Warn("cannot make response", util.SlErr(err))
	}
}

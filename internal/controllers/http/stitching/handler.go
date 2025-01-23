package stitching

import (
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/5aradise/adtelligent-test-task/pkg/api"
	"github.com/5aradise/adtelligent-test-task/pkg/middleware"
	"github.com/5aradise/adtelligent-test-task/pkg/util"
)

type stitchingService interface {
	ModifyPlaylist(sourceId int, playlist io.Reader) ([]byte, error)
}

type handler struct {
	ss stitchingService
	l  *slog.Logger
}

func New(ss stitchingService, logger *slog.Logger) *handler {
	return &handler{ss, logger}
}

// /stitching.m3u8?sourceID={source_id}
func (h *handler) Init(router *http.ServeMux) {
	router.HandleFunc("POST /stitching.m3u8", h.handleStitching)
}

func (h *handler) handleStitching(w http.ResponseWriter, r *http.Request) {
	const op = "handler.handleStitching"
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
	sourceID, err := strconv.Atoi(sourceIDStr)
	if err != nil {
		err := api.WriteError(w, http.StatusBadRequest, "sourceID query must be a number")
		if err != nil {
			l.Warn("cannot make response", util.SlErr(err))
		}
		return
	}

	newPlaylist, err := h.ss.ModifyPlaylist(sourceID, r.Body)
	if err != nil {
		err := api.WriteError(w, http.StatusBadRequest, err.Error())
		if err != nil {
			l.Warn("cannot make response", util.SlErr(err))
		}
		return
	}

	err = api.WriteM3U8(w, http.StatusOK, newPlaylist)
	if err != nil {
		l.Warn("cannot make response", util.SlErr(err))
	}
}

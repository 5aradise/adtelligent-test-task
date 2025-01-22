package stitching

import (
	"log/slog"
	"net/http"
)

type stitchingService interface {
	ModifyPlaylist(playlist []byte) []byte
}

type handler struct {
	s stitchingService
	l *slog.Logger
}

func New(s stitchingService, logger *slog.Logger) *handler {
	return &handler{s, logger}
}

// /stitching.m3u8?sourceID={source_id} 
func (h *handler) Init(router *http.ServeMux) {
	router.HandleFunc("POST /stitching.m3u8", h.handleStitching)
}

func (h *handler) handleStitching(w http.ResponseWriter, r *http.Request) {
}

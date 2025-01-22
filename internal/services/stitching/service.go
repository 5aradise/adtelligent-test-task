package stitching

import (
	"log/slog"
)

type service struct {
	url string
	l   *slog.Logger
}

func New(auctionURL string, logger *slog.Logger) *service {
	return &service{
		url: auctionURL,
		l:   logger,
	}
}

func (s *service) ModifyPlaylist(playlist []byte) []byte {
	return nil
}

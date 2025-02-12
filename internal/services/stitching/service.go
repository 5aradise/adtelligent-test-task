package stitching

import (
	"errors"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/5aradise/adtelligent-test-task/pkg/m3u8"
	"github.com/5aradise/adtelligent-test-task/pkg/util"
)

var (
	ErrAdStartUnfound = errors.New("ad start is unfound in playlist")
	ErrAdEndUnfound   = errors.New("ad end is unfound in playlist")
)

type service struct {
	url    string
	l      *slog.Logger
	client *http.Client
}

func New(auctionURL string, logger *slog.Logger, requestTimeout time.Duration) *service {
	return &service{
		url:    auctionURL,
		l:      logger,
		client: &http.Client{Timeout: requestTimeout},
	}
}

func (s *service) ModifyPlaylist(sourceId int, playlist io.Reader) ([]byte, error) {
	const op = "service.ModifyPlaylist"
	l := s.l.With(
		slog.String("op", op),
		slog.Int("source_id", sourceId),
	)

	var (
		updatedPlaylist = make([]byte, 0, 1024)
		playlistLines   = m3u8.NewLines(playlist)
	)

	updatedPlaylist, adStart := m3u8.SearchAdStart(playlistLines, updatedPlaylist)
	if adStart == nil {
		l.Info("ad start is unfound in playlist", slog.String("playlist", string(updatedPlaylist)))
		return nil, util.OpWrap(op, ErrAdStartUnfound)
	}
	adDuration, err := m3u8.ExtractAdDuration(adStart)
	if err != nil {
		l.Info("failed to extract ad duration", util.SlErr(err), slog.Any("ad_line", adStart))
		return nil, util.OpWrap(op, err)
	}

	adPlaylist, err := s.getAdPlaylist(sourceId, adDuration)
	if err != nil {
		l.Warn("failed to retrieve ad playlist", util.SlErr(err), slog.Duration("ad_duration", adDuration))
		return nil, util.OpWrap(op, err)
	}
	updatedPlaylist = m3u8.AppendLine(updatedPlaylist, adPlaylist)

	// skip replaced part and search for end of ad
	_, adEnd := m3u8.SearchAdEnd(playlistLines, nil)
	if adEnd == nil {
		l.Info("ad end is unfound in playlist", slog.String("playlist", string(updatedPlaylist)))
		return nil, util.OpWrap(op, ErrAdEndUnfound)
	}
	updatedPlaylist = m3u8.AppendLine(updatedPlaylist, adEnd)

	// add rest
	updatedPlaylist = m3u8.AppendLines(updatedPlaylist, playlistLines)

	return updatedPlaylist, nil
}

package stitching

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/5aradise/adtelligent-test-task/pkg/util"
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

	var updatedPlaylist []byte
	var startAdLine []byte
	var isStartFound, isEndFound bool

	scanner := bufio.NewScanner(playlist)

	// search for start of ad
	for scanner.Scan() {
		line := scanner.Bytes()

		if bytes.HasPrefix(line, []byte("#EXT-X-CUE-OUT:")) {
			startAdLine = bytes.TrimRight(line, " ")
			isStartFound = true
			break
		}

		updatedPlaylist = appendLine(updatedPlaylist, line)
	}
	if !isStartFound {
		l.Debug("ad start is unfound in playlist", slog.String("playlist", string(updatedPlaylist)))
		return updatedPlaylist, nil
	}

	// insert ad
	adDurationInSec, err := strconv.ParseFloat(string(startAdLine[15:]), 64)
	if err != nil {
		return nil, err
	}
	maxDuration := time.Duration(adDurationInSec) * time.Second
	adPlaylist, err := s.getAdPlaylist(sourceId, maxDuration)
	if err != nil {
		l.Debug("failed to retrieve ad playlist", util.SlErr(err), slog.Duration("maxDuration", maxDuration))
		return nil, err
	}
	updatedPlaylist = appendLine(updatedPlaylist, adPlaylist)

	// skip replaced part and search for end of ad
	for scanner.Scan() {
		line := scanner.Bytes()

		if bytes.Equal(line, []byte("#EXT-X-CUE-IN")) {
			updatedPlaylist = appendLine(updatedPlaylist, line)
			isEndFound = true
			break
		}
	}
	if !isEndFound {
		l.Debug("ad end is unfound in playlist", slog.String("playlist", string(updatedPlaylist)))
		return nil, err
	}

	// add rest
	for scanner.Scan() {
		line := scanner.Bytes()

		updatedPlaylist = appendLine(updatedPlaylist, line)
	}

	return updatedPlaylist, nil
}

func appendLine(sl, line []byte) []byte {
	return append(sl, append(line, byte('\n'))...)
}

type AdPlaylistResponse struct {
	DurationInMs int    `json:"duracion_in_ms"`
	AdPlaylist   []byte `json:"hls_playlist"`
}

func (s *service) getAdPlaylist(sourceId int, maxDuration time.Duration) ([]byte, error) {
	const op = "service.GetProfitCreative"

	resp, err := s.client.Get(s.url + fmt.Sprintf("?sourceID=%d&maxDuration=%s", sourceId, maxDuration.String()))
	if err != nil {
		return nil, util.OpWrap(op, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errorData, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, util.OpWrap(op, err)
		}
		return nil, util.OpWrap(op, errors.New(string(errorData)))
	}

	var structResp AdPlaylistResponse
	err = json.NewDecoder(resp.Body).Decode(&structResp)
	if err != nil {
		return nil, util.OpWrap(op, err)
	}

	header := []byte("#EXT-X-CUE-OUT:" +
		strconv.Itoa(structResp.DurationInMs/1000) + "." + fmt.Sprintf("%.3d", structResp.DurationInMs%1000) +
		"\n")

	return append(header, structResp.AdPlaylist...), nil
}

package stitching

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/5aradise/adtelligent-test-task/pkg/m3u8"
	"github.com/5aradise/adtelligent-test-task/pkg/ops"
)

type AdPlaylistResponse struct {
	DurationInMs int    `json:"duracion_in_ms"`
	AdPlaylist   []byte `json:"hls_playlist"`
}

func (s *service) getAdPlaylist(sourceId int, maxDuration time.Duration) ([]byte, error) {
	const op = "service.getAdPlaylist"

	resp, err := s.client.Get(s.url + fmt.Sprintf("?sourceID=%d&maxDuration=%s", sourceId, maxDuration.String()))
	if err != nil {
		return nil, ops.Wrap(op, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errorData, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, ops.Wrap(op, err)
		}
		return nil, ops.Wrap(op, errors.New(string(errorData)))
	}

	var structResp AdPlaylistResponse
	err = json.NewDecoder(resp.Body).Decode(&structResp)
	if err != nil {
		return nil, ops.Wrap(op, err)
	}

	return append(m3u8.AdHeader(structResp.DurationInMs), structResp.AdPlaylist...), nil
}

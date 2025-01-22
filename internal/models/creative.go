package models

import "time"

type Creative struct {
	ID              int           `json:"id"`
	Price           Price         `json:"price"`
	Duration        time.Duration `json:"duracion"`
	HlsPlaylistPath string        `json:"hls_playlist"`
}

func (c Creative) Id() int {
	return c.ID
}

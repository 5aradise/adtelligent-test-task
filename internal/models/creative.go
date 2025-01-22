package models

import "time"

type Creative struct {
	ID              int           `json:"id"`
	CampaignID      int           `json:"campaign_id"`
	Price           Price         `json:"price"`
	Duration        time.Duration `json:"duracion"`
	HlsPlaylistPath string        `json:"hls_playlist"`
}

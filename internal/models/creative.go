package models

type Creative struct {
	ID           int    `json:"id"`
	CampaignID   int    `json:"campaign_id"`
	Price        Price  `json:"price"`
	DurationInMs int    `json:"duracion_in_ms"`
	HlsPlaylist  []byte `json:"hls_playlist"`
}

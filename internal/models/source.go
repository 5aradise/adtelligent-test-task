package models

type Source struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	IsActive    bool   `json:"is_active"`
	CampaignIds []int  `json:"campaign_ids,omitempty"`
}

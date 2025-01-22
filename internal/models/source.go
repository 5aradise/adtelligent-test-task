package models

type Source struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	CampaignIds []int  `json:"campaign_ids,omitempty"`
}

func (s Source) Id() int {
	return s.ID
}

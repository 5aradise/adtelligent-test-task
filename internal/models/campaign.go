package models

import "time"

type Campaign struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	CreativeIds []int     `json:"creative_ids,omitempty"`
}

func (c Campaign) Id() int {
	return c.ID
}

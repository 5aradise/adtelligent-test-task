package models

import "time"

type Campaign struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
}

func (c Campaign) IsActive(relativeTime time.Time) bool {
	return c.StartTime.Before(relativeTime) && c.EndTime.After(relativeTime)
}

package models

import "time"

type Quote struct {
	ID        int64     `json:"id"`
	Author    string    `json:"author"`
	Text      string    `json:"quote"`
	CreatedAt time.Time `json:"created_at"`
}

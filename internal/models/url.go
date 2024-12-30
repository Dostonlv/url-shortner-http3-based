package models

import "time"

type URL struct {
	ID          int64     `json:"id" db:"id"`
	ShortCode   string    `json:"short_code" db:"short_code"`
	OriginalURL string    `json:"original_url" db:"original_url"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	Clicks      int       `json:"clicks" db:"clicks"`
}

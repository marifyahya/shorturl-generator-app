package model

import "time"

type URL struct {
	ID          int       `json:"id" db:"id"`
	ShortCode   string    `json:"short_code" db:"short_code"`
	OriginalURL string    `json:"original_url" db:"original_url"`
	Hits        int       `json:"hits" db:"hits"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

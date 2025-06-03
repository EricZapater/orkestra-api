package searches

import "github.com/google/uuid"

type SearchResult struct {
	Tipus string    `json:"tipus" db:"tipus"`
	ID    uuid.UUID `json:"id" db:"id"`
	Text  string    `json:"text" db:"text"`
	Link  string    `json:"link" db:"link"`
}
package llm

import (
	"time"

	"github.com/google/uuid"
)

type QueryRequest struct {
	Question string `json:"question" binding:"required"`
	UserID   uuid.UUID `json:"user_id"`
}

type QueryResponse struct {
	Answer    string `json:"answer"`
	Query     string `json:"query,omitempty"`
	Data      interface{} `json:"data,omitempty"`
	Timestamp time.Time `json:"timestamp"`
}

type QueryContext struct {
	Tables      []TableInfo `json:"tables"`
	Relationships []Relationship `json:"relationships"`
}

type TableInfo struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Columns     []ColumnInfo `json:"columns"`
}

type ColumnInfo struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Description string `json:"description"`
}

type Relationship struct {
	FromTable  string `json:"from_table"`
	FromColumn string `json:"from_column"`
	ToTable    string `json:"to_table"`
	ToColumn   string `json:"to_column"`
	Type       string `json:"type"` // "one_to_many", "many_to_one", "many_to_many"
}

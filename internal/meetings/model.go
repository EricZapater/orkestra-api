package meetings

import "github.com/google/uuid"

type Topic struct {
	ID              string            `json:"id" db:"id"`
	Title           string            `json:"title" db:"title"`
	MeetingID       string            `json:"meeting_id" db:"meeting_id"`
	CreatedAt       string            `json:"created_at" db:"created_at"`
	TopicAgreements *[]TopicAgreement `json:"topic_agreements" db:"topic_agreements"`
}
type TopicAgreement struct {
	ID             string `json:"id" db:"id"`
	Title          string `json:"title" db:"title"`
	MeetingTopicID string `json:"meeting_topic_id" db:"meeting_topic_id"`
	CreatedBy      string `json:"created_by" db:"created_by"`
	CreatedAt      string `json:"created_at" db:"created_at"`
}
type Participant struct {
	ID        string `json:"id" db:"id"`
	MeetingID string `json:"meeting_id" db:"meeting_id"`
	UserID    string `json:"user_id" db:"user_id"`
}

type Meeting struct {
	ID           uuid.UUID      `json:"id" db:"id"`
	GroupID      uuid.UUID      `json:"group_id" db:"group_id"`
	Title        string         `json:"title" db:"title"`
	Description  string         `json:"description" db:"description"`
	StartTime    string         `json:"start_time" db:"start_time"`
	CreatedBy    string         `json:"created_by" db:"created_by"`
	CreatedAt    string         `json:"created_at" db:"created_at"`
	Participants *[]Participant `json:"participants" db:"participants"`
	Topics       *[]Topic       `json:"topics" db:"topics"`
}

type MeetingSummary struct {
	ID		  uuid.UUID `json:"id" db:"id"`
	GroupName string `json:"group_name" db:"group_name"`
	Title       string `json:"title" db:"title"`
	StartTime   string `json:"start_time" db:"start_time"`
	CreatedBy   string `json:"created_by" db:"created_by"`
	NumTopics   int    `json:"num_topics" db:"num_topics"`
	NumParticipants int `json:"num_participants" db:"num_participants"`
	HasAgreements bool `json:"has_agreements" db:"has_agreements"`
}
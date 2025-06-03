package meetings

type CreateMeetingRequest struct {
	GroupID     string `json:"group_id" binding:"required"`
	Title       string `json:"title" binding:"required"`
	Description string `json:"description" binding:"required"`
	StartTime   string `json:"start_time" binding:"required"`
	CreatedBy   string `json:"created_by" binding:"required"`
}
type UpdateMeetingRequest struct {
	GroupID     string `json:"group_id" binding:"required"`
	Title       string `json:"title" binding:"required"`
	Description string `json:"description" binding:"required"`
	StartTime   string `json:"start_time" binding:"required"`
}
type MeetingParticipantRequest struct {
	MeetingID string `json:"meeting_id" binding:"required"`
	UserID    string `json:"user_id" binding:"required"`
}

type MeetingTopicsRequest struct {
	MeetingID string `json:"meeting_id" binding:"required"`
	Title     string `json:"title" binding:"required"`
}
type MeetingTopicAgreementRequest struct {
	MeetingTopicId string `json:"meeting_topic_id" binding:"required"`
	Title          string `json:"title" binding:"required"`
	CreatedBy      string `json:"created_by" binding:"required"`
}

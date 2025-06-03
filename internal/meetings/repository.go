package meetings

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type MeetingRepository interface {
	Create(ctx context.Context, request *Meeting) (Meeting, error)
	Update(ctx context.Context, request *Meeting) (Meeting, error)
	Delete(ctx context.Context, id uuid.UUID) error
	FindByID(ctx context.Context, id uuid.UUID) (Meeting, error)
	FindByGroupID(ctx context.Context, groupID uuid.UUID) ([]MeetingSummary, error)
	FindBetweenDates(ctx context.Context, startDate, endDate time.Time, groupID uuid.UUID) ([]MeetingSummary, error)
	AddParticipant(ctx context.Context, id uuid.UUID, request MeetingParticipantRequest) error
	RemoveParticipant(ctx context.Context, id uuid.UUID) error
	AddTopics(ctx context.Context, id uuid.UUID, request MeetingTopicsRequest) (Topic, error)
	RemoveTopics(ctx context.Context, id uuid.UUID) error
	AddTopicAgreements(ctx context.Context, id uuid.UUID, request MeetingTopicAgreementRequest) (TopicAgreement, error)
	UpdateTopicAgreements(ctx context.Context, id uuid.UUID, request MeetingTopicAgreementRequest) (TopicAgreement, error)
	RemoveTopicAgreements(ctx context.Context, id uuid.UUID, ) error
}

type meetingRepository struct {
	db *sql.DB
}

func NewMeetingRepository(db *sql.DB) MeetingRepository {
	return &meetingRepository{db}
}

func (r *meetingRepository) Create(ctx context.Context, request *Meeting) (Meeting, error) {
	_, err := r.db.ExecContext(ctx, `
	INSERT INTO meetings (id, group_id, created_at, start_time, created_by, title, description)
	VALUES ($1, $2, NOW(), $3, $4, $5, $6)`,
		request.ID, request.GroupID, request.StartTime, request.CreatedBy, request.Title, request.Description,
	)
	if err != nil {
		return Meeting{}, fmt.Errorf("error inserting meeting: %w", err)
	}
	return *request, nil
}

func (r *meetingRepository) Update(ctx context.Context, request *Meeting) (Meeting, error) {
	_, err := r.db.ExecContext(ctx, `
	UPDATE meetings
	SET title = $1, description = $2, start_time = $3
	WHERE id = $4`,
		request.Title, request.Description, request.StartTime, request.ID,
	)
	if err != nil {
		return Meeting{}, fmt.Errorf("error updating meeting: %w", err)
	}
	return *request, nil
}

func (r *meetingRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `
	DELETE FROM meetings
	WHERE id = $1`,
		id,
	)
	if err != nil {
		return fmt.Errorf("error deleting meeting: %w", err)
	}
	return nil
}

func (r *meetingRepository) FindByID(ctx context.Context, id uuid.UUID) (Meeting, error) {
	query := `
		SELECT
			m.id, m.group_id, m.title, m.description, m.start_time, m.created_by, m.created_at,
			p.id, p.user_id,
			t.id, t.title, t.created_at,
			a.id, a.title, a.created_by, a.created_at
		FROM meetings m
		LEFT JOIN meeting_participants p ON m.id = p.meeting_id
		LEFT JOIN meeting_topics t ON m.id = t.meeting_id
		LEFT JOIN meeting_topic_agreements a ON t.id = a.meeting_topic_id
		WHERE m.id = $1
	`

	rows, err := r.db.QueryContext(ctx,query, id)
	if err != nil {
		return Meeting{}, err
	}
	defer rows.Close()

	meeting := &Meeting{
		Participants: &[]Participant{},
		Topics:       &[]Topic{},
	}
	participantsMap := map[string]struct{}{}
	topicsMap := map[string]*Topic{}

	for rows.Next() {
		var (
			meetingID, groupID, title, description, startTime, createdBy, createdAt string
			participantID, participantUserID                                        sql.NullString
			topicID, topicTitle, topicCreatedAt                                    sql.NullString
			agreementID, agreementTitle, agreementCreatedBy, agreementCreatedAt    sql.NullString
		)

		err := rows.Scan(
			&meetingID, &groupID, &title, &description, &startTime, &createdBy, &createdAt,
			&participantID, &participantUserID,
			&topicID, &topicTitle, &topicCreatedAt,
			&agreementID, &agreementTitle, &agreementCreatedBy, &agreementCreatedAt,
		)
		if err != nil {
			return Meeting{}, err
		}

		// Només assignem una vegada
		if meeting.ID == uuid.Nil {
			meeting.ID = uuid.MustParse(meetingID)
			meeting.GroupID = uuid.MustParse(groupID)
			meeting.Title = title
			meeting.Description = description
			meeting.StartTime = startTime
			meeting.CreatedBy = createdBy
			meeting.CreatedAt = createdAt
		}

		// Participants
		if participantID.Valid && participantUserID.Valid {
			if _, exists := participantsMap[participantID.String]; !exists {
				*meeting.Participants = append(*meeting.Participants, Participant{
					ID:        participantID.String,
					UserID:    participantUserID.String,
					MeetingID: meetingID,
				})
				participantsMap[participantID.String] = struct{}{}
			}
		}

		// Temes
		var currentTopic *Topic
		if topicID.Valid {
			topic, exists := topicsMap[topicID.String]
			if !exists {
				topic = &Topic{
					ID:              topicID.String,
					Title:           topicTitle.String,
					MeetingID:       meetingID,
					CreatedAt:       topicCreatedAt.String,
					TopicAgreements: &[]TopicAgreement{},
				}
				topicsMap[topicID.String] = topic
			}
			currentTopic = topic
		}

		// Acords
		if agreementID.Valid && currentTopic != nil {
			*currentTopic.TopicAgreements = append(*currentTopic.TopicAgreements, TopicAgreement{
				ID:             agreementID.String,
				Title:          agreementTitle.String,
				CreatedBy:      agreementCreatedBy.String,
				CreatedAt:      agreementCreatedAt.String,
				MeetingTopicID: topicID.String,
			})
		}
	}

	if err := rows.Err(); err != nil {
		return Meeting{}, err
	}

	// Finalment, afegim els topics
	for _, t := range topicsMap {
		*meeting.Topics = append(*meeting.Topics, *t)
	}

	return *meeting, nil
}

func (r *meetingRepository) FindByGroupID(ctx context.Context, id uuid.UUID) ([]MeetingSummary, error) {
	query := `
		SELECT
			m.id, m.group_id, m.title, m.description, m.start_time, m.created_by, m.created_at,
			COUNT(DISTINCT p.id) AS num_participants,
			COUNT(DISTINCT t.id) AS num_topics,
			COUNT(DISTINCT a.id) > 0 AS has_agreements
		FROM meetings m		
		LEFT JOIN meeting_participants p ON m.id = p.meeting_id
		LEFT JOIN meeting_topics t ON m.id = t.meeting_id
		LEFT JOIN meeting_topic_agreements a ON t.id = a.meeting_topic_id
		WHERE m.group_id = $1
		GROUP BY m.id, m.group_id, m.title, m.description, m.start_time, m.created_by, m.created_at
	`
	rows, err := r.db.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var meetings []MeetingSummary
	for rows.Next() {
		var meeting MeetingSummary
		if err := rows.Scan(
			&meeting.ID, &meeting.GroupName, &meeting.Title, &meeting.StartTime,
			&meeting.CreatedBy, &meeting.NumTopics, &meeting.NumParticipants,
			&meeting.HasAgreements,
		); err != nil {
			return nil, err
		}
		meetings = append(meetings, meeting)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return meetings, nil
}

func (r *meetingRepository) FindBetweenDates(ctx context.Context, startDate, endDate time.Time, groupId uuid.UUID) ([]MeetingSummary, error) {
	query := `
		SELECT
			m.id, 
			g.name as group_name, 
			m.title, 			
			m.start_time, 
			m.created_by, 			
			COUNT(DISTINCT p.id) AS num_participants,
			COUNT(DISTINCT t.id) AS num_topics,
			COUNT(DISTINCT a.id) > 0 AS has_agreements
		FROM meetings m		
			INNER JOIN groups g ON m.group_id = g.id
		LEFT JOIN meeting_participants p ON m.id = p.meeting_id
		LEFT JOIN meeting_topics t ON m.id = t.meeting_id
		LEFT JOIN meeting_topic_agreements a ON t.id = a.meeting_topic_id
		WHERE m.start_time BETWEEN $1 AND $2 
		GROUP BY m.id, g.name, m.title, m.start_time, m.created_by
	`
	rows, err := r.db.Query(query, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var meetings []MeetingSummary
	for rows.Next() {
		var meeting MeetingSummary
		if err := rows.Scan(
			&meeting.ID, &meeting.GroupName, &meeting.Title, &meeting.StartTime,
			&meeting.CreatedBy, &meeting.NumParticipants, &meeting.NumTopics,
			&meeting.HasAgreements,
		); err != nil {
			return nil, err
		}
		meetings = append(meetings, meeting)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return meetings, nil
}

func (r *meetingRepository) AddParticipant(ctx context.Context, id uuid.UUID, request MeetingParticipantRequest) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO meeting_participants (id, meeting_id, user_id)
		VALUES ($1, $2, $3)`,
		id, request.MeetingID, request.UserID,
	)
	if err != nil {
		return fmt.Errorf("error adding participant: %w", err)
	}
	return nil
}

func (r *meetingRepository) RemoveParticipant(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `
		DELETE FROM meeting_participants
		WHERE id = $1`,
		id,
	)
	if err != nil {
		return fmt.Errorf("error removing participant: %w", err)
	}
	return nil
}

func (r *meetingRepository) AddTopics(ctx context.Context, id uuid.UUID, request MeetingTopicsRequest) (Topic, error) {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO meeting_topics (id, meeting_id, title)
		VALUES ($1, $2, $3)`,
		id, request.MeetingID, request.Title,
	)
	if err != nil {
		return Topic{}, fmt.Errorf("error adding topic: %w", err)
	}
	return Topic{Title: request.Title}, nil
}

func (r *meetingRepository) RemoveTopics(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `
		DELETE FROM meeting_topics
		WHERE id = $1`,
		id,
	)
	if err != nil {
		return fmt.Errorf("error removing topic: %w", err)
	}
	return nil
}

func (r *meetingRepository) AddTopicAgreements(ctx context.Context, id uuid.UUID, request MeetingTopicAgreementRequest) (TopicAgreement, error) {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO meeting_topic_agreements (id, meeting_topic_id, title, created_by)
		VALUES ($1, $2, $3, $4)`,
		id, request.MeetingTopicId, request.Title, request.CreatedBy,
	)
	if err != nil {
		return TopicAgreement{}, fmt.Errorf("error adding topic agreement: %w", err)
	}
	return TopicAgreement{Title: request.Title}, nil
}

func (r *meetingRepository) UpdateTopicAgreements(ctx context.Context, id uuid.UUID, request MeetingTopicAgreementRequest) (TopicAgreement, error) {
	_, err := r.db.ExecContext(ctx, `
		UPDATE meeting_topic_agreements
		SET title = $1, created_by = $2
		WHERE id = $3`,
		request.Title, request.CreatedBy, id,
	)
	if err != nil {
		return TopicAgreement{}, fmt.Errorf("error updating topic agreement: %w", err)
	}
	return TopicAgreement{Title: request.Title}, nil
}
func (r *meetingRepository) RemoveTopicAgreements(ctx context.Context, id uuid.UUID, ) error {
	_, err := r.db.ExecContext(ctx, `
		DELETE FROM meeting_topic_agreements
		WHERE id = $1`,
		id,
	)
	if err != nil {
		return fmt.Errorf("error removing topic agreement: %w", err)
	}
	return nil
}
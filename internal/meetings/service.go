package meetings

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type MeetingService interface {
	Create(ctx context.Context, request CreateMeetingRequest) (Meeting, error)
	Update(ctx context.Context, id string, request UpdateMeetingRequest) (Meeting, error)
	Delete(ctx context.Context, id string) error
	FindByID(ctx context.Context, id string) (Meeting, error)
	FindByGroupID(ctx context.Context, groupID string) ([]MeetingSummary, error)
	FindBetweenDates(ctx context.Context, startDate, endDate string, groupID string) ([]MeetingSummary, error)
	AddParticipant(ctx context.Context, request MeetingParticipantRequest) error
	RemoveParticipant(ctx context.Context, id string) error
	AddTopics(ctx context.Context, request MeetingTopicsRequest) (Topic, error)
	RemoveTopics(ctx context.Context, id string) error
	AddTopicAgreements(ctx context.Context, request MeetingTopicAgreementRequest) (TopicAgreement, error)
	UpdateTopicAgreements(ctx context.Context, id string, request MeetingTopicAgreementRequest) (TopicAgreement, error)
	RemoveTopicAgreements(ctx context.Context, id string) error
}

type meetingService struct {
	meetingRepository MeetingRepository
}

func NewMeetingService(meetingRepository MeetingRepository) MeetingService {
	return &meetingService{
		meetingRepository: meetingRepository,
	}
}
func (s *meetingService) Create(ctx context.Context, request CreateMeetingRequest) (Meeting, error) {
	meeting := Meeting{
		ID:          uuid.New(),
		Title:       request.Title,
		Description: request.Description,
		StartTime:   request.StartTime,		
		GroupID:     uuid.MustParse(request.GroupID),
		CreatedBy:   request.CreatedBy,
		CreatedAt:   time.Now().String(),
	}
	return s.meetingRepository.Create(ctx, &meeting)
}

func (s *meetingService) Update(ctx context.Context, id string, request UpdateMeetingRequest) (Meeting, error) {
	meetingID, err := uuid.Parse(id)
	if err != nil {
		return Meeting{}, ErrInvalidID
	}

	// Validate the request
	if request.Title == "" || request.Description == "" || request.StartTime == "" {
		return Meeting{}, ErrInvalidRequest
	}

	meeting, err := s.meetingRepository.FindByID(ctx, meetingID)
	if err != nil {
		return Meeting{}, err
	}

	meeting.Title = request.Title
	meeting.Description = request.Description
	meeting.StartTime = request.StartTime

	return s.meetingRepository.Update(ctx, &meeting)
}

func (s *meetingService) Delete(ctx context.Context, id string) error {
	meetingID, err := uuid.Parse(id)
	if err != nil {
		return ErrInvalidID
	}

	// Check if the meeting exists
	_, err = s.meetingRepository.FindByID(ctx, meetingID)
	if err != nil {
		return err
	}

	return s.meetingRepository.Delete(ctx, meetingID)
}

func (s *meetingService) FindByID(ctx context.Context, id string) (Meeting, error) {
	meetingID, err := uuid.Parse(id)
	if err != nil {
		return Meeting{}, ErrInvalidID
	}

	meeting, err := s.meetingRepository.FindByID(ctx, meetingID)
	if err != nil {
		return Meeting{}, err
	}

	return meeting, nil
}

func (s *meetingService) FindByGroupID(ctx context.Context, groupID string) ([]MeetingSummary, error) {
	groupUUID, err := uuid.Parse(groupID)
	if err != nil {
		return nil, ErrInvalidID
	}

	meetings, err := s.meetingRepository.FindByGroupID(ctx, groupUUID)
	if err != nil {
		return nil, err
	}

	return meetings, nil
}

func (s *meetingService) FindBetweenDates(ctx context.Context, startDate, endDate string, groupID string) ([]MeetingSummary, error) {
	const layoutISO = "2006-01-02"

	startTime, err := time.Parse(layoutISO, startDate)
	if err != nil {
		return nil, ErrInvalidRequest
	}

	endTime, err := time.Parse(layoutISO, endDate)
	if err != nil {
		return nil, ErrInvalidRequest
	}
	groupUUID, err := uuid.Parse(groupID)
	if err != nil {
		groupUUID = uuid.Nil
	}

	meetings, err := s.meetingRepository.FindBetweenDates(ctx, startTime, endTime, groupUUID)
	if err != nil {
		return nil, err
	}

	return meetings, nil
}	

func (s *meetingService) AddParticipant(ctx context.Context, request MeetingParticipantRequest) error {
	_, err := uuid.Parse(request.MeetingID)
	if err != nil {
		return ErrInvalidID
	}

	_ , err = uuid.Parse(request.UserID)
	if err != nil {
		return ErrInvalidID
	}
	id := uuid.New()

	return s.meetingRepository.AddParticipant(ctx, id, request)
}

func (s *meetingService) RemoveParticipant(ctx context.Context, id string) error {	
	idUUID , err := uuid.Parse(id)
	if err != nil {
		return ErrInvalidID
	}

	return s.meetingRepository.RemoveParticipant(ctx, idUUID)
}

func (s *meetingService) AddTopics(ctx context.Context, request MeetingTopicsRequest) (Topic, error) {
	_, err := uuid.Parse(request.MeetingID)
	if err != nil {
		return Topic{}, ErrInvalidID
	}

	// Validate the request
	if request.Title == "" {
		return Topic{}, ErrInvalidRequest
	}

	id := uuid.New()	

	return s.meetingRepository.AddTopics(ctx, id, request)
}

func (s *meetingService) RemoveTopics(ctx context.Context, id string) error {
	idUUID , err := uuid.Parse(id)
	if err != nil {
		return ErrInvalidID
	}

	return s.meetingRepository.RemoveTopics(ctx, idUUID)
}

func (s *meetingService) AddTopicAgreements(ctx context.Context, request MeetingTopicAgreementRequest) (TopicAgreement, error) {
	_, err := uuid.Parse(request.MeetingTopicId)
	if err != nil {
		return TopicAgreement{}, ErrInvalidID
	}

	// Validate the request
	if request.Title == "" || request.CreatedBy == "" {
		return TopicAgreement{}, ErrInvalidRequest
	}

	id := uuid.New()

	return s.meetingRepository.AddTopicAgreements(ctx, id, request)
}

func (s *meetingService) UpdateTopicAgreements(ctx context.Context, id string, request MeetingTopicAgreementRequest) (TopicAgreement, error) {
	idUUID , err := uuid.Parse(id)
	if err != nil {
		return TopicAgreement{}, ErrInvalidID
	}

	// Validate the request
	if request.Title == "" || request.CreatedBy == "" {
		return TopicAgreement{}, ErrInvalidRequest
	}

	return s.meetingRepository.UpdateTopicAgreements(ctx, idUUID, request)
}

func (s *meetingService) RemoveTopicAgreements(ctx context.Context, id string) error {
	idUUID , err := uuid.Parse(id)
	if err != nil {
		return ErrInvalidID
	}

	return s.meetingRepository.RemoveTopicAgreements(ctx, idUUID)
}
package meetings

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type MeetingHandler struct {
	MeetingService MeetingService
}

func NewMeetingHandler(meetingService MeetingService) *MeetingHandler {
	return &MeetingHandler{
		MeetingService: meetingService,
	}
}

func (h *MeetingHandler) CreateMeeting(c *gin.Context) {
	var request CreateMeetingRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	meeting, err := h.MeetingService.Create(c.Request.Context(), request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, meeting)
}

func (h *MeetingHandler) UpdateMeeting(c *gin.Context) {
	id := c.Param("id")
	var request UpdateMeetingRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	meeting, err := h.MeetingService.Update(c.Request.Context(), id, request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, meeting)
}

func (h *MeetingHandler) DeleteMeeting(c *gin.Context) {
	id := c.Param("id")
	err := h.MeetingService.Delete(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func (h *MeetingHandler) GetMeetingByID(c *gin.Context) {
	id := c.Param("id")
	meeting, err := h.MeetingService.FindByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, meeting)
}

func (h *MeetingHandler) GetMeetingsByGroupID(c *gin.Context) {
	id := c.Param("id")
	meetings, err := h.MeetingService.FindByGroupID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, meetings)
}

func (h *MeetingHandler) GetMeetingsBetweenDates(c *gin.Context) {
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	groupID := c.Query("group_id")	
	
	if startDate == "" || endDate == ""  {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Falten paràmetres de consulta"})
		return
	}

	meetings, err := h.MeetingService.FindBetweenDates(c.Request.Context(), startDate, endDate, groupID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, meetings)
}

func (h *MeetingHandler) AddParticipant(c *gin.Context) {
	var request MeetingParticipantRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.MeetingService.AddParticipant(c.Request.Context(), request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, nil)
}

func (h *MeetingHandler) RemoveParticipant(c *gin.Context) {
	id := c.Param("id")
	err := h.MeetingService.RemoveParticipant(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

func (h *MeetingHandler) AddTopics(c *gin.Context) {	
	var request MeetingTopicsRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	topic, err := h.MeetingService.AddTopics(c.Request.Context(), request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, topic)
}

func (h *MeetingHandler) RemoveTopics(c *gin.Context) {
	id := c.Param("id")
	err := h.MeetingService.RemoveTopics(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

func (h *MeetingHandler) AddTopicAgreements(c *gin.Context) {
	var request MeetingTopicAgreementRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	topicAgreement, err := h.MeetingService.AddTopicAgreements(c.Request.Context(), request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, topicAgreement)
}

func (h *MeetingHandler) UpdateTopicAgreements(c *gin.Context) {
	id := c.Param("id")
	var request MeetingTopicAgreementRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	topicAgreement, err := h.MeetingService.UpdateTopicAgreements(c.Request.Context(), id, request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, topicAgreement)
}

func (h *MeetingHandler) RemoveTopicAgreements(c *gin.Context) {
	id := c.Param("id")
	err := h.MeetingService.RemoveTopicAgreements(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
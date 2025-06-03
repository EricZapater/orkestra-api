package meetings

import "github.com/gin-gonic/gin"

func RegisterRoutes(router *gin.RouterGroup, handler *MeetingHandler) {
	router.POST("/meetings", handler.CreateMeeting)
	router.PUT("/meetings/:id", handler.UpdateMeeting)
	router.DELETE("/meetings/:id", handler.DeleteMeeting)
	router.GET("/meetings/:id", handler.GetMeetingByID)
	router.GET("/meetings/group/:id", handler.GetMeetingsByGroupID)
	router.GET("/meetings/dates", handler.GetMeetingsBetweenDates)

	router.POST("/meetings/participants", handler.AddParticipant)
	router.DELETE("/meetings/participants/:id", handler.RemoveParticipant)

	router.POST("/meetings/topics", handler.AddTopics)
	router.DELETE("/meetings/topics/:id", handler.RemoveTopics)

	router.POST("/meetings/topic-agreements", handler.AddTopicAgreements)
	router.PUT("/meetings/topic-agreements/:id", handler.UpdateTopicAgreements)
	router.DELETE("/meetings/topic-agreements/:id", handler.RemoveTopicAgreements)
}

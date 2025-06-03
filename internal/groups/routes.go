package groups

import "github.com/gin-gonic/gin"

func RegisterRoutes(router *gin.RouterGroup, handler *GroupHandler) {
	router.POST("/groups", handler.CreateGroup)
	router.PUT("/groups/:id", handler.UpdateGroup)
	router.DELETE("/groups/:id", handler.DeleteGroup)
	router.GET("/groups/:id", handler.GetGroupByID)
	router.GET("/groups", handler.GetAllGroups)
	router.GET("/groups/user/:userID", handler.GetGroupsByUserID)
	router.PUT("/groups/:id/users/:userID", handler.AddUserToGroup)
	router.DELETE("/groups/:id/users/:userID", handler.RemoveUserFromGroup)
}
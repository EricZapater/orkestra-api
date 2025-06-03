package projects

import "github.com/gin-gonic/gin"

func RegisterRoutes(router *gin.RouterGroup, handler *ProjectHandler) {
	router.POST("/projects", handler.CreateProject)
	router.PUT("/projects/:id", handler.UpdateProject)
	router.DELETE("/projects/:id", handler.DeleteProject)
	router.GET("/projects/:id", handler.GetProjectByID)
	router.GET("/projects", handler.GetAllProjects)
	router.GET("/projects/dates", handler.GetProjectsBetweenDates)
}
package projects

import "github.com/gin-gonic/gin"

func RegisterRoutes(router *gin.RouterGroup, handler *ProjectHandler) {
	router.POST("/projects", handler.CreateProject)
	router.PUT("/projects/:id", handler.UpdateProject)
	router.DELETE("/projects/:id", handler.DeleteProject)
	router.GET("/projects/:id", handler.GetProjectByID)
	router.GET("/projects", handler.GetAllProjects)
	router.GET("/projects/dates", handler.GetProjectsBetweenDates)
	router.GET("/projects/calendar/dates", handler.GetProjectsCalendarBetweenDates)
	router.POST("/projects/operators", handler.AddOperatorToProject)
	router.DELETE("/projects/operators/:id", handler.RemoveOperatorFromProject)
	router.GET("/projects/operators/:project_id", handler.GetOperatorsByProjectID)
	router.GET("/projects/operators/calendar/dates", handler.GetOperatorsCalendarBetweenDates)
	router.POST("/projects/costitems", handler.AddCostItem)
	router.DELETE("/projects/costitems/:id", handler.RemoveCostItem)
	router.GET("/projects/costitems/project/:project_id", handler.GetCostItemByProjectID)
}

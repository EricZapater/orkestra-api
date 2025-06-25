package tasks

import "github.com/gin-gonic/gin"

func RegisterRoutes(router *gin.RouterGroup, handler *TaskHandler){
	router.POST("/tasks", handler.CreateTask)
	router.PUT("/tasks/:id", handler.UpdateTask)
	router.DELETE("/tasks/:id", handler.DeleteTask)
	router.GET("/tasks/:id", handler.GetTaskByID)
	router.GET("/tasks/status/:status", handler.GetTaskByStatus)
	router.GET("/tasks/user/:userid", handler.GetTaskByUserID)
	router.GET("/tasks/project/:projectid", handler.GetTaskByProjectID)
	router.GET("/tasks/priority/:priority", handler.GetTaskByPriority)
	router.GET("/tasks", handler.GetAllTask)
}
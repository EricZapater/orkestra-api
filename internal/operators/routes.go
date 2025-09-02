package operators

import "github.com/gin-gonic/gin"

func RegisterRoutes(router *gin.RouterGroup, handler *OperatorHandler) {
	router.POST("/operators", handler.Create)
	router.PUT("/operators/:id", handler.Update)
	router.DELETE("/operators/:id", handler.Delete)
	router.GET("/operators/:id", handler.FindByID)
	router.GET("/operators", handler.FindAll)
}
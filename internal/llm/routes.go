package llm

import (
	"github.com/gin-gonic/gin"
)

//func RegisterRoutes(router *gin.RouterGroup, handler *OperatorHandler) {
func RegisterRoutes(router *gin.RouterGroup, handler *Handler) {
	llmGroup := router.Group("/llm")
	{
		llmGroup.POST("/query", handler.Query)
		llmGroup.GET("/schema", handler.GetSchema)
	}
}

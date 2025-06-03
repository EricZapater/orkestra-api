package searches

import "github.com/gin-gonic/gin"

func RegisterRoutes(router *gin.RouterGroup, handler *SearchHandler){
	router.POST("/search", handler.GetByText)
}
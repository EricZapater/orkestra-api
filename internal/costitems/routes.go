package costitems

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.RouterGroup, handler *CostItemHandler){
	prefix := "/costitems"
	router.POST(prefix, handler.Create)
	router.PUT(fmt.Sprintf("%s/:id", prefix), handler.Update)
	router.DELETE(fmt.Sprintf("%s/:id", prefix), handler.Delete)
	router.GET(fmt.Sprintf("%s/:id", prefix), handler.GetByID)
	router.GET(fmt.Sprintf("%s/project/:projectid", prefix), handler.GetByProjectID)
	router.GET(prefix, handler.GetAll)
}
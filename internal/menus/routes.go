package menus

import "github.com/gin-gonic/gin"

func RegisterRoutes(router *gin.RouterGroup, handler *MenuHandler) {
	router.POST("/menus", handler.Create)
	router.PUT("/menus/:id", handler.Update)
	router.DELETE("/menus/:id", handler.Delete)
	router.GET("/menus/:id", handler.GetByID)
	router.GET("/menus", handler.GetAll)
	router.GET("/menus/profile/:profile_id", handler.GetByProfileID)

	router.POST("/menus/addprofile", handler.AddMenuToProfile)
	router.POST("/menus/removeprofile", handler.RemoveMenuFromProfile)
}
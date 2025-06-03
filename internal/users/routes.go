package users

import "github.com/gin-gonic/gin"

func RegisterRoutes(router *gin.RouterGroup, handler *UserHandler) {
	roles := router.Group("/users")
	{		
		roles.PUT("/:id", handler.Update)
		roles.DELETE("/:id", handler.Delete)
		roles.POST("/change-password", handler.ChangePassword)				
		roles.GET("/username/:username", handler.GetByUsername)
		roles.GET("/phone/:phone_number", handler.GetByPhoneNumber)
		roles.GET("/:id", handler.GetByID)
		roles.GET("", handler.GetAll)
		roles.GET("/group/:group_id", handler.GetByGroupID)
	}
}

func RegisterPublicRoutes(router *gin.RouterGroup, handler *UserHandler) {
    router.POST("/register", handler.Create) // Ruta pública per crear usuaris
	router.POST("/verify", handler.VerifyUser)
}
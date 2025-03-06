package http

import (
	"auth-service/internal/http/handler"
	"github.com/gin-gonic/gin"
)

func SetupRouter(controller *handler.Controller) *gin.Engine {
	r := gin.Default()

	api := r.Group("/api/v1/auth")
	{
		api.GET("/health", controller.Health)

		api.POST("/login", controller.Login)
		api.POST("/consent", controller.Consent)
	}

	return r
}

package http

import (
	"auth-service/internal/http/handler"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter(controller *handler.Controller) *gin.Engine {
	r := gin.Default()
	r.Use(cors.Default())
	r.Any("/check", controller.Health)

	api := r.Group("/api/auth")
	{
		api.GET("/health", controller.Health)

		api.POST("/login", controller.Login)
		api.POST("/consent", controller.Consent)
	}

	return r
}

package httpengine

import (
	"github.com/gin-gonic/gin"
	"ticket-service/internal/httpengine/handler"
)

func NewRouter(handler handler.Handler) *gin.Engine {
	router := gin.New()

	router.Use(gin.Logger())

	router.POST("/api/ticket", handler.CreateTicket)
	router.PUT("/api/ticket/:uuid/status", handler.UpdateTicketStatus)
	router.PUT("/api/ticket/:uuid/responsible", handler.AssignTicketResponsible)
	router.GET("/api/ticket", handler.GetTickets)
	router.GET("/check", handler.Check)

	return router
}

package http

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"notification-service/internal/usecase"
)

func RegisterMailReceiver(useCase *usecase.RegisterMailReceiverUseCase) func(c *gin.Context) {
	return func(c *gin.Context) {
		var registerMailReceiverDto RegisterMailReceiverDto

		if err := c.BindJSON(&registerMailReceiverDto); err != nil {
			c.JSON(http.StatusBadRequest, Error{Message: err.Error()})
			return
		}

		if err := useCase.Execute(
			c.Request.Context(),
			registerMailReceiverDto.UserUUID,
			registerMailReceiverDto.Mail,
		); err != nil {
			c.JSON(http.StatusInternalServerError, Error{Message: err.Error()})
			return
		}

		c.Status(http.StatusOK)
	}
}

func HealthCheck() func(c *gin.Context) {
	return func(c *gin.Context) {
		c.Status(http.StatusOK)
	}
}

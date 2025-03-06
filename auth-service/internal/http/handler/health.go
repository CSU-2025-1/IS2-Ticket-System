package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (c *Controller) Health(ctx *gin.Context) {
	ctx.Status(http.StatusOK)
}

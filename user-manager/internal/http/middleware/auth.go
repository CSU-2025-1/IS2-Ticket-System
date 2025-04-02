package middleware

import (
	"github.com/gin-gonic/gin"
	"user-mananger/internal/domain/constant"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user_uuid := ctx.GetHeader(constant.HeaderUserUUID)
		ctx.Set(constant.CtxUserUUID, user_uuid)
	}
}

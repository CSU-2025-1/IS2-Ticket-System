package handler

import (
	"auth-service/internal/domain/constant"
	"auth-service/internal/domain/errors/service"
	"auth-service/internal/http/dto"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
)

func (c *Controller) Login(ctx *gin.Context) {
	var req dto.AuthenticateRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.NewApiError("не удалось спарсить тело запроса"))
		return
	}

	redirectURL, err := c.Service.Auth.Authenticate(ctx, req.Challenge, req.Login, req.Password)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			ctx.Redirect(http.StatusMovedPermanently, fmt.Sprintf(constant.Oauth2LoginChallenge, req.Challenge, err.Error()))
			return
		}

		slog.Error(ctx.FullPath(), slog.String(constant.ErrorField, err.Error()))
		ctx.Status(http.StatusInternalServerError)
		return
	}

	slog.Info("login redirect", slog.String("url", redirectURL))
	ctx.Redirect(http.StatusMovedPermanently, redirectURL)
}

func (c *Controller) Consent(ctx *gin.Context) {
	var req dto.ConsentRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.NewApiError("не удалось спарсить тело запроса"))
		return
	}

	redirectURL, err := c.Service.Auth.Consent(ctx, req.Challenge, []string{"offline_access"})
	if err != nil {
		slog.Error(ctx.FullPath(), slog.String(constant.ErrorField, err.Error()))
		ctx.Status(http.StatusInternalServerError)
		return
	}

	slog.Info("consent redirect", slog.String("url", redirectURL))
	ctx.Redirect(http.StatusMovedPermanently, redirectURL)
}

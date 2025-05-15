package handler

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log/slog"
	"net/http"
	"user-mananger/internal/domain/entity"
	"user-mananger/internal/http/dto"
)

type UserSaver interface {
	Save(ctx context.Context, authData entity.User) error
}

type UserController interface {
	CreateUser(ctx context.Context, user entity.User) (uuid.UUID, error)
	DeleteUser(ctx context.Context, uuid uuid.UUID) error
	GetUsers(ctx context.Context) ([]entity.User, error)
}

type UsersHandler struct {
	saver UserSaver
	user  UserController
}

func NewUsersHandler(user UserController, saver UserSaver) *UsersHandler {
	return &UsersHandler{
		saver: saver,
		user:  user,
	}
}

func (u *UsersHandler) CreateUser(c *gin.Context) {
	var req dto.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Error("create user parse", slog.String("error", err.Error()))
		c.JSON(http.StatusBadRequest, gin.H{"error": "чет муть какая-то"})
		return
	}

	userUUID, err := u.user.CreateUser(c, entity.User{
		Login:    req.Login,
		Password: req.Password,
	})
	if err != nil {
		slog.Error("create user", slog.String("error", err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "чет не-то, я сломался"})
		return
	}

	err = u.saver.Save(c, entity.User{
		UUID:     userUUID,
		Login:    req.Login,
		Password: req.Password,
	})
	if err != nil {
		slog.Error("save auth data", slog.String("error", err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "чет не-то, я сломался"})
		return
	}

	c.JSON(http.StatusCreated, dto.CreateUserResponse{
		UUID: userUUID.String(),
	})
}

func (u *UsersHandler) UpdateUser(c *gin.Context) {
	c.Status(http.StatusNotImplemented)
}

func (u *UsersHandler) DeleteUser(c *gin.Context) {
	user := c.Param("user_uuid")

	userUUID, err := uuid.Parse(user)
	if err != nil {
		slog.Warn("incorrect user_uuid", slog.String("error", err.Error()))
		c.AbortWithStatusJSON(http.StatusBadRequest, "че с ююидом, это не ююид, уходи")
	}

	err = u.user.DeleteUser(c, userUUID)
	if err != nil {
		slog.Error("create user", slog.String("error", err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "чет не-то, я сломался"})
		return
	}

	c.Status(http.StatusNoContent)
}

func (u *UsersHandler) GetUsers(c *gin.Context) {
	users, err := u.user.GetUsers(c)
	if err != nil {
		slog.Error("create user", slog.String("error", err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "чет не-то, я сломался"})
		return
	}

	response := make([]dto.User, len(users))
	for i, user := range users {
		response[i] = dto.User{
			UUID:  user.UUID.String(),
			Login: user.Login,
		}
	}

	c.JSON(http.StatusOK, response)
}

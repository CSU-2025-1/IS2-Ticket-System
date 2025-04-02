// internal/http/handler/groups.go
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

type GroupController interface {
	CreateGroup(ctx context.Context, name string) (uuid.UUID, error)
	DeleteGroup(ctx context.Context, id uuid.UUID) error
	GetGroups(ctx context.Context) ([]entity.Group, error)
	AddUsersToGroup(ctx context.Context, groupID uuid.UUID, userIDs []uuid.UUID) error
	RemoveUsersFromGroup(ctx context.Context, groupID uuid.UUID, userIDs []uuid.UUID) error
	GetGroupUsers(ctx context.Context, groupID uuid.UUID) ([]uuid.UUID, error)
}

type GroupsHandler struct {
	group GroupController
}

func NewGroupsHandler(group GroupController) *GroupsHandler {
	return &GroupsHandler{
		group: group,
	}
}

func (g *GroupsHandler) CreateGroup(ctx *gin.Context) {
	var req dto.CreateGroupRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		slog.Error("create group parse", slog.String("error", err.Error()))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	groupUUID, err := g.group.CreateGroup(ctx, req.Name)
	if err != nil {
		slog.Error("create group", slog.String("error", err.Error()))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create group"})
		return
	}

	ctx.JSON(http.StatusCreated, dto.CreateGroupResponse{
		UUID: groupUUID.String(),
	})
}

func (g *GroupsHandler) DeleteGroup(ctx *gin.Context) {
	groupID := ctx.Param("group_id")

	groupUUID, err := uuid.Parse(groupID)
	if err != nil {
		slog.Warn("incorrect group_id", slog.String("error", err.Error()))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid group id"})
		return
	}

	if err := g.group.DeleteGroup(ctx, groupUUID); err != nil {
		slog.Error("delete group", slog.String("error", err.Error()))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete group"})
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (g *GroupsHandler) GetGroups(ctx *gin.Context) {
	groups, err := g.group.GetGroups(ctx)
	if err != nil {
		slog.Error("get groups", slog.String("error", err.Error()))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get groups"})
		return
	}

	response := make([]dto.Group, len(groups))
	for i, group := range groups {
		response[i] = dto.Group{
			UUID: group.UUID.String(),
			Name: group.Name,
		}
	}

	ctx.JSON(http.StatusOK, response)
}

func (g *GroupsHandler) AddUserToGroup(ctx *gin.Context) {
	groupID := ctx.Param("group_id")
	groupUUID, err := uuid.Parse(groupID)
	if err != nil {
		slog.Warn("incorrect group_id", slog.String("error", err.Error()))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid group id"})
		return
	}

	var req dto.AddUsersToGroupRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		slog.Error("add users to group parse", slog.String("error", err.Error()))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	userUUIDs := make([]uuid.UUID, len(req.Users))
	for i, userID := range req.Users {
		userUUID, err := uuid.Parse(userID)
		if err != nil {
			slog.Warn("incorrect user_id", slog.String("error", err.Error()))
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
			return
		}
		userUUIDs[i] = userUUID
	}

	if err := g.group.AddUsersToGroup(ctx, groupUUID, userUUIDs); err != nil {
		slog.Error("add users to group", slog.String("error", err.Error()))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to add users to group"})
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (g *GroupsHandler) RemoveUserFromGroup(ctx *gin.Context) {
	groupID := ctx.Param("group_id")
	groupUUID, err := uuid.Parse(groupID)
	if err != nil {
		slog.Warn("incorrect group_id", slog.String("error", err.Error()))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid group id"})
		return
	}

	var req dto.RemoveUsersFromGroupRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		slog.Error("remove users from group parse", slog.String("error", err.Error()))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	userUUIDs := make([]uuid.UUID, len(req.Users))
	for i, userID := range req.Users {
		userUUID, err := uuid.Parse(userID)
		if err != nil {
			slog.Warn("incorrect user_id", slog.String("error", err.Error()))
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
			return
		}
		userUUIDs[i] = userUUID
	}

	if err := g.group.RemoveUsersFromGroup(ctx, groupUUID, userUUIDs); err != nil {
		slog.Error("remove users from group", slog.String("error", err.Error()))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to remove users from group"})
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (g *GroupsHandler) ListUsers(ctx *gin.Context) {
	groupID := ctx.Param("group_id")
	groupUUID, err := uuid.Parse(groupID)
	if err != nil {
		slog.Warn("incorrect group_id", slog.String("error", err.Error()))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid group id"})
		return
	}

	userUUIDs, err := g.group.GetGroupUsers(ctx, groupUUID)
	if err != nil {
		slog.Error("get group users", slog.String("error", err.Error()))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get group users"})
		return
	}

	response := make([]string, len(userUUIDs))
	for i, userUUID := range userUUIDs {
		response[i] = userUUID.String()
	}

	ctx.JSON(http.StatusOK, gin.H{"users": response})
}

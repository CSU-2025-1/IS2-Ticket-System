package handler

import "github.com/gin-gonic/gin"

type GroupsHandler struct {
}

func NewGroupsHandler() *GroupsHandler {
	return &GroupsHandler{}
}

func (g *GroupsHandler) CreateGroup(ctx *gin.Context) {

}

func (g *GroupsHandler) DeleteGroup(ctx *gin.Context) {

}

func (g *GroupsHandler) GetGroups(ctx *gin.Context) {

}

func (g *GroupsHandler) AddUserToGroup(ctx *gin.Context) {

}

func (g *GroupsHandler) RemoveUserFromGroup(ctx *gin.Context) {

}

func (g *GroupsHandler) ListUsers(ctx *gin.Context) {

}

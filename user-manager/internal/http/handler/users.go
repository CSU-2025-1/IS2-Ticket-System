package handler

import "github.com/gin-gonic/gin"

type UsersHandler struct {
}

func NewUsersHandler() *UsersHandler {
	return &UsersHandler{}
}

func (u *UsersHandler) CreateUser(c *gin.Context) {

}

func (u *UsersHandler) UpdateUser(c *gin.Context) {

}

func (u *UsersHandler) DeleteUser(c *gin.Context) {

}

func (u *UsersHandler) GetUsers(c *gin.Context) {

}

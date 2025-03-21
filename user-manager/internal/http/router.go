package http

import (
	"github.com/gin-gonic/gin"
	"user-mananger/internal/http/handler"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	health := handler.NewHealthHandler()
	r.Any("/health", health.Health)

	api := r.Group("/api/v1/users")

	usersHandler := handler.NewUsersHandler()
	users := api.Group("")
	{
		users.GET("", usersHandler.GetUsers)
		users.POST("", usersHandler.CreateUser)
		users.PUT("", usersHandler.UpdateUser)
		users.DELETE("", usersHandler.DeleteUser)
	}

	groupsHandler := handler.NewGroupsHandler()
	groups := api.Group("/groups")
	{
		groups.GET("", groupsHandler.GetGroups)
		groups.POST("", groupsHandler.CreateGroup)
		groups.DELETE("", groupsHandler.DeleteGroup)

		groups.POST("/:group_id/add-user", groupsHandler.AddUserToGroup)
		groups.DELETE("/:group_id/remove-user", groupsHandler.RemoveUserFromGroup)

		groups.GET("/:group_id/list", groupsHandler.ListUsers)

	}

	return r
}

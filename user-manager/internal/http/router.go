package http

import (
	"github.com/gin-gonic/gin"
	"user-mananger/internal/http/handler"
	"user-mananger/internal/http/middleware"
	"user-mananger/internal/repository"
)

func SetupRouter(repository repository.Manager) *gin.Engine {
	r := gin.Default()

	health := handler.NewHealthHandler()
	r.Any("/health", health.Health)
	r.Any("/check", health.Health)

	api := r.Group("/api/users")
	api.Use(middleware.AuthMiddleware())

	usersHandler := handler.NewUsersHandler(repository.User, repository.AuthData)
	users := api.Group("")
	{
		users.GET("", usersHandler.GetUsers)
		users.POST("", usersHandler.CreateUser)
		users.PUT("", usersHandler.UpdateUser)
		users.DELETE("/:user_uuid", usersHandler.DeleteUser)
	}

	groupsHandler := handler.NewGroupsHandler(repository.Group)
	groups := api.Group("/groups")
	{
		groups.GET("", groupsHandler.GetGroups)
		groups.POST("", groupsHandler.CreateGroup)
		groups.DELETE("/:group_id", groupsHandler.DeleteGroup)

		groups.POST("/:group_id/add-user", groupsHandler.AddUserToGroup)
		groups.DELETE("/:group_id/remove-user", groupsHandler.RemoveUserFromGroup)

		groups.GET("/:group_id/list", groupsHandler.ListUsers)
	}

	return r
}

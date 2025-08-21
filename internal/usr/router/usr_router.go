package usr_router

import (
	"github.com/gin-gonic/gin"
	"github.com/williamkoller/cloud-architecture-golang/internal/usr/handler"
)

func RegisterUserRoutes(group *gin.RouterGroup, h *handler.UserHandler) {
	users := group.Group("/users")
	{
		users.POST("", h.CreateUser)
		users.GET("", h.ListUsers)
		users.GET("/:email", h.GetUser)
		users.PATCH("/:email", h.UpdateUser)
		users.DELETE("/:email", h.DeleteUser)
	}
}

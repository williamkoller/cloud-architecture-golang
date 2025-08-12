package usr_router

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type User struct {
	Name  string
	Email string
}

type UserResponse struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func mapToUserResponse(u User) UserResponse {
	return UserResponse{
		Name:  u.Name,
		Email: u.Email,
	}
}

func mapUsersToResponse(users []User) []UserResponse {
	out := make([]UserResponse, 0, len(users))
	for _, u := range users {
		out = append(out, mapToUserResponse(u))
	}
	return out
}

func handlerUsers(c *gin.Context) {
	users := []User{
		{Name: "William K", Email: "william@mail.com"},
		{Name: "Novo user test", Email: "novo-user@mail.com"},
	}

	response := mapUsersToResponse(users)
	c.JSON(http.StatusOK, response)
}

func RegisterUserRoutes(r *gin.Engine) {
	r.GET("/users", handlerUsers)
}

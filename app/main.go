package main

import (
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
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

var (
	router      *gin.Engine
	ginLambdaV2 *ginadapter.GinLambdaV2
)

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

func init() {
	gin.SetMode(gin.ReleaseMode)
	router = gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})
	router.GET("/users", handlerUsers)

	ginLambdaV2 = ginadapter.NewV2(router)
}

func handlerUsers(c *gin.Context) {
	users := []User{
		{Name: "William K", Email: "william@mail.com"},
		{Name: "Novo user test", Email: "novo-user@mail.com"},
	}

	response := mapUsersToResponse(users)
	c.JSON(http.StatusOK, response)
}

func main() {
	if os.Getenv("LOCAL") == "true" {
		log.Println("Running locally on :8080")
		if err := router.Run(":8080"); err != nil {
			log.Fatal(err)
		}
		return
	}
	lambda.Start(ginLambdaV2.ProxyWithContext)
}

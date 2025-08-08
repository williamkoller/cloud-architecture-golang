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

var ginLambda *ginadapter.GinLambda
var router *gin.Engine

func mapToUserResponse(u User) UserResponse {
	return UserResponse{
		Name:  u.Name,
		Email: u.Email,
	}
}

func init() {
	router = gin.Default()
	router.GET("/health", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})
	router.GET("/user", handlerUsers)

	ginLambda = ginadapter.New(router)
}

func handlerUsers(c *gin.Context) {
	user := User{
		Name:  "William K",
		Email: "william@mail.com",
	}

	response := mapToUserResponse(user)
	c.JSON(http.StatusOK, response)
}

func main() {
	if os.Getenv("LOCAL") == "true" {
		log.Println("Running locally on :8080")
		router.Run(":8080")
	} else {
		lambda.Start(ginLambda.Proxy)
	}
}

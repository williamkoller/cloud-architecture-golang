package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})
	router.GET("/users", handlerUsers)

	return router
}

func TestHealthEndpoint(t *testing.T) {
	router := setupTestRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/health", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "OK", w.Body.String())
}

func TestUsersEndpoint(t *testing.T) {
	router := setupTestRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "William K")
	assert.Contains(t, w.Body.String(), "william@mail.com")
}

func TestMapToUserResponse(t *testing.T) {
	user := User{
		Name:  "Test User",
		Email: "test@example.com",
	}

	response := mapToUserResponse(user)

	assert.Equal(t, "Test User", response.Name)
	assert.Equal(t, "test@example.com", response.Email)
}

func TestMapUsersToResponse(t *testing.T) {
	users := []User{
		{Name: "User 1", Email: "user1@example.com"},
		{Name: "User 2", Email: "user2@example.com"},
	}

	responses := mapUsersToResponse(users)

	assert.Len(t, responses, 2)
	assert.Equal(t, "User 1", responses[0].Name)
	assert.Equal(t, "user1@example.com", responses[0].Email)
	assert.Equal(t, "User 2", responses[1].Name)
	assert.Equal(t, "user2@example.com", responses[1].Email)
}

package handler

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/williamkoller/cloud-architecture-golang/internal/usr/domain"
	"github.com/williamkoller/cloud-architecture-golang/internal/usr/domain/vo"
	"github.com/williamkoller/cloud-architecture-golang/internal/usr/dtos"
	"github.com/williamkoller/cloud-architecture-golang/internal/usr/mappers"
	"github.com/williamkoller/cloud-architecture-golang/internal/usr/repository"
	"github.com/williamkoller/cloud-architecture-golang/internal/usr/validation"
)

type UserHandler struct {
	repo repository.UserRepository
}

func NewUserHandler(repo repository.UserRepository) *UserHandler {
	return &UserHandler{repo: repo}
}

// Timeoutzinho helper para operações rápidas
func (h *UserHandler) ctx(c *gin.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(c.Request.Context(), 2*time.Second)
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var req dtos.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		validation.RespondValidationError(c, err)
		return
	}

	active := true
	if req.Active != nil {
		active = *req.Active
	}

	u, err := domain.NewUser(
		strings.TrimSpace(req.Name),
		strings.TrimSpace(req.Email),
		req.Password,
		active,
		domain.UserType(strings.TrimSpace(req.UserType)),
	)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := h.ctx(c)
	defer cancel()

	if err := h.repo.Create(ctx, u); err != nil {
		if err == repository.ErrAlreadyExists {
			c.JSON(http.StatusConflict, gin.H{"error": "user already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, mappers.ToUserResponse(u))
}

func (h *UserHandler) ListUsers(c *gin.Context) {
	ctx, cancel := h.ctx(c)
	defer cancel()

	users, err := h.repo.List(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	resp := make([]mappers.UserResponse, 0, len(users))
	for _, u := range users {
		resp = append(resp, mappers.ToUserResponse(u))
	}
	c.JSON(http.StatusOK, resp)
}

func (h *UserHandler) GetUser(c *gin.Context) {
	emailParam := strings.TrimSpace(c.Param("email"))
	email, err := vo.NewEmail(emailParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid email"})
		return
	}

	ctx, cancel := h.ctx(c)
	defer cancel()

	u, ok, err := h.repo.GetByEmail(ctx, email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, mappers.ToUserResponse(u))
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	emailParam := strings.TrimSpace(c.Param("email"))
	email, err := vo.NewEmail(emailParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid email"})
		return
	}

	var req dtos.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		validation.RespondValidationError(c, err)
		return
	}

	ctx, cancel := h.ctx(c)
	defer cancel()

	current, ok, err := h.repo.GetByEmail(ctx, email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	// aplica mudanças parciais
	name := current.Name
	if req.Name != nil {
		name = strings.TrimSpace(*req.Name)
	}
	active := current.Active
	if req.Active != nil {
		active = *req.Active
	}
	userType := current.UserType
	if req.UserType != nil {
		userType = domain.UserType(strings.TrimSpace(*req.UserType))
	}
	password := string(current.Password) // vo.Password é alias de string
	if req.Password != nil {
		password = *req.Password
	}

	updated, err := domain.NewUser(name, string(email), password, active, userType)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	if err := h.repo.Update(ctx, updated); err != nil {
		if err == repository.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, mappers.ToUserResponse(updated))
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	emailParam := strings.TrimSpace(c.Param("email"))
	email, err := vo.NewEmail(emailParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid email"})
		return
	}

	ctx, cancel := h.ctx(c)
	defer cancel()

	if err := h.repo.Delete(ctx, email); err != nil {
		if err == repository.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

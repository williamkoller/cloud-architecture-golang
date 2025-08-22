package handler

import (
	"context"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/williamkoller/cloud-architecture-golang/internal/metrics"
	"github.com/williamkoller/cloud-architecture-golang/internal/usr/domain"
	"github.com/williamkoller/cloud-architecture-golang/internal/usr/domain/vo"
	"github.com/williamkoller/cloud-architecture-golang/internal/usr/dtos"
	"github.com/williamkoller/cloud-architecture-golang/internal/usr/mappers"
	"github.com/williamkoller/cloud-architecture-golang/internal/usr/repository"
	"github.com/williamkoller/cloud-architecture-golang/internal/usr/validation"
)

// CacheItem representa um item no cache com TTL
type CacheItem struct {
	Value     mappers.UserResponse
	ExpiresAt time.Time
}

// IsExpired verifica se o item do cache expirou
func (ci *CacheItem) IsExpired() bool {
	return time.Now().After(ci.ExpiresAt)
}

type UserHandler struct {
	repo  repository.UserRepository
	cache sync.Map // map[string]*CacheItem
	// Configurações de performance
	cacheTTL       time.Duration
	requestTimeout time.Duration
}

func NewUserHandler(repo repository.UserRepository) *UserHandler {
	handler := &UserHandler{
		repo:           repo,
		cacheTTL:       30 * time.Second, // TTL mais curto para consistência
		requestTimeout: 5 * time.Second,  // Timeout mais generoso
	}

	return handler
}

// startCacheCleanup inicia limpeza periódica do cache (removido - pode causar crashes)

// cleanExpiredCache remove itens expirados do cache
func (h *UserHandler) cleanExpiredCache() {
	h.cache.Range(func(key, value interface{}) bool {
		if item, ok := value.(*CacheItem); ok && item.IsExpired() {
			h.cache.Delete(key)
		}
		return true
	})
}

// ctx cria um contexto com timeout otimizado
func (h *UserHandler) ctx(c *gin.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(c.Request.Context(), h.requestTimeout)
}

// getCachedUser busca usuário no cache
func (h *UserHandler) getCachedUser(email string) (mappers.UserResponse, bool) {
	if cached, ok := h.cache.Load(email); ok {
		if item, ok := cached.(*CacheItem); ok {
			if !item.IsExpired() {
				return item.Value, true
			}
			// Remove item expirado de forma segura
			h.cache.Delete(email)
		}
	}
	return mappers.UserResponse{}, false
}

// setCachedUser armazena usuário no cache
func (h *UserHandler) setCachedUser(email string, user mappers.UserResponse) {
	item := &CacheItem{
		Value:     user,
		ExpiresAt: time.Now().Add(h.cacheTTL),
	}
	h.cache.Store(email, item)
}

// invalidateCache remove usuário do cache
func (h *UserHandler) invalidateCache(email string) {
	h.cache.Delete(email)
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var req dtos.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		validation.RespondValidationError(c, err)
		return
	}

	// Validação de entrada mais rigorosa
	if strings.TrimSpace(req.Name) == "" || strings.TrimSpace(req.Email) == "" || req.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name, email and password are required"})
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

	// Atualizar métricas de forma síncrona e eficiente
	metrics.UsersCreatedInc()

	response := mappers.ToUserResponse(u)
	// Cachear o usuário criado
	h.setCachedUser(string(u.Email), response)

	c.JSON(http.StatusCreated, response)
}

func (h *UserHandler) ListUsers(c *gin.Context) {
	ctx, cancel := h.ctx(c)
	defer cancel()

	users, err := h.repo.List(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Processar resposta de forma otimizada com pré-alocação
	resp := make([]mappers.UserResponse, 0, len(users))
	for _, u := range users {
		resp = append(resp, mappers.ToUserResponse(u))
	}

	c.JSON(http.StatusOK, resp)
}

func (h *UserHandler) GetUser(c *gin.Context) {
	emailParam := strings.TrimSpace(c.Param("email"))
	if emailParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email parameter is required"})
		return
	}

	// Verificar cache primeiro
	if cached, found := h.getCachedUser(emailParam); found {
		c.JSON(http.StatusOK, cached)
		return
	}

	email, err := vo.NewEmail(emailParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid email format"})
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

	userResp := mappers.ToUserResponse(u)
	// Cachear resultado
	h.setCachedUser(emailParam, userResp)

	c.JSON(http.StatusOK, userResp)
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	emailParam := strings.TrimSpace(c.Param("email"))
	if emailParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email parameter is required"})
		return
	}

	email, err := vo.NewEmail(emailParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid email format"})
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

	// Aplicar mudanças parciais
	name := current.Name
	if req.Name != nil && strings.TrimSpace(*req.Name) != "" {
		name = strings.TrimSpace(*req.Name)
	}
	active := current.Active
	if req.Active != nil {
		active = *req.Active
	}
	userType := current.UserType
	if req.UserType != nil && strings.TrimSpace(*req.UserType) != "" {
		userType = domain.UserType(strings.TrimSpace(*req.UserType))
	}
	password := string(current.Password)
	if req.Password != nil && *req.Password != "" {
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

	// Invalidar cache e atualizar métricas
	h.invalidateCache(emailParam)
	metrics.UsersUpdatedInc()

	response := mappers.ToUserResponse(updated)
	// Cachear o usuário atualizado
	h.setCachedUser(emailParam, response)

	c.JSON(http.StatusOK, response)
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	emailParam := strings.TrimSpace(c.Param("email"))
	if emailParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email parameter is required"})
		return
	}

	email, err := vo.NewEmail(emailParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid email format"})
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

	// Invalidar cache e atualizar métricas
	h.invalidateCache(emailParam)
	metrics.UsersDeletedInc()

	c.Status(http.StatusNoContent)
}

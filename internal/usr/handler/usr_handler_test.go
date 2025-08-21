package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/williamkoller/cloud-architecture-golang/internal/usr/domain"
	"github.com/williamkoller/cloud-architecture-golang/internal/usr/domain/vo"
	"github.com/williamkoller/cloud-architecture-golang/internal/usr/mappers"
	"github.com/williamkoller/cloud-architecture-golang/internal/usr/repository"
)

// ---- stub repo ----

type stubRepo struct {
	createFn func(ctx context.Context, u domain.User) error
	getFn    func(ctx context.Context, email vo.Email) (domain.User, bool, error)
	listFn   func(ctx context.Context) ([]domain.User, error)
	updateFn func(ctx context.Context, u domain.User) error
	deleteFn func(ctx context.Context, email vo.Email) error
}

func (s *stubRepo) Create(ctx context.Context, u domain.User) error {
	if s.createFn != nil {
		return s.createFn(ctx, u)
	}
	return nil
}
func (s *stubRepo) GetByEmail(ctx context.Context, email vo.Email) (domain.User, bool, error) {
	if s.getFn != nil {
		return s.getFn(ctx, email)
	}
	return domain.User{}, false, nil
}
func (s *stubRepo) List(ctx context.Context) ([]domain.User, error) {
	if s.listFn != nil {
		return s.listFn(ctx)
	}
	return nil, nil
}
func (s *stubRepo) Update(ctx context.Context, u domain.User) error {
	if s.updateFn != nil {
		return s.updateFn(ctx, u)
	}
	return nil
}
func (s *stubRepo) Delete(ctx context.Context, email vo.Email) error {
	if s.deleteFn != nil {
		return s.deleteFn(ctx, email)
	}
	return nil
}


func routerWithUserRoutes(h *UserHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/users", h.CreateUser)
	r.GET("/users", h.ListUsers)
	r.GET("/users/:email", h.GetUser)
	r.PATCH("/users/:email", h.UpdateUser)
	r.DELETE("/users/:email", h.DeleteUser)
	return r
}

func doJSON(t *testing.T, r http.Handler, method, path string, body any) *httptest.ResponseRecorder {
	t.Helper()
	var buf *bytes.Buffer
	if body != nil {
		bts, err := json.Marshal(body)
		if err != nil {
			t.Fatalf("marshal body: %v", err)
		}
		buf = bytes.NewBuffer(bts)
	} else {
		buf = bytes.NewBuffer(nil)
	}
	req := httptest.NewRequest(method, path, buf)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func mustUser(t *testing.T, name, email string, active bool, ut domain.UserType) domain.User {
	t.Helper()
	u, err := domain.NewUser(name, email, "secret123", active, ut)
	if err != nil {
		t.Fatalf("NewUser: %v", err)
	}
	return u
}


func TestCreateUser_Success(t *testing.T) {
	var captured domain.User
	repo := &stubRepo{
		createFn: func(ctx context.Context, u domain.User) error {
			captured = u
			return nil
		},
	}
	h := NewUserHandler(repo)
	r := routerWithUserRoutes(h)

	body := map[string]any{
		"name":     "Ana",
		"email":    "ana@example.com",
		"password": "secret123",
		"userType": "User",
	}
	w := doJSON(t, r, http.MethodPost, "/users", body)

	if w.Code != http.StatusCreated {
		t.Fatalf("status: got %d, want %d", w.Code, http.StatusCreated)
	}
	var resp mappers.UserResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if resp.Name != "Ana" || resp.Email != "ana@example.com" || resp.UserType != domain.UserTypeUser || resp.Active != true {
		t.Fatalf("response mismatch: %+v", resp)
	}
	// senha deve ser hash (não igual ao raw)
	if string(captured.Password) == "secret123" || !strings.HasPrefix(string(captured.Password), "$2") {
		t.Fatalf("password should be bcrypt hash; got %q", string(captured.Password))
	}
}

func TestCreateUser_Conflict(t *testing.T) {
	repo := &stubRepo{
		createFn: func(ctx context.Context, u domain.User) error {
			return repository.ErrAlreadyExists
		},
	}
	h := NewUserHandler(repo)
	r := routerWithUserRoutes(h)

	body := map[string]any{
		"name":     "Ana",
		"email":    "ana@example.com",
		"password": "secret123",
		"userType": "User",
	}
	w := doJSON(t, r, http.MethodPost, "/users", body)

	if w.Code != http.StatusConflict {
		t.Fatalf("status: got %d, want %d", w.Code, http.StatusConflict)
	}
}

func TestCreateUser_BindingError_Returns400(t *testing.T) {
	repo := &stubRepo{}
	h := NewUserHandler(repo)
	r := routerWithUserRoutes(h)

	body := map[string]any{
		"name": "Ana",
	}
	w := doJSON(t, r, http.MethodPost, "/users", body)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status: got %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestCreateUser_DomainError_Returns422(t *testing.T) {
	repo := &stubRepo{}
	h := NewUserHandler(repo)
	r := routerWithUserRoutes(h)

	body := map[string]any{
		"name":     "   ",
		"email":    "ana@example.com",
		"password": "secret123",
		"userType": "User",
	}
	w := doJSON(t, r, http.MethodPost, "/users", body)
	if w.Code != http.StatusUnprocessableEntity {
		t.Fatalf("status: got %d, want %d", w.Code, http.StatusUnprocessableEntity)
	}
}

func TestCreateUser_InternalError_Returns500(t *testing.T) {
	repo := &stubRepo{
		createFn: func(ctx context.Context, u domain.User) error {
			return errors.New("boom")
		},
	}
	h := NewUserHandler(repo)
	r := routerWithUserRoutes(h)

	body := map[string]any{
		"name":     "Ana",
		"email":    "ana@example.com",
		"password": "secret123",
		"userType": "User",
	}
	w := doJSON(t, r, http.MethodPost, "/users", body)
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("status: got %d, want %d", w.Code, http.StatusInternalServerError)
	}
}

func TestListUsers_Success(t *testing.T) {
	u1 := mustUser(t, "Ana", "ana@example.com", true, domain.UserTypeUser)
	u2 := mustUser(t, "Bob", "bob@example.com", false, domain.UserTypeAdmin)

	repo := &stubRepo{
		listFn: func(ctx context.Context) ([]domain.User, error) {
			return []domain.User{u1, u2}, nil
		},
	}
	h := NewUserHandler(repo)
	r := routerWithUserRoutes(h)

	w := doJSON(t, r, http.MethodGet, "/users", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("status: got %d, want %d", w.Code, http.StatusOK)
	}
	var resp []mappers.UserResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(resp) != 2 {
		t.Fatalf("len: got %d, want 2", len(resp))
	}
}

func TestListUsers_InternalError(t *testing.T) {
	repo := &stubRepo{
		listFn: func(ctx context.Context) ([]domain.User, error) {
			return nil, errors.New("db down")
		},
	}
	h := NewUserHandler(repo)
	r := routerWithUserRoutes(h)

	w := doJSON(t, r, http.MethodGet, "/users", nil)
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("status: got %d, want %d", w.Code, http.StatusInternalServerError)
	}
}

func TestGetUser_Success(t *testing.T) {
	u := mustUser(t, "Ana", "ana@example.com", true, domain.UserTypeUser)
	repo := &stubRepo{
		getFn: func(ctx context.Context, email vo.Email) (domain.User, bool, error) {
			return u, true, nil
		},
	}
	h := NewUserHandler(repo)
	r := routerWithUserRoutes(h)

	w := doJSON(t, r, http.MethodGet, "/users/ana@example.com", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("status: got %d, want %d", w.Code, http.StatusOK)
	}
	var resp mappers.UserResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if resp.Email != "ana@example.com" {
		t.Fatalf("email mismatch: %q", resp.Email)
	}
}

func TestGetUser_InvalidEmail_Returns400(t *testing.T) {
	repo := &stubRepo{}
	h := NewUserHandler(repo)
	r := routerWithUserRoutes(h)

	w := doJSON(t, r, http.MethodGet, "/users/invalid-email", nil)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status: got %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestGetUser_NotFoundAndInternalError(t *testing.T) {
	repoNF := &stubRepo{
		getFn: func(ctx context.Context, email vo.Email) (domain.User, bool, error) {
			return domain.User{}, false, nil
		},
	}
	h1 := NewUserHandler(repoNF)
	r1 := routerWithUserRoutes(h1)
	if w := doJSON(t, r1, http.MethodGet, "/users/ana@example.com", nil); w.Code != http.StatusNotFound {
		t.Fatalf("not found: got %d, want %d", w.Code, http.StatusNotFound)
	}

	repoErr := &stubRepo{
		getFn: func(ctx context.Context, email vo.Email) (domain.User, bool, error) {
			return domain.User{}, false, errors.New("boom")
		},
	}
	h2 := NewUserHandler(repoErr)
	r2 := routerWithUserRoutes(h2)
	if w := doJSON(t, r2, http.MethodGet, "/users/ana@example.com", nil); w.Code != http.StatusInternalServerError {
		t.Fatalf("internal: got %d, want %d", w.Code, http.StatusInternalServerError)
	}
}

func TestUpdateUser_Success(t *testing.T) {
	current := mustUser(t, "Ana", "ana@example.com", true, domain.UserTypeUser)

	repo := &stubRepo{
		getFn: func(ctx context.Context, email vo.Email) (domain.User, bool, error) {
			return current, true, nil
		},
		updateFn: func(ctx context.Context, u domain.User) error {
			if u.Name != "Ana Paula" || u.Active != false || u.UserType != domain.UserTypeAdmin {
				return errors.New("unexpected update data")
			}
			return nil
		},
	}
	h := NewUserHandler(repo)
	r := routerWithUserRoutes(h)

	body := map[string]any{
		"name":     "Ana Paula",
		"active":   false,
		"userType": "Admin",
	}
	w := doJSON(t, r, http.MethodPatch, "/users/ana@example.com", body)
	if w.Code != http.StatusOK {
		t.Fatalf("status: got %d, want %d", w.Code, http.StatusOK)
	}
}

func TestUpdateUser_BindError_And_DomainError(t *testing.T) {
	current := mustUser(t, "Ana", "ana@example.com", true, domain.UserTypeUser)

	repo1 := &stubRepo{
		getFn: func(ctx context.Context, email vo.Email) (domain.User, bool, error) {
			return current, true, nil
		},
	}
	h1 := NewUserHandler(repo1)
	r1 := routerWithUserRoutes(h1)
	w1 := doJSON(t, r1, http.MethodPatch, "/users/ana@example.com", map[string]any{
		"password": "123", 
	})
	if w1.Code != http.StatusBadRequest {
		t.Fatalf("bind error: got %d, want %d", w1.Code, http.StatusBadRequest)
	}

	repo2 := &stubRepo{
		getFn: func(ctx context.Context, email vo.Email) (domain.User, bool, error) {
			return current, true, nil
		},
	}
	h2 := NewUserHandler(repo2)
	r2 := routerWithUserRoutes(h2)
	w2 := doJSON(t, r2, http.MethodPatch, "/users/ana@example.com", map[string]any{
		"name": "   ",
	})
	if w2.Code != http.StatusUnprocessableEntity {
		t.Fatalf("domain error: got %d, want %d", w2.Code, http.StatusUnprocessableEntity)
	}
}

func TestUpdateUser_NotFound_And_InternalError(t *testing.T) {
	repo1 := &stubRepo{
		getFn: func(ctx context.Context, email vo.Email) (domain.User, bool, error) {
			return domain.User{}, false, nil
		},
	}
	h1 := NewUserHandler(repo1)
	r1 := routerWithUserRoutes(h1)
	if w := doJSON(t, r1, http.MethodPatch, "/users/ana@example.com", map[string]any{"name": "Ana"}); w.Code != http.StatusNotFound {
		t.Fatalf("get not found: got %d, want %d", w.Code, http.StatusNotFound)
	}

	repo2 := &stubRepo{
		getFn: func(ctx context.Context, email vo.Email) (domain.User, bool, error) {
			return domain.User{}, false, errors.New("boom")
		},
	}
	h2 := NewUserHandler(repo2)
	r2 := routerWithUserRoutes(h2)
	if w := doJSON(t, r2, http.MethodPatch, "/users/ana@example.com", map[string]any{"name": "Ana"}); w.Code != http.StatusInternalServerError {
		t.Fatalf("get internal: got %d, want %d", w.Code, http.StatusInternalServerError)
	}

	current := mustUser(t, "Ana", "ana@example.com", true, domain.UserTypeUser)
	repo3 := &stubRepo{
		getFn: func(ctx context.Context, email vo.Email) (domain.User, bool, error) {
			return current, true, nil
		},
		updateFn: func(ctx context.Context, u domain.User) error {
			return repository.ErrNotFound
		},
	}
	h3 := NewUserHandler(repo3)
	r3 := routerWithUserRoutes(h3)
	if w := doJSON(t, r3, http.MethodPatch, "/users/ana@example.com", map[string]any{"name": "Ana"}); w.Code != http.StatusNotFound {
		t.Fatalf("update not found: got %d, want %d", w.Code, http.StatusNotFound)
	}

	repo4 := &stubRepo{
		getFn: func(ctx context.Context, email vo.Email) (domain.User, bool, error) {
			return current, true, nil
		},
		updateFn: func(ctx context.Context, u domain.User) error {
			return errors.New("boom")
		},
	}
	h4 := NewUserHandler(repo4)
	r4 := routerWithUserRoutes(h4)
	if w := doJSON(t, r4, http.MethodPatch, "/users/ana@example.com", map[string]any{"name": "Ana"}); w.Code != http.StatusInternalServerError {
		t.Fatalf("update internal: got %d, want %d", w.Code, http.StatusInternalServerError)
	}
}

func TestDeleteUser_Scenarios(t *testing.T) {
	repo0 := &stubRepo{}
	h0 := NewUserHandler(repo0)
	r0 := routerWithUserRoutes(h0)
	if w := doJSON(t, r0, http.MethodDelete, "/users/invalid-email", nil); w.Code != http.StatusBadRequest {
		t.Fatalf("invalid email: got %d, want %d", w.Code, http.StatusBadRequest)
	}

	repo1 := &stubRepo{
		deleteFn: func(ctx context.Context, email vo.Email) error {
			return repository.ErrNotFound
		},
	}
	h1 := NewUserHandler(repo1)
	r1 := routerWithUserRoutes(h1)
	if w := doJSON(t, r1, http.MethodDelete, "/users/ana@example.com", nil); w.Code != http.StatusNotFound {
		t.Fatalf("not found: got %d, want %d", w.Code, http.StatusNotFound)
	}

	repo2 := &stubRepo{
		deleteFn: func(ctx context.Context, email vo.Email) error {
			return errors.New("boom")
		},
	}
	h2 := NewUserHandler(repo2)
	r2 := routerWithUserRoutes(h2)
	if w := doJSON(t, r2, http.MethodDelete, "/users/ana@example.com", nil); w.Code != http.StatusInternalServerError {
		t.Fatalf("internal: got %d, want %d", w.Code, http.StatusInternalServerError)
	}

	repo3 := &stubRepo{
		deleteFn: func(ctx context.Context, email vo.Email) error {
			return nil
		},
	}
	h3 := NewUserHandler(repo3)
	r3 := routerWithUserRoutes(h3)
	if w := doJSON(t, r3, http.MethodDelete, "/users/ana@example.com", nil); w.Code != http.StatusNoContent {
		t.Fatalf("success: got %d, want %d", w.Code, http.StatusNoContent)
	}
}

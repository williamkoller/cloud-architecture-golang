package usr_router

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	m.Run()
}

func TestMapToUserResponse(t *testing.T) {
	tests := []struct {
		name     string
		input    User
		expected UserResponse
	}{
		{
			name: "mapeamento_usuario_completo",
			input: User{
				Name:  "João Silva",
				Email: "joao@exemplo.com",
			},
			expected: UserResponse{
				Name:  "João Silva",
				Email: "joao@exemplo.com",
			},
		},
		{
			name: "mapeamento_usuario_campos_vazios",
			input: User{
				Name:  "",
				Email: "",
			},
			expected: UserResponse{
				Name:  "",
				Email: "",
			},
		},
		{
			name: "mapeamento_usuario_nome_vazio",
			input: User{
				Name:  "",
				Email: "teste@exemplo.com",
			},
			expected: UserResponse{
				Name:  "",
				Email: "teste@exemplo.com",
			},
		},
		{
			name: "mapeamento_usuario_email_vazio",
			input: User{
				Name:  "Maria Santos",
				Email: "",
			},
			expected: UserResponse{
				Name:  "Maria Santos",
				Email: "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mapToUserResponse(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestMapUsersToResponse(t *testing.T) {
	tests := []struct {
		name     string
		input    []User
		expected []UserResponse
	}{
		{
			name: "lista_usuarios_multiplos",
			input: []User{
				{Name: "João Silva", Email: "joao@exemplo.com"},
				{Name: "Maria Santos", Email: "maria@exemplo.com"},
				{Name: "Pedro Oliveira", Email: "pedro@exemplo.com"},
			},
			expected: []UserResponse{
				{Name: "João Silva", Email: "joao@exemplo.com"},
				{Name: "Maria Santos", Email: "maria@exemplo.com"},
				{Name: "Pedro Oliveira", Email: "pedro@exemplo.com"},
			},
		},
		{
			name:     "lista_usuarios_vazia",
			input:    []User{},
			expected: []UserResponse{},
		},
		{
			name: "lista_usuario_unico",
			input: []User{
				{Name: "Usuário Único", Email: "unico@exemplo.com"},
			},
			expected: []UserResponse{
				{Name: "Usuário Único", Email: "unico@exemplo.com"},
			},
		},
		{
			name:     "lista_usuarios_nil",
			input:    nil,
			expected: []UserResponse{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mapUsersToResponse(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestHandlerUsers(t *testing.T) {
	router := gin.New()
	RegisterUserRoutes(router)

	tests := []struct {
		name           string
		method         string
		path           string
		expectedStatus int
		expectedUsers  []UserResponse
	}{
		{
			name:           "get_users_sucesso",
			method:         http.MethodGet,
			path:           "/users",
			expectedStatus: http.StatusOK,
			expectedUsers: []UserResponse{
				{Name: "William K", Email: "william@mail.com"},
				{Name: "Novo user test", Email: "novo-user@mail.com"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(tt.method, tt.path, nil)
			require.NoError(t, err)

			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))

			var response []UserResponse
			err = json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)
			assert.Equal(t, tt.expectedUsers, response)
		})
	}
}

func TestHandlerUsersMetodoInvalido(t *testing.T) {
	router := gin.New()
	RegisterUserRoutes(router)

	metodosInvalidos := []string{
		http.MethodPost,
		http.MethodPut,
		http.MethodDelete,
		http.MethodPatch,
	}

	for _, metodo := range metodosInvalidos {
		t.Run("metodo_"+metodo+"_nao_permitido", func(t *testing.T) {
			req, err := http.NewRequest(metodo, "/users", nil)
			require.NoError(t, err)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusNotFound, w.Code)
		})
	}
}

func TestRegisterUserRoutes(t *testing.T) {
	tests := []struct {
		name string
		path string
	}{
		{
			name: "registro_rota_users",
			path: "/users",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()

			RegisterUserRoutes(router)

			req, err := http.NewRequest(http.MethodGet, tt.path, nil)
			require.NoError(t, err)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.NotEqual(t, http.StatusNotFound, w.Code)
		})
	}
}

func BenchmarkMapToUserResponse(b *testing.B) {
	user := User{
		Name:  "Benchmark User",
		Email: "benchmark@exemplo.com",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mapToUserResponse(user)
	}
}

func BenchmarkMapUsersToResponse(b *testing.B) {
	users := []User{
		{Name: "User 1", Email: "user1@exemplo.com"},
		{Name: "User 2", Email: "user2@exemplo.com"},
		{Name: "User 3", Email: "user3@exemplo.com"},
		{Name: "User 4", Email: "user4@exemplo.com"},
		{Name: "User 5", Email: "user5@exemplo.com"},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mapUsersToResponse(users)
	}
}

func BenchmarkHandlerUsers(b *testing.B) {
	router := gin.New()
	RegisterUserRoutes(router)

	req, _ := http.NewRequest(http.MethodGet, "/users", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
	}
}

// Testes de estrutura de dados

func TestUserStruct(t *testing.T) {
	user := User{
		Name:  "Test User",
		Email: "test@exemplo.com",
	}

	assert.Equal(t, "Test User", user.Name)
	assert.Equal(t, "test@exemplo.com", user.Email)
}

func TestUserResponseStruct(t *testing.T) {
	userResponse := UserResponse{
		Name:  "Test User",
		Email: "test@exemplo.com",
	}

	assert.Equal(t, "Test User", userResponse.Name)
	assert.Equal(t, "test@exemplo.com", userResponse.Email)

	jsonData, err := json.Marshal(userResponse)
	require.NoError(t, err)

	var unmarshaled map[string]interface{}
	err = json.Unmarshal(jsonData, &unmarshaled)
	require.NoError(t, err)

	assert.Equal(t, "Test User", unmarshaled["name"])
	assert.Equal(t, "test@exemplo.com", unmarshaled["email"])
}

func TestMapUsersToResponseCasosExtremos(t *testing.T) {
	t.Run("lista_grande_usuarios", func(t *testing.T) {
		var users []User
		for i := 0; i < 1000; i++ {
			users = append(users, User{
				Name:  "User " + string(rune(i)),
				Email: "user" + string(rune(i)) + "@exemplo.com",
			})
		}

		result := mapUsersToResponse(users)
		assert.Len(t, result, 1000)
	})

	t.Run("usuarios_com_caracteres_especiais", func(t *testing.T) {
		users := []User{
			{Name: "João José", Email: "joão@exemplo.com"},
			{Name: "María García", Email: "maria@domínio.com"},
			{Name: "张三", Email: "zhang@example.com"},
		}

		result := mapUsersToResponse(users)
		assert.Len(t, result, 3)
		assert.Equal(t, "João José", result[0].Name)
		assert.Equal(t, "María García", result[1].Name)
		assert.Equal(t, "张三", result[2].Name)
	})
}

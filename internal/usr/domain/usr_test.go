package domain

import (
	"strings"
	"testing"

	"github.com/williamkoller/cloud-architecture-golang/internal/usr/domain/vo"
)

func TestUserValidate_MissingFields_OrderAndMessage(t *testing.T) {
	u := User{
		Name:     "   ",           // TrimSpace → vazio
		Email:    vo.Email(""),    // vazio
		Password: vo.Password(""), // vazio
		Active:   false,           // não é obrigatório
		UserType: UserType(""),    // vazio
	}

	err := u.Validate()
	if err == nil {
		t.Fatalf("expected error for missing fields, got nil")
	}

	got := err.Error()
	want := "the following fields are required: Name, Email, Password, UserType"
	if got != want {
		t.Fatalf("error message mismatch:\n got: %q\nwant: %q", got, want)
	}
}

func TestUserValidate_Success(t *testing.T) {
	u := User{
		Name:     "Ana",
		Email:    vo.Email("ana@example.com"),
		Password: vo.Password("$2a$10$fakehashjustforvalidate"), // só precisa ser não-vazio
		Active:   true,
		UserType: UserTypeUser,
	}

	if err := u.Validate(); err != nil {
		t.Fatalf("Validate() unexpected error: %v", err)
	}
}

func TestNewUser_Success(t *testing.T) {
	u, err := NewUser("Ana", "ana@example.com", "secret123", true, UserTypeAdmin)
	if err != nil {
		t.Fatalf("NewUser unexpected error: %v", err)
	}

	if u.Name != "Ana" {
		t.Fatalf("Name: got %q, want %q", u.Name, "Ana")
	}
	if string(u.Email) != "ana@example.com" {
		t.Fatalf("Email: got %q, want %q", string(u.Email), "ana@example.com")
	}

	// >>> Correções aqui: não comparar com texto puro <<<
	// 1) Não deve ser igual ao raw
	if string(u.Password) == "secret123" {
		t.Fatalf("Password must be hashed, but equals raw password")
	}
	// 2) Deve ter cara de bcrypt
	if !strings.HasPrefix(string(u.Password), "$2") {
		t.Fatalf("Password hash does not look like bcrypt: %q", string(u.Password))
	}
	// 3) Compare deve validar a senha crua
	if ok := u.Password.Compare("secret123"); !ok {
		t.Fatalf("Password.Compare should return true for the correct raw password")
	}

	if u.Active != true {
		t.Fatalf("Active: got %v, want %v", u.Active, true)
	}
	if u.UserType != UserTypeAdmin {
		t.Fatalf("UserType: got %v, want %v", u.UserType, UserTypeAdmin)
	}
}

func TestNewUser_InvalidEmail_ReturnsError(t *testing.T) {
	// Email malformado para forçar erro em vo.NewEmail
	if _, err := NewUser("Ana", "invalid-email", "secret123", true, UserTypeUser); err == nil {
		t.Fatalf("expected error for invalid email, got nil")
	}
}

func TestNewUser_EmptyPassword_ReturnsError(t *testing.T) {
	if _, err := NewUser("Ana", "ana@example.com", "", true, UserTypeUser); err == nil {
		t.Fatalf("expected error for empty password, got nil")
	}
}

func TestNewUser_NameOnlySpaces_ReturnsError(t *testing.T) {
	if _, err := NewUser(strings.Repeat(" ", 3), "ana@example.com", "secret123", true, UserTypeUser); err == nil {
		t.Fatalf("expected error for name with only spaces, got nil")
	}
}

package dtos

import (
	"testing"

	"github.com/go-playground/validator/v10"
)

func newValidator() *validator.Validate {
	v := validator.New()
	v.SetTagName("binding")
	return v
}

func TestCreateUserRequest_Valid(t *testing.T) {
	v := newValidator()

	req := CreateUserRequest{
		Name:     "Ana",
		Email:    "ana@example.com",
		Password: "secret123",
		Active:   nil,     
		UserType: "Admin",
	}

	if err := v.Struct(req); err != nil {
		t.Fatalf("expected valid CreateUserRequest, got error: %v", err)
	}

	req.UserType = "User"
	if err := v.Struct(req); err != nil {
		t.Fatalf("expected valid CreateUserRequest (User type), got error: %v", err)
	}
}

func TestCreateUserRequest_InvalidCases(t *testing.T) {
	v := newValidator()

	tests := []struct {
		name string
		req  CreateUserRequest
	}{
		{
			name: "missing all required",
			req:  CreateUserRequest{},
		},
		{
			name: "invalid email",
			req: CreateUserRequest{
				Name:     "Ana",
				Email:    "invalid-email",
				Password: "secret123",
				UserType: "User",
			},
		},
		{
			name: "short password",
			req: CreateUserRequest{
				Name:     "Ana",
				Email:    "ana@example.com",
				Password: "12345", // < 6
				UserType: "User",
			},
		},
		{
			name: "invalid user type",
			req: CreateUserRequest{
				Name:     "Ana",
				Email:    "ana@example.com",
				Password: "secret123",
				UserType: "Root",
			},
		},
		{
			name: "empty name",
			req: CreateUserRequest{
				Name:     "",
				Email:    "ana@example.com",
				Password: "secret123",
				UserType: "User",
			},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			if err := v.Struct(tc.req); err == nil {
				t.Fatalf("expected validation error, got nil")
			}
		})
	}
}

func TestUpdateUserRequest_AllOptionalNil_IsValid(t *testing.T) {
	v := newValidator()

	req := UpdateUserRequest{
		Name:     nil,
		Password: nil,
		Active:   nil,
		UserType: nil,
	}

	if err := v.Struct(req); err != nil {
		t.Fatalf("expected valid UpdateUserRequest (all nil), got error: %v", err)
	}
}

func TestUpdateUserRequest_FieldRules(t *testing.T) {
	v := newValidator()

	name := ""
	req := UpdateUserRequest{Name: &name}
	if err := v.Struct(req); err == nil {
		t.Fatalf("expected error when Name is present but empty")
	}

	pass := "12345"
	req = UpdateUserRequest{Password: &pass}
	if err := v.Struct(req); err == nil {
		t.Fatalf("expected error when Password is present but short")
	}

	ut := "Root"
	req = UpdateUserRequest{UserType: &ut}
	if err := v.Struct(req); err == nil {
		t.Fatalf("expected error when UserType is invalid")
	}

	validName := "Ana"
	validPass := "secret123"
	validUT := "Admin"
	active := new(bool)
	*active = true

	req = UpdateUserRequest{
		Name:     &validName,
		Password: &validPass,
		Active:   active,
		UserType: &validUT,
	}
	if err := v.Struct(req); err != nil {
		t.Fatalf("expected valid UpdateUserRequest with all fields present, got error: %v", err)
	}
}

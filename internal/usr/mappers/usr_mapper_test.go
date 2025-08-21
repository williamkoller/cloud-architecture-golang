package mappers

import (
	"testing"

	"github.com/williamkoller/cloud-architecture-golang/internal/usr/domain"
)

func TestToUserResponse_MapsAllFields(t *testing.T) {
	u, err := domain.NewUser(
		"Ana",
		"ana@example.com",
		"secret123",
		true,
		domain.UserTypeAdmin,
	)
	if err != nil {
		t.Fatalf("domain.NewUser: %v", err)
	}

	resp := ToUserResponse(u)

	if resp.Name != "Ana" {
		t.Fatalf("Name: got %q, want %q", resp.Name, "Ana")
	}
	if resp.Email != "ana@example.com" {
		t.Fatalf("Email: got %q, want %q", resp.Email, "ana@example.com")
	}
	if resp.Active != true {
		t.Fatalf("Active: got %v, want %v", resp.Active, true)
	}
	if resp.UserType != domain.UserTypeAdmin {
		t.Fatalf("UserType: got %v, want %v", resp.UserType, domain.UserTypeAdmin)
	}
}

func TestToUserResponse_MapsDifferentValues(t *testing.T) {
	u, err := domain.NewUser(
		"Bruno",
		"bruno@example.com",
		"anotherSecret!",
		false,
		domain.UserTypeUser,
	)
	if err != nil {
		t.Fatalf("domain.NewUser: %v", err)
	}

	resp := ToUserResponse(u)

	if resp.Name != "Bruno" {
		t.Fatalf("Name: got %q, want %q", resp.Name, "Bruno")
	}
	if resp.Email != "bruno@example.com" {
		t.Fatalf("Email: got %q, want %q", resp.Email, "bruno@example.com")
	}
	if resp.Active != false {
		t.Fatalf("Active: got %v, want %v", resp.Active, false)
	}
	if resp.UserType != domain.UserTypeUser {
		t.Fatalf("UserType: got %v, want %v", resp.UserType, domain.UserTypeUser)
	}
}

func TestToUserResponse_ZeroValueUser_DoesNotPanic(t *testing.T) {
	// domain.User zero value — ToUserResponse não deve panicar
	var zero domain.User
	resp := ToUserResponse(zero)

	if resp.Name != "" || resp.Email != "" || resp.Active != false || resp.UserType != "" {
		t.Fatalf("unexpected response from zero value user: %+v", resp)
	}
}

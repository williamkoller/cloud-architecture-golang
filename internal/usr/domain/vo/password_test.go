package vo

import (
	"strings"
	"testing"
)

func TestNewPassword_HashAndCompare(t *testing.T) {
	raw := "secret123"

	p, err := NewPassword(raw)
	if err != nil {
		t.Fatalf("NewPassword error: %v", err)
	}

	// Não deve ser o texto puro
	if string(p) == raw {
		t.Fatalf("hash must not equal raw password")
	}
	// Deve parecer um hash bcrypt ($2*)
	if !strings.HasPrefix(string(p), "$2") {
		t.Fatalf("hash does not look like bcrypt: %q", string(p))
	}

	// Compare correto/errado
	if !p.Compare(raw) {
		t.Fatalf("Compare with correct password should be true")
	}
	if p.Compare("wrongpass") {
		t.Fatalf("Compare with wrong password should be false")
	}
}

func TestNewPassword_MinLength(t *testing.T) {
	_, err := NewPassword("12345") // < 6
	if err == nil {
		t.Fatalf("expected error for short password, got nil")
	}
	if !strings.Contains(err.Error(), "at least 6") {
		t.Fatalf("unexpected error message: %v", err)
	}
}

func TestCompare_InvalidOrZeroValue(t *testing.T) {
	// Hash inválido não deve panicar e deve retornar false
	var invalid Password = "not-a-bcrypt-hash"
	if invalid.Compare("anything") {
		t.Fatalf("Compare on invalid hash should be false")
	}

	// Zero value
	var zero Password
	if zero.Compare("secret") {
		t.Fatalf("Compare on zero value should be false")
	}
}

func TestNewPassword_ProducesDifferentHashes(t *testing.T) {
	raw := "samepassword"

	p1, err := NewPassword(raw)
	if err != nil {
		t.Fatalf("NewPassword #1: %v", err)
	}
	p2, err := NewPassword(raw)
	if err != nil {
		t.Fatalf("NewPassword #2: %v", err)
	}

	// Bcrypt usa salt → hashes devem diferir praticamente sempre
	if string(p1) == string(p2) {
		t.Fatalf("hashes for the same password should differ (got equal)")
	}

	// Ambos devem validar o raw
	if !p1.Compare(raw) || !p2.Compare(raw) {
		t.Fatalf("both hashes should validate the original password")
	}
}

func TestPassword_String(t *testing.T) {
	p := Password("$2a$10$somesalthash...............")
	if got := p.String(); got != string(p) {
		t.Fatalf("String() = %q, want %q", got, string(p))
	}
}

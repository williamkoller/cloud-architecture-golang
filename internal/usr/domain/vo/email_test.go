package vo

import (
	"testing"
)

func TestNewEmail_ValidAddresses(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "simple address",
			input: "ana@example.com",
			want:  "ana@example.com",
		},
		{
			name:  "with display name",
			input: "Ana Silva <ana@example.com>",
			want:  "ana@example.com",
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got, err := NewEmail(tc.input)
			if err != nil {
				t.Fatalf("NewEmail(%q) unexpected error: %v", tc.input, err)
			}
			if string(got) != tc.want {
				t.Fatalf("email value: got %q, want %q", string(got), tc.want)
			}
		})
	}
}

func TestNewEmail_InvalidAddresses(t *testing.T) {
	invalids := []string{
		"",                    
		"plainaddress",        
		"ana@",                
		"@example.com",        
		"ana@example,com",     
		"ana@ example.com",    
		"ana example@example",
	}

	for _, in := range invalids {
		_, err := NewEmail(in)
		if err == nil {
			t.Fatalf("NewEmail(%q) expected error, got nil", in)
		}
	}
}

func TestEmail_String(t *testing.T) {
	e := Email("user@example.com")
	if got := e.String(); got != "user@example.com" {
		t.Fatalf("String(): got %q, want %q", got, "user@example.com")
	}
}

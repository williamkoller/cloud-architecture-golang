package vo

import (
	"errors"
	"fmt"
	"net/mail"
	"strings"
)

type Email string

func NewEmail(value string) (Email, error) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return "", errors.New("invalid email: empty")
	}

	addr, err := mail.ParseAddress(trimmed)
	if err != nil {
		return "", fmt.Errorf("invalid email: %w", err)
	}

	addrSpec := trimmed
	if i := strings.Index(trimmed, "<"); i != -1 {
		if j := strings.Index(trimmed[i:], ">"); j != -1 {
			addrSpec = strings.TrimSpace(trimmed[i+1 : i+j])
		}
	}

	if strings.ContainsAny(addrSpec, " \t\r\n") {
		return "", errors.New("invalid email: whitespace inside address")
	}

	return Email(addr.Address), nil
}

func (e Email) String() string {
	return string(e)
}

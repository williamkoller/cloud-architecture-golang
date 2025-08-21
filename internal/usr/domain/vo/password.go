package vo

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type Password string

func NewPassword(raw string) (Password, error) {
	if len(raw) < 6 {
		return "", errors.New("password must be at least 6 characters long")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(raw), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return Password(hash), nil
}

func (p Password) Compare(raw string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(p), []byte(raw))
	return err == nil
}

func (p Password) String() string {
	return string(p)
}

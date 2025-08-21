package domain

import (
	"fmt"
	"strings"

	"github.com/williamkoller/cloud-architecture-golang/internal/usr/domain/vo"
)

type UserType string

const (
	UserTypeAdmin UserType = "Admin"
	UserTypeUser  UserType = "User"
)

type User struct {
	Name     string
	Email    vo.Email
	Password     vo.Password
	Active   bool
	UserType UserType
}

func (u User) Validate() error {
	var missing []string

	if strings.TrimSpace(u.Name) == "" {
		missing = append(missing, "Name")
	}
	if u.Email == "" {
		missing = append(missing, "Email")
	}
	if u.Password == "" {
		missing = append(missing, "Password")
	}
	if u.UserType == "" {
		missing = append(missing, "UserType")
	}

	if len(missing) > 0 {
		return fmt.Errorf("the following fields are required: %s", strings.Join(missing, ", "))
	}
	return nil
}

func NewUser(name, emailRaw, passRaw string, active bool, userType UserType) (User, error) {
	email, err := vo.NewEmail(emailRaw)
	if err != nil {
		return User{}, err
	}

	pass, err := vo.NewPassword(passRaw)
	if err != nil {
		return User{}, err
	}
	u := User{
		Name:     name,
		Email:    email,
		Password:     pass,
		Active:   active,
		UserType: userType,
	}

	if err := u.Validate(); err != nil {
		return User{}, err
	}

	return u, nil
}

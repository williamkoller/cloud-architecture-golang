package mappers

import "github.com/williamkoller/cloud-architecture-golang/internal/usr/domain"

type UserResponse struct {
	Name     string          `json:"name"`
	Email    string          `json:"email"`
	Active   bool            `json:"active"`
	UserType domain.UserType `json:"userType"`
}

func ToUserResponse(u domain.User) UserResponse {
	return UserResponse{
		Name:     u.Name,
		Email:    string(u.Email),
		Active:   u.Active,
		UserType: u.UserType,
	}
}

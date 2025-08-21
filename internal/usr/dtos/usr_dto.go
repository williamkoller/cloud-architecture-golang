package dtos

type CreateUserRequest struct {
	Name     string `json:"name" binding:"required,min=1"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Active   *bool  `json:"active"`
	UserType string `json:"userType" binding:"required,oneof=Admin User"`
}

type UpdateUserRequest struct {
	Name     *string `json:"name" binding:"omitempty,min=1"`
	Password *string `json:"password" binding:"omitempty,min=6"`
	Active   *bool   `json:"active"`
	UserType *string `json:"userType" binding:"omitempty,oneof=Admin User"`
}

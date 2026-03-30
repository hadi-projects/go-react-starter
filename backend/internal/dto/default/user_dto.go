package dto

type RegisterRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type CreateUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	RoleID   uint   `json:"role_id" binding:"required"`
	Status   string `json:"status" binding:"omitempty,oneof=active freezed pending"`
}

type UpdateUserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email" binding:"omitempty,email"`
	Password string `json:"password" binding:"omitempty,min=6"`
	RoleID   uint   `json:"role_id"`
	Status   string `json:"status" binding:"omitempty,oneof=active freezed pending"`
}

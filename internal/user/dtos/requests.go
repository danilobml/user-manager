package dtos

type RegisterRequest struct {
	Email string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6,max=20"`
	Roles []string `json:"roles" validate:"omitempty,dive,oneof=user admin"`
}

type LoginRequest struct {
	Email string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6,max=20"`
}

type LogoutRequest struct {
	
}

type UnregisterRequest struct {
	Email string `json:"email" validate:"required,email"`
	Token string `json:"token" validate:"required,min=1"`
}

type CheckUserRequest struct {
	
}

type UpdateUserRequest struct {
	Email string `json:"email" validate:"required,email"`
	Roles []string `json:"roles" validate:"omitempty,dive,oneof=user admin"`
}

type ChangePasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6,max=20"`
}

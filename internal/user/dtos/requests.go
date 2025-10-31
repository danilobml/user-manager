package dtos

type RegisterRequest struct {
	Email string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6,max=20"`
	Roles []string `json:"roles" validate:"omitempty,dive,oneof=user admin"`
}

type LoginRequest struct {
	
}

type LogoutRequest struct {
	
}

type UnregisterRequest struct {
	
}

type CheckUserRequest struct {
	
}

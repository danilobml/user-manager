package dtos

import (
	"github.com/google/uuid"
)

type RegisterRequest struct {
	Email    string   `json:"email" validate:"required,email"`
	Password string   `json:"password" validate:"required,min=6,max=20"`
	Roles    []string `json:"roles" validate:"omitempty,dive,oneof=user admin"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6,max=20"`
}

type LogoutRequest struct {
}

type UnregisterRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type CheckUserRequest struct {
}

type UpdateUserRequest struct {
	ID    uuid.UUID `json:"-"`
	Email string    `json:"email" validate:"omitempty,email"`
	Roles []string  `json:"roles" validate:"omitempty,dive,oneof=user admin"`
}

type ChangePasswordRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6,max=20"`
}

type RequestPasswordChangeRequest struct {
	Email string `json:"email" validate:"required,email"`
}

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

type UnregisterRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type CheckUserRequest struct {
	Token  string `json:"token" validate:"required"`
}

type UpdateUserRequest struct {
	ID    uuid.UUID `json:"-"`
	Email string    `json:"email" validate:"omitempty,email"`
	Roles []string  `json:"roles" validate:"omitempty,dive,oneof=user admin"`
}

type RequestPasswordResetRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type ResetPasswordRequest struct {
	Email      string `json:"email" validate:"required,email"`
	Password   string `json:"password" validate:"required,min=6,max=20"`
	ResetToken string `json:"reset_token,omitempty"`
}

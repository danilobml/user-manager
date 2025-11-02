package dtos

import "github.com/google/uuid"

type ResponseUser struct {
	ID       uuid.UUID `json:"id"`
	Email    string    `json:"email"`
	Roles    []string  `json:"roles"`
	IsActive bool      `json:"is_active"`
}

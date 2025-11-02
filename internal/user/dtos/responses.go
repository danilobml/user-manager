package dtos

import "github.com/danilobml/user-manager/internal/user/model"

type RegisterResponse struct {
	Token string `json:"token,omitempty"`
}

type LoginResponse struct {
	Token string `json:"token,omitempty"`
}

type CheckUserResponse struct {
	IsValid bool       `json:"is_valid"`
	User    model.User `json:"user"`
}

type GetAllUsersResponse = []ResponseUser

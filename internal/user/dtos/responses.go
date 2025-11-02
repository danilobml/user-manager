package dtos

type RegisterResponse struct {
	Token string `json:"token,omitempty"`
}

type LoginResponse struct {
	Token string `json:"token,omitempty"`
}

type CheckUserResponse struct {
	IsValid bool `json:"is_valid"`
}

type GetAllUsersResponse = []ResponseUser

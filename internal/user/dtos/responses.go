package dtos

type RegisterResponse struct {
	Token  string `json:"token,omitempty"`
}

type LoginResponse struct {
	Token  string `json:"token,omitempty"`
}

type LogoutResponse struct {
}

type CheckUserResponse struct {
}

type GetAllUsersResponse = []ResponseUser

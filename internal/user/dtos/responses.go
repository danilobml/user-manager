package dtos

type RegisterResponse struct {
	Token  string `json:"token,omitempty"`
}

type LoginResponse struct {
}

type LogoutResponse struct {
}

type UnregisterResponse struct {
}

type CheckUserResponse struct {
}

package handler

import (
	"github.com/danilobml/user-manager/internal/user/service"
)

type UserHandler struct {
	UserService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{
		UserService: userService, 
	}
}


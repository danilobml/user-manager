package service

import (
	"github.com/danilobml/user-manager/internal/user/dtos"
	"github.com/danilobml/user-manager/internal/user/repository"
)

type UserService interface {
	Register(resisterReq dtos.RegisterRequest) (dtos.RegisterResponse, error)
	Login(loginReq dtos.LoginRequest) (dtos.LoginResponse, error)
	Logout(logoutReq dtos.LogoutRequest) (dtos.LogoutResponse, error)
	Unregister(unregisterReq dtos.UnregisterRequest) (dtos.UnregisterResponse, error)
	CheckUser(checkUserReq dtos.CheckUserRequest) (dtos.CheckUserResponse, error)
}

type UserServiceImpl struct {
	UserRepository repository.UserRepository
}

func NewUserserviceImpl(userRepository repository.UserRepository) *UserServiceImpl {
	return &UserServiceImpl{
		UserRepository: userRepository, 
	}
}

func (us *UserServiceImpl) Register(resisterReq dtos.RegisterRequest) (dtos.RegisterResponse, error) {
	return dtos.RegisterResponse{}, nil
}
	
func (us *UserServiceImpl) Login(loginReq dtos.LoginRequest) (dtos.LoginResponse, error) {
	return dtos.LoginResponse{}, nil
}
	
func (us *UserServiceImpl) Logout(logoutReq dtos.LogoutRequest) (dtos.LogoutResponse, error) {
	return dtos.LogoutResponse{}, nil
}
	
func (us *UserServiceImpl) Unregister(unregisterReq dtos.UnregisterRequest) (dtos.UnregisterResponse, error) {
	return dtos.UnregisterResponse{}, nil
}

func (us *UserServiceImpl) CheckUser(checkUserReq dtos.CheckUserRequest) (dtos.CheckUserResponse, error) {
	return dtos.CheckUserResponse{}, nil
}

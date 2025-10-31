package service

import (
	"context"

	"github.com/danilobml/user-manager/internal/helpers"
	"github.com/danilobml/user-manager/internal/user/dtos"
	"github.com/danilobml/user-manager/internal/user/jwt"
	"github.com/danilobml/user-manager/internal/user/model"
	passwordhasher "github.com/danilobml/user-manager/internal/user/password_hasher"
	"github.com/danilobml/user-manager/internal/user/repository"
	"github.com/google/uuid"
)

type UserService interface {
	Register(ctx context.Context, registerReq dtos.RegisterRequest) (dtos.RegisterResponse, error)
	Login(ctx context.Context, loginReq dtos.LoginRequest) (dtos.LoginResponse, error)
	Logout(ctx context.Context, logoutReq dtos.LogoutRequest) (dtos.LogoutResponse, error)
	Unregister(ctx context.Context, unregisterReq dtos.UnregisterRequest) (dtos.UnregisterResponse, error)
	CheckUser(ctx context.Context, checkUserReq dtos.CheckUserRequest) (dtos.CheckUserResponse, error)
	ListAllUsers(ctx context.Context) ([]*model.User, error)
}

type UserServiceImpl struct {
	userRepository repository.UserRepository
	jwtManager *jwt.JwtManager
	passwordHasher passwordhasher.PasswordHasher
}

func NewUserserviceImpl(userRepository repository.UserRepository, jwtManager *jwt.JwtManager) *UserServiceImpl {
	return &UserServiceImpl{
		userRepository: userRepository,
		jwtManager: jwtManager,
		passwordHasher: passwordhasher.NewPasswordHasher(),
	}
}

func (us *UserServiceImpl) Register(ctx context.Context, registerReq dtos.RegisterRequest) (dtos.RegisterResponse, error) {
	hashedPassword, err := us.passwordHasher.HashPassword(registerReq.Password)
	if err != nil {
		return dtos.RegisterResponse{}, err
	}

	id := uuid.New()
	parsedRoles, err := helpers.ParseRoles(registerReq.Roles)
	if err != nil {
		return dtos.RegisterResponse{}, err
	}

	user := model.User{
		ID: id,
		HashedPassword: hashedPassword,
		Email: registerReq.Email,
		Roles: parsedRoles,
	}
	err = us.userRepository.Create(ctx, user)
	if err != nil {
		return dtos.RegisterResponse{}, err
	}

	jwt, err := us.jwtManager.CreateToken(user.Email, registerReq.Roles)
	if err != nil {
		return dtos.RegisterResponse{}, err
	}
	
	return dtos.RegisterResponse{
		Token: jwt,
	}, nil
}

// TODO: implement
func (us *UserServiceImpl) Login(ctx context.Context, loginReq dtos.LoginRequest) (dtos.LoginResponse, error) {
	return dtos.LoginResponse{}, nil
}

// TODO: implement
func (us *UserServiceImpl) Logout(ctx context.Context, logoutReq dtos.LogoutRequest) (dtos.LogoutResponse, error) {
	return dtos.LogoutResponse{}, nil
}

// TODO: implement
func (us *UserServiceImpl) Unregister(ctx context.Context, unregisterReq dtos.UnregisterRequest) (dtos.UnregisterResponse, error) {
	return dtos.UnregisterResponse{}, nil
}

// TODO: implement
func (us *UserServiceImpl) CheckUser(ctx context.Context, checkUserReq dtos.CheckUserRequest) (dtos.CheckUserResponse, error) {
	return dtos.CheckUserResponse{}, nil
}

func (us *UserServiceImpl) ListAllUsers(ctx context.Context) ([]*model.User, error) {
	users, err := us.userRepository.List(ctx)
	if err != nil {
		return nil, err
	}

	return users, nil
}

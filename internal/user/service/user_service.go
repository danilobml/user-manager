package service

import (
	"context"

	"github.com/danilobml/user-manager/internal/errs"
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
	Unregister(ctx context.Context, unregisterRequest dtos.UnregisterRequest) error
	CheckUser(ctx context.Context, checkUserReq dtos.CheckUserRequest) (dtos.CheckUserResponse, error)
	ChangePassword(ctx context.Context, changePassRequest dtos.ChangePasswordRequest) error
	UpdateUserData(ctx context.Context, checkUserReq dtos.CheckUserRequest) error
	ListAllUsers(ctx context.Context) (dtos.GetAllUsersResponse, error)
}

type UserServiceImpl struct {
	userRepository repository.UserRepository
	jwtManager     *jwt.JwtManager
	passwordHasher passwordhasher.PasswordHasher
}

func NewUserserviceImpl(userRepository repository.UserRepository, jwtManager *jwt.JwtManager) *UserServiceImpl {
	return &UserServiceImpl{
		userRepository: userRepository,
		jwtManager:     jwtManager,
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
		ID:             id,
		HashedPassword: hashedPassword,
		Email:          registerReq.Email,
		Roles:          parsedRoles,
		IsActive:       true,
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

func (us *UserServiceImpl) Login(ctx context.Context, loginReq dtos.LoginRequest) (dtos.LoginResponse, error) {
	user, err := us.userRepository.FindByEmail(ctx, loginReq.Email)
	if err != nil {
		return dtos.LoginResponse{}, err
	}

	isPasswordValid := us.passwordHasher.CheckPasswordHash(loginReq.Password, user.HashedPassword)

	if !isPasswordValid {
		return dtos.LoginResponse{}, errs.ErrInvalidCredentials
	}

	rolesStr := helpers.GetRoleNames(user.Roles)

	jwt, err := us.jwtManager.CreateToken(user.Email, rolesStr)
	if err != nil {
		return dtos.LoginResponse{}, err
	}

	return dtos.LoginResponse{
		Token: jwt,
	}, nil
}

// TODO: implement
func (us *UserServiceImpl) Logout(ctx context.Context, logoutReq dtos.LogoutRequest) (dtos.LogoutResponse, error) {
	return dtos.LogoutResponse{}, nil
}

func (us *UserServiceImpl) Unregister(ctx context.Context, unregisterRequest dtos.UnregisterRequest) error {
	user, err := us.userRepository.FindByEmail(ctx, unregisterRequest.Email)
	if err != nil {
		return err
	}

	user.IsActive = false
	err = us.userRepository.Update(ctx, *user)
	if err != nil {
		return err
	}

	return nil
}

// TODO: implement
func (us *UserServiceImpl) CheckUser(ctx context.Context, checkUserReq dtos.CheckUserRequest) (dtos.CheckUserResponse, error) {
	return dtos.CheckUserResponse{}, nil
}

// TODO: implement
func (us *UserServiceImpl) ChangePassword(ctx context.Context, changePassRequest dtos.ChangePasswordRequest) error {
	return nil
}

// TODO: implement
func (us *UserServiceImpl) UpdateUserData(ctx context.Context, checkUserReq dtos.CheckUserRequest) error {
	return nil
}

func (us *UserServiceImpl) ListAllUsers(ctx context.Context) (dtos.GetAllUsersResponse, error) {
	users, err := us.userRepository.List(ctx)
	if err != nil {
		return nil, err
	}

	var respUsers dtos.GetAllUsersResponse
	for _, user := range users {
		roleNames := helpers.GetRoleNames(user.Roles)
		respUser := dtos.ResponseUser{
			ID:    user.ID,
			Email: user.Email,
			Roles: roleNames,
		}
		respUsers = append(respUsers, respUser)
	}

	return respUsers, nil
}

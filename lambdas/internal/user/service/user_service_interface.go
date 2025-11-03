package service

import (
	"context"

	"github.com/danilobml/user-manager/internal/user/dtos"
	"github.com/danilobml/user-manager/internal/user/model"
	"github.com/google/uuid"
)

type UserService interface {
	Register(ctx context.Context, registerReq dtos.RegisterRequest) (dtos.RegisterResponse, error)
	Login(ctx context.Context, loginReq dtos.LoginRequest) (dtos.LoginResponse, error)
	Unregister(ctx context.Context, unregisterRequest dtos.UnregisterRequest) error
	GetUserData(ctx context.Context) (dtos.ResponseUser, error)
	RequestPasswordReset(ctx context.Context, requestPassResetReq dtos.RequestPasswordResetRequest) error
	ResetPassword(ctx context.Context, resetPassRequest dtos.ResetPasswordRequest) error
	ListAllUsers(ctx context.Context) (dtos.GetAllUsersResponse, error)
	UpdateUserData(ctx context.Context, updateUserRequest dtos.UpdateUserRequest) error
	RemoveUser(ctx context.Context, id uuid.UUID) error
	GetUser(ctx context.Context, id uuid.UUID) (*model.User, error)
	CheckUser(ctx context.Context, checkUserReq dtos.CheckUserRequest) (dtos.CheckUserResponse, error)
}

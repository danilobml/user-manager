package service

import (
	"context"
	"fmt"
	"log"

	"github.com/danilobml/user-manager/internal/errs"
	"github.com/danilobml/user-manager/internal/helpers"
	"github.com/danilobml/user-manager/internal/httpx/middleware"
	mailer "github.com/danilobml/user-manager/internal/mailer/service"
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

type UserServiceImpl struct {
	userRepository repository.UserRepository
	jwtManager     *jwt.JwtManager
	passwordHasher passwordhasher.PasswordHasher
	emailService   mailer.Mailer
	baseUrl        string
}

func NewUserserviceImpl(userRepository repository.UserRepository, jwtManager *jwt.JwtManager, emailService mailer.Mailer, baseUrl string) *UserServiceImpl {
	return &UserServiceImpl{
		userRepository: userRepository,
		jwtManager:     jwtManager,
		passwordHasher: passwordhasher.NewPasswordHasher(),
		emailService:   emailService,
		baseUrl:        baseUrl,
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

	jwt, err := us.jwtManager.CreateToken(user.Email, user.Roles)
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
	if user == nil {
		return dtos.LoginResponse{}, errs.ErrInvalidCredentials
	}

	if !user.IsActive {
		return dtos.LoginResponse{}, errs.ErrInvalidCredentials
	}

	isPasswordValid := us.passwordHasher.CheckPasswordHash(loginReq.Password, user.HashedPassword)
	if !isPasswordValid {
		return dtos.LoginResponse{}, errs.ErrInvalidCredentials
	}

	j, err := us.jwtManager.CreateToken(user.Email, user.Roles)
	if err != nil {
		return dtos.LoginResponse{}, err
	}

	return dtos.LoginResponse{Token: j}, nil
}

func (us *UserServiceImpl) GetUserData(ctx context.Context) (dtos.ResponseUser, error) {
	claims, ok := middleware.GetClaimsFromContext(ctx)
	if !ok {
		return dtos.ResponseUser{}, errs.ErrInvalidToken
	}

	user, err := us.userRepository.FindByEmail(ctx, claims.Email)
	if err != nil {
		return dtos.ResponseUser{}, errs.ErrNotFound
	}

	roleNames := helpers.GetRoleNames(user.Roles)

	respUser := dtos.ResponseUser{
		ID: user.ID,
		Email: user.Email,
		Roles: roleNames,
		IsActive: user.IsActive,
	}

	return respUser, nil
}

func (us *UserServiceImpl) Unregister(ctx context.Context, unregisterRequest dtos.UnregisterRequest) error {
	user, err := us.userRepository.FindByEmail(ctx, unregisterRequest.Email)
	if err != nil {
		return err
	}

	// Only the user themselves, or admins can unregister
	if !us.IsUserOwner(ctx, user.Email) && !us.IsUserAdmin(ctx) {
		return errs.ErrUnauthorized
	}

	userToUnregister := model.User{
		ID:             user.ID,
		Email:          user.Email,
		HashedPassword: user.HashedPassword,
		Roles:          user.Roles,
		IsActive:       false,
	}

	err = us.userRepository.Update(ctx, userToUnregister)
	if err != nil {
		return err
	}

	return nil
}

func (us *UserServiceImpl) RequestPasswordReset(ctx context.Context, requestPassResetReq dtos.RequestPasswordResetRequest) error {
	user, err := us.userRepository.FindByEmail(ctx, requestPassResetReq.Email)
	if err != nil || user == nil {
		log.Println("Error sending email: ", err)
		return nil
	}

	token, err := us.jwtManager.CreateResetToken(user.ID.String())
	if err != nil {
		log.Println("Error sending email: ", err)
		return err
	}

	link := fmt.Sprintf("%s/change-password?token=%s", us.baseUrl, token)
	subject := "Password Reset Request"
	body := fmt.Sprintf("Click the link below to reset your password:\r\n\r\n%s\r\n\r\nThis link expires in 15 minutes.", link)

	if err := us.emailService.SendMail([]string{user.Email}, subject, body); err != nil {
		log.Println("Error sending email: ", err)
	}
	return nil
}

func (us *UserServiceImpl) ResetPassword(ctx context.Context, resetPassRequest dtos.ResetPasswordRequest) error {
	userID, err := us.jwtManager.VerifyResetToken(resetPassRequest.ResetToken)
	if err != nil {
		return errs.ErrInvalidToken
	}
	uid, err := uuid.Parse(userID)
	if err != nil {
		return errs.ErrInvalidToken
	}

	user, err := us.userRepository.FindById(ctx, uid)
	if err != nil {
		return errs.ErrInvalidToken
	}

	newHashedPassword, err := us.passwordHasher.HashPassword(resetPassRequest.Password)
	if err != nil {
		return err
	}

	userWithNewPassword := model.User{
		ID:             user.ID,
		Email:          user.Email,
		HashedPassword: newHashedPassword,
		Roles:          user.Roles,
		IsActive:       user.IsActive,
	}

	err = us.userRepository.Update(ctx, userWithNewPassword)
	if err != nil {
		return err
	}

	return nil
}

func (us *UserServiceImpl) UpdateUserData(ctx context.Context, updateUserRequest dtos.UpdateUserRequest) error {
	user, err := us.userRepository.FindById(ctx, updateUserRequest.ID)
	if err != nil {
		return err
	}

	// Only the user themselves, or admins can update data
	if !us.IsUserOwner(ctx, user.Email) && !us.IsUserAdmin(ctx) {
		return errs.ErrUnauthorized
	}

	dbRoles, err := helpers.ParseRoles(updateUserRequest.Roles)
	if err != nil {
		return errs.ErrParsingRoles
	}

	userToUnregister := model.User{
		ID:             user.ID,
		Email:          updateUserRequest.Email,
		HashedPassword: user.HashedPassword,
		Roles:          dbRoles,
		IsActive:       user.IsActive,
	}

	err = us.userRepository.Update(ctx, userToUnregister)
	if err != nil {
		return err
	}

	return nil
}

// Admin only
func (us *UserServiceImpl) ListAllUsers(ctx context.Context) (dtos.GetAllUsersResponse, error) {
	if !us.IsUserAdmin(ctx) {
		return dtos.GetAllUsersResponse{}, errs.ErrUnauthorized
	}

	users, err := us.userRepository.List(ctx)
	if err != nil {
		return nil, err
	}

	if len(users) == 0 {
		return []dtos.ResponseUser{}, nil
	}

	var respUsers dtos.GetAllUsersResponse
	for _, user := range users {
		roleNames := helpers.GetRoleNames(user.Roles)
		respUser := dtos.ResponseUser{
			ID:       user.ID,
			Email:    user.Email,
			Roles:    roleNames,
			IsActive: user.IsActive,
		}
		respUsers = append(respUsers, respUser)
	}

	return respUsers, nil
}

// Admin only
func (us *UserServiceImpl) RemoveUser(ctx context.Context, id uuid.UUID) error {
	// Only admins can remove (delete from DB) an user
	if !us.IsUserAdmin(ctx) {
		return errs.ErrUnauthorized
	}

	err := us.userRepository.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

// For external services:
func (us *UserServiceImpl) CheckUser(ctx context.Context, checkUserReq dtos.CheckUserRequest) (dtos.CheckUserResponse, error) {
	claims, err := us.jwtManager.ParseAndValidateToken(checkUserReq.Token)
	if err != nil {
		return dtos.CheckUserResponse{
			IsValid: false,
		}, err
	}

	user, err := us.userRepository.FindByEmail(ctx, claims.Email)
	if err != nil {
		return dtos.CheckUserResponse{
			IsValid: false,
		}, err
	}

	if !user.IsActive {
		return dtos.CheckUserResponse{
			IsValid: false,
			User:    *user,
		}, nil
	}

	return dtos.CheckUserResponse{
		IsValid: true,
		User:    *user,
	}, nil
}

// Not exposed
func (us *UserServiceImpl) GetUser(ctx context.Context, id uuid.UUID) (*model.User, error) {
	user, err := us.userRepository.FindById(ctx, id)
	if err != nil {
		return nil, errs.ErrNotFound
	}

	return user, nil
}

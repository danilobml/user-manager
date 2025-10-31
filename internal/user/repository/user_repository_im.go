package repository

import (
	"context"
	"slices"

	"github.com/google/uuid"

	"github.com/danilobml/user-manager/internal/errs"
	"github.com/danilobml/user-manager/internal/user/model"
)

type UserRepositoryInMemory struct {
	data []model.User
}

func NewUserRepositoryInMemory() *UserRepositoryInMemory {
	return &UserRepositoryInMemory{
		data: make([]model.User, 0),
	}
}

func (ur *UserRepositoryInMemory) List(ctx context.Context) ([]*model.User, error) {
	var usersResp []*model.User
	for _, user := range ur.data {
		usersResp = append(usersResp, &user)
	}

	return usersResp, nil
}

func (ur *UserRepositoryInMemory) FindById(ctx context.Context, id uuid.UUID) (*model.User, error) {
	for _, user := range ur.data {
		if user.ID == id {
			return &user, nil
		}
	}

	return nil, errs.ErrNotFound
}

func (ur *UserRepositoryInMemory) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	for _, user := range ur.data {
		if user.Email == email {
			return &user, nil
		}
	}

	return nil, errs.ErrNotFound
}

func (ur *UserRepositoryInMemory) Create(ctx context.Context, user model.User) error {
	existingUser, _ := ur.FindByEmail(ctx, user.Email)
	if existingUser != nil {
		return errs.ErrAlreadyExists
	}

	ur.data = append(ur.data, user)

	return nil
}

func (ur *UserRepositoryInMemory) Update(ctx context.Context, user model.User) error {
	existingUser, err := ur.FindByEmail(ctx, user.Email)
	if err != nil {
		return err
	}

	existingUser.Email = user.Email
	existingUser.HashedPassword = user.HashedPassword
	existingUser.Roles = user.Roles

	return nil
}

func (ur *UserRepositoryInMemory) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := ur.FindById(ctx, id)
	if err != nil {
		return err
	}

	ur.data = slices.DeleteFunc(ur.data, func(user model.User) bool {
		return user.ID == id
	})

	return nil
}

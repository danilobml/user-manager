package repository

import (
	"context"

	"github.com/danilobml/user-manager/internal/user/model"
	"github.com/google/uuid"
)

type UserRepository interface {
	List(ctx context.Context) ([]*model.User, error)
	FindById(ctx context.Context, id uuid.UUID) (*model.User, error)
	FindByEmail(ctx context.Context, email string) (*model.User, error)
	Create(ctx context.Context, user model.User) error
	Update(ctx context.Context, user model.User) error
	Delete(ctx context.Context, id uuid.UUID) error
}

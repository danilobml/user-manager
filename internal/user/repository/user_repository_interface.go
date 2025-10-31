package repository

import (
	"github.com/danilobml/user-manager/internal/user/model"
	"github.com/google/uuid"
)

type UserRepository interface {
	List() []*model.User
	FindById(id uuid.UUID) (*model.User, error)
	FindByEmail(email string) (*model.User, error)
	Save(user model.User) (*model.User, error)
	Delete(id uuid.UUID) error
}

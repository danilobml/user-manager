package repository

import (
	"github.com/danilobml/user-manager/internal/user/model"
	"github.com/google/uuid"
)

// TODO: implement logic for DynamoDb
type UserRepositoryDdb struct {
	
}

func NewUserRepositoryDdb() *UserRepositoryDdb {
	return &UserRepositoryDdb{}
}
	
func (ur *UserRepositoryDdb) List() []*model.User {
	return []*model.User{}
}

func (ur *UserRepositoryDdb) FindById(id uuid.UUID) (*model.User, error) {
	return nil, nil
}

func (ur *UserRepositoryDdb) FindByEmail(email string) (*model.User, error) {
	return nil, nil
}

func (ur *UserRepositoryDdb) Save(user model.User) (*model.User, error) {
	return nil, nil
}

func (ur *UserRepositoryDdb) Delete(id uuid.UUID) error {
	return nil
}

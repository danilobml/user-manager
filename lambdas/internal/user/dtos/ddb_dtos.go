package dtos

import (
	"github.com/danilobml/user-manager/internal/user/model"
	"github.com/google/uuid"
)

type UserDDB struct {
	ID             string   `dynamodbav:"id"`
	Email          string   `dynamodbav:"email"`
	HashedPassword string   `dynamodbav:"hashed_password"`
	Roles          []string `dynamodbav:"roles"`
	IsActive       bool     `dynamodbav:"is_active"`
}

func ToDDB(u model.User) UserDDB {
	roleNames := make([]string, 0, len(u.Roles))
	for _, r := range u.Roles {
		roleNames = append(roleNames, r.GetName())
	}
	return UserDDB{
		ID:             u.ID.String(),
		Email:          u.Email,
		HashedPassword: u.HashedPassword,
		Roles:          roleNames,
		IsActive:       u.IsActive,
	}
}

func FromDDB(d UserDDB) (model.User, error) {
	id, err := uuid.Parse(d.ID)
	if err != nil {
		return model.User{}, err
	}
	roles := make([]model.Role, 0, len(d.Roles))
	for _, name := range d.Roles {
		r, err := model.ParseRole(name)
		if err != nil {
			return model.User{}, err
		}
		roles = append(roles, r)
	}
	return model.User{
		ID:             id,
		Email:          d.Email,
		HashedPassword: d.HashedPassword,
		Roles:          roles,
		IsActive:       d.IsActive,
	}, nil
}

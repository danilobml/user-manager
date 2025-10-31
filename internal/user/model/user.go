package model

import (
	"github.com/danilobml/user-manager/internal/errs"
	"github.com/google/uuid"
)

type Role int

const (
	Admin Role = iota
	AppUser
)

var roleName = map[Role]string{
	Admin:   "admin",
	AppUser: "user",
}

func (r Role) GetName() string {
	return roleName[r]
}

func ParseRole(s string) (Role, error) {
	for r, name := range roleName {
		if name == s {
			return r, nil
		}
	}
	return 0, errs.ErrParsingRoles
}

type User struct {
	ID             uuid.UUID `dynamodbav:"id" json:"id"`
	Email          string    `dynamodbav:"email" json:"email"`
	HashedPassword string    `dynamodbav:"hashed_password" json:"-"`
	Roles          []Role    `dynamodbav:"roles" json:"roles"`
}

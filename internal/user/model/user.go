package model

import "github.com/google/uuid"

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

type User struct {
	ID             uuid.UUID `dynamodbav:"id" json:"id"`
	Email          string    `dynamodbav:"email" json:"email"`
	HashedPassword string    `dynamodbav:"hashed_password" json:"hashed_password"`
	Roles          []Role    `dynamodbav:"roles" json:"roles"`
}

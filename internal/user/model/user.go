package model

import "github.com/google/uuid"

type Role int

const (
	Admin = iota
	AppUser
)

var roleName = map[Role]string{
	Admin: "admin",
	AppUser: "user",
}

func (r Role) GetName() string {
	return roleName[r]
}

type User struct {
	Id uuid.UUID
	Email string
	HasedPassword string
	Roles []Role
}

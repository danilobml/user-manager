package service

import (
	"context"

	"github.com/danilobml/user-manager/internal/httpx/middleware"
	"github.com/danilobml/user-manager/internal/user/model"
)

func IsUserOwner(ctx context.Context, user *model.User) bool {
	claims, ok := middleware.GetClaimsFromContext(ctx)
	if !ok {
		return false
	}

	if claims.Email != user.Email {
		return false
	}

	return true
}

func IsUserAdmin(ctx context.Context) bool {
	claims, ok := middleware.GetClaimsFromContext(ctx)
	if !ok {
		return false
	}

	for _, role := range claims.Roles {
		if role.GetName() == "admin" {
			return true
		}
	}

	return false
}

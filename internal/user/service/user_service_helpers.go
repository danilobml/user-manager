package service

import (
	"context"

	"github.com/danilobml/user-manager/internal/httpx/middleware"
)

// Helpers
func (us *UserServiceImpl) IsUserOwner(ctx context.Context, userEmail string) bool {
	_, err := us.userRepository.FindByEmail(ctx, userEmail)
	if err != nil {
		return false
	}

	claims, ok := middleware.GetClaimsFromContext(ctx)
	if !ok {
		return false
	}

	if claims.Email != userEmail {
		return false
	}

	return true
}

func (us *UserServiceImpl) IsUserAdmin(ctx context.Context) bool {
	claims, ok := middleware.GetClaimsFromContext(ctx)
	if !ok {
		return false
	}

	_, err := us.userRepository.FindByEmail(ctx, claims.Email)
	if err != nil {
		return false
	}

	for _, role := range claims.Roles {
		if role.GetName() == "admin" {
			return true
		}
	}

	return false
}

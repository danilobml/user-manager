package helpers

import "github.com/danilobml/user-manager/internal/user/model"

func ParseRoles(names []string) ([]model.Role, error) {
	roles := make([]model.Role, 0, len(names))
	for _, name := range names {
		role, err := model.ParseRole(name)
		if err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}

	return roles, nil
}

func GetRoleNames(roles []model.Role) ([]string) {
	names := make([]string, 0, len(roles))
	for _, role := range roles {
		roleName := role.GetName()
		names = append(names, roleName)
	}

	return names
}
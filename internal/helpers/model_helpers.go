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

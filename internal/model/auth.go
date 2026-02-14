package model

type Auth struct {
	// Login user id
	ID          string
	Roles       []string
	Permissions []string
}

func (a *Auth) HasRole(role string) bool {
	for _, r := range a.Roles {
		if r == role {
			return true
		}
	}
	return false
}

func (a *Auth) HasPermission(permission string) bool {
	for _, p := range a.Permissions {
		if p == permission {
			return true
		}
	}
	return false
}

func (a *Auth) HasAnyRole(roles ...string) bool {
	for _, role := range roles {
		if a.HasRole(role) {
			return true
		}
	}
	return false
}

func (a *Auth) HasAnyPermission(permissions ...string) bool {
	for _, permission := range permissions {
		if a.HasPermission(permission) {
			return true
		}
	}
	return false
}


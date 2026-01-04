package auth

import "e-commerce-api/internal/modules/user"

type LoginResult struct {
	user *user.User
}

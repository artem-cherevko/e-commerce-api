package auth

import (
	"e-commerce-api/internal/modules/models"

	"github.com/appleboy/gin-jwt/v3/core"
)

type LoginResult struct {
	tokenPair *core.Token
}

type RegisterResult struct {
	user      *models.User
	tokenPair *core.Token
}

type RegisterInput struct {
	Name     string `json:"name"`
	UserName string `json:"user_name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

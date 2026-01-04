package auth

import (
	jwt "github.com/appleboy/gin-jwt/v3"
	"gorm.io/gorm"
)

type Service struct {
	db  *gorm.DB
	jwt *jwt.GinJWTMiddleware
}

func NewAuthService(db *gorm.DB, jwt *jwt.GinJWTMiddleware) *Service {
	return &Service{db: db, jwt: jwt}
}

func (s *Service) Login() (*LoginResult, error) {
	return nil, nil
}

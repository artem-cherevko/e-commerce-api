package auth

import (
	"context"
	"crypto/sha256"
	"e-commerce-api/internal/modules/models"
	"encoding/base64"
	"errors"

	jwt "github.com/appleboy/gin-jwt/v3"
	"github.com/appleboy/gin-jwt/v3/core"
	"github.com/matthewhartstonge/argon2"
	"gorm.io/gorm"
)

type Service struct {
	db    *gorm.DB
	jwt   *jwt.GinJWTMiddleware
	argon *argon2.Config
}

func NewAuthService(db *gorm.DB, jwt *jwt.GinJWTMiddleware, argon *argon2.Config) *Service {
	return &Service{db: db, jwt: jwt, argon: argon}
}

func (s *Service) Login(ctx context.Context, email string, password string) (*core.Token, error) {
	var user models.User
	if err := s.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}

	ok, err := argon2.VerifyEncoded([]byte(password), []byte(user.PasswordHash))
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, errors.New("passwords not matches")
	}

	tokenPair, err := s.jwt.TokenGenerator(ctx, &user)
	if err != nil {
		return nil, err
	}

	hash := sha256.Sum256([]byte(tokenPair.RefreshToken))
	refreshTokenHash := base64.StdEncoding.EncodeToString(hash[:])

	session := models.UserSessions{
		UserID:           user.ID,
		RefreshTokenHash: string(refreshTokenHash),
	}

	if err := s.db.WithContext(ctx).Create(&session).Error; err != nil {
		return nil, err
	}

	return tokenPair, nil
}

func (s *Service) Register(ctx context.Context, name string, userName string, email string, password string) (*RegisterResult, error) {
	err := s.db.WithContext(ctx).Where("email = ?", email).First(&models.User{}).Error
	if err == nil {
		return nil, ErrUserAlreadyExists
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	err = s.db.WithContext(ctx).Select("id").Where("user_name = ?", userName).First(&models.User{}).Error
	if err == nil {
		return nil, ErrUserNameTaken
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	passwordHash, err := s.argon.HashEncoded([]byte(password))
	if err != nil {
		return nil, err
	}

	user := models.User{
		Name:         name,
		UserName:     userName,
		Email:        email,
		PasswordHash: string(passwordHash),
	}
	if err := s.db.WithContext(ctx).Create(&user).Error; err != nil {
		return nil, err
	}

	tokenPair, err := s.jwt.TokenGenerator(ctx, &user)
	if err != nil {
		return nil, err
	}

	hash := sha256.Sum256([]byte(tokenPair.RefreshToken))
	refreshTokenHash := base64.StdEncoding.EncodeToString(hash[:])

	userSession := models.UserSessions{
		UserID:           user.ID,
		RefreshTokenHash: string(refreshTokenHash),
	}

	if err := s.db.WithContext(ctx).Create(&userSession).Error; err != nil {
		return nil, err
	}

	return &RegisterResult{
		user:      &user,
		tokenPair: tokenPair,
	}, nil
}

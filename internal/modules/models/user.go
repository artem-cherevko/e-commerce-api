package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID uuid.UUID `json:"id" gorm:"type:uuid;default:gen_random_uuid()"`

	Name         string `json:"name"`
	Email        string `json:"email" gorm:"type:varchar(255);uniqueIndex:idx_users_email"`
	UserName     string `json:"user_name" gorm:"uniqueIndex:idx_users_user_name"`
	PasswordHash string `json:"-"`

	Sessions []UserSessions `gorm:"foreignKey:User"`
}

type UserSessions struct {
	ID uuid.UUID `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`

	UserID           uuid.UUID `gorm:"type:uuid;not null;index"`
	User             User      `json:"-"`
	RefreshTokenHash string

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

package user

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID uuid.UUID `json:"id" gorm:"type:uuid;default:uuid_generate_v4()"`

	Name     string `json:"name"`
	Email    string `json:"email" gorm:"type:varchar(255);uniqueIndex:idx_users_email"`
	UserName string `json:"user_name" gorm:"uniqueIndex:idx_users_user_name"`

	Sessions []UserSessions
}

type UserSessions struct {
	ID               uuid.UUID `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	User             User      `json:"-"`
	RefreshTokenHash string

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

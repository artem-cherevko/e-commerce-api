package user

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID uuid.UUID `json:"id" gorm:"type:uuid;default:uuid_generate_v4()"`

	Name     string `json:"name"`
	UserName string `json:"user_name" gorm:"unique"`

	Sessions []UserSessions
}

type UserSessions struct {
	ID               uuid.UUID `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	User             User      `json:"-"`
	RefreshTokenHash string

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

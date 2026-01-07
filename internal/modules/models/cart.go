package models

import "github.com/google/uuid"

type Cart struct {
	ID     uuid.UUID `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID uuid.UUID `json:"user_id" gorm:"type:uuid;not null;uniqueIndex"`

	CartItems []CartItem
}

type CartItem struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	CartID    uuid.UUID `gorm:"not null;index;uniqueIndex:idx_cart_product"`
	ProductID uuid.UUID `gorm:"not null;index;uniqueIndex:idx_cart_product"`
	Quantity  uint      `gorm:"not null"`

	Product Product `gorm:"foreignKey:ProductID"`
}

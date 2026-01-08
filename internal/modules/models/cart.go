package models

import "github.com/google/uuid"

type Cart struct {
	ID     uuid.UUID `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID uuid.UUID `json:"user_id" gorm:"type:uuid;not null;uniqueIndex"`

	CartItems []CartItem `json:"cart_items"`
}

type CartItem struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	CartID    uuid.UUID `json:"cart_id" gorm:"not null;index;uniqueIndex:idx_cart_product"`
	ProductID uuid.UUID `json:"product_id" gorm:"not null;index;uniqueIndex:idx_cart_product"`
	Quantity  uint      `json:"quantity" gorm:"not null"`

	Product Product `json:"product" gorm:"foreignKey:ProductID"`
}

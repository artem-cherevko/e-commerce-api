package models

import "github.com/google/uuid"

type Product struct {
	ID uuid.UUID `json:"id" gorm:"type:uuid;default:gen_random_uuid()"`

	Name     string `json:"name"`
	Price    uint   `json:"price"`
	Quantity uint   `json:"quantity"`
}

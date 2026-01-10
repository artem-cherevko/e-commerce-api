package payments

import (
	"context"
	"e-commerce-api/internal/modules/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Service struct {
	db *gorm.DB
}

func NewPaymentsService(db *gorm.DB) *Service {
	return &Service{db: db}
}

func (s *Service) CreateCheckout(ctx context.Context, userID uuid.UUID, cartID uuid.UUID) (string, error) {
	var cart models.Cart
	if err := s.db.WithContext(ctx).Preload("CartItems").Preload("CartItems.Product").Where("id = ?", cartID).First(&cart).Error; err != nil {
		return "", err
	}

	link, err := getCheckoutSessionLink(userID.String(), cartID.String(), cart.CartItems)
	if err != nil {
		return "", err
	}

	return link, nil
}

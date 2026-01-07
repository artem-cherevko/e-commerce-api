package cart

import (
	"context"
	"e-commerce-api/internal/modules/models"
	"e-commerce-api/internal/modules/products"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Service struct {
	db              *gorm.DB
	productsService *products.Service
}

func NewCartService(db *gorm.DB, productsService *products.Service) *Service {
	return &Service{db: db, productsService: productsService}
}

func (s *Service) AddProductToCart(
	ctx context.Context,
	userID uuid.UUID,
	productID uuid.UUID,
	quantity uint,
) (*models.Cart, error) {

	returnCart := &models.Cart{}

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		product, err := s.productsService.GetProduct(ctx, productID)
		if err != nil {
			return err
		}

		if quantity > product.Quantity {
			return errors.New("quantity more than available")
		}

		if err := tx.
			Where("user_id = ?", userID).
			FirstOrCreate(returnCart, models.Cart{UserID: userID}).Error; err != nil {
			return err
		}

		var item models.CartItem
		err = tx.
			Where("cart_id = ? AND product_id = ?", returnCart.ID, product.ID).
			First(&item).Error

		if errors.Is(err, gorm.ErrRecordNotFound) {
			item = models.CartItem{
				CartID:    returnCart.ID,
				ProductID: product.ID,
				Quantity:  quantity,
			}
			return tx.Create(&item).Error
		}

		if err != nil {
			return err
		}

		item.Quantity += quantity
		return tx.Save(&item).Error
	})

	if err != nil {
		return nil, err
	}

	return returnCart, nil
}

func (s *Service) RemoveProductFromCart(ctx context.Context, userID uuid.UUID, productID uuid.UUID) error {
	var cart models.Cart
	if err := s.db.WithContext(ctx).Where("user_id = ?", userID).First(&cart).Error; err != nil {
		return err
	}

	res := s.db.WithContext(ctx).
		Delete(&models.CartItem{}, "cart_id = ? AND product_id = ?", cart.ID, productID)

	if res.Error != nil {
		return res.Error
	}

	if res.RowsAffected == 0 {
		return ErrProductNotFoundInCart
	}

	return nil
}

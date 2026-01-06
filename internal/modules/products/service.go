package products

import (
	"context"
	"e-commerce-api/internal/modules/models"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Service struct {
	db *gorm.DB
}

func NewProductsService(db *gorm.DB) *Service {
	return &Service{db: db}
}

func (s *Service) AddProduct(ctx context.Context, name string, price uint, quantity uint) (*models.Product, error) {

	product := models.Product{
		Name:     name,
		Price:    price,
		Quantity: quantity,
	}

	if err := s.db.WithContext(ctx).Create(&product).Error; err != nil {
		return nil, err
	}

	return &product, nil
}

func (s *Service) GetProducts(ctx context.Context) (*[]models.Product, error) {
	var products []models.Product
	result := s.db.WithContext(ctx).First(&products)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, ErrProductsNotFound
	}
	if result.Error != nil {
		return nil, result.Error
	}

	return &products, nil
}

func (s *Service) GetProduct(ctx context.Context, id uuid.UUID) (*models.Product, error) {
	var product models.Product
	err := s.db.WithContext(ctx).Where("id = ?", id).First(&product).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrProductNotFound
	}
	if err != nil {
		return nil, err
	}

	return &product, nil
}

func (s *Service) DeleteProduct(ctx context.Context, id uuid.UUID) error {
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var product models.Product
		errProduct := tx.Where("id = ?", id).First(&product).Error
		if errProduct != nil {
			return errProduct
		}

		if err := tx.Delete(&product).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

package products

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	service *Service
}

func NewProductsHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) AddProduct(c *gin.Context) {
	var input struct {
		Name     string `json:"name"`
		Price    uint   `json:"price"`
		Quantity uint   `json:"quantity"`
	}
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "invalid input"})
		return
	}

	product, err := h.service.AddProduct(c.Request.Context(), input.Name, input.Price, input.Quantity)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "an error occurred while adding product"})
		log.Printf("error while adding product: %s", err.Error())
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": "ok", "product": product})
}

func (h *Handler) GetProducts(c *gin.Context) {
	products, err := h.service.GetProducts(c.Request.Context())
	if errors.Is(err, ErrProductsNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "products not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "an error occurred while finding products"})
		log.Printf("error while getting products: %s", err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok", "products": products})
}

func (h *Handler) GetProduct(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "provide product id"})
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "id not valid"})
		return
	}

	product, err := h.service.GetProduct(c.Request.Context(), id)
	if errors.Is(err, ErrProductNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "product not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "an error occurred while finding product"})
		log.Printf("error while getting products: %s", err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok", "product": product})
}

func (h *Handler) DeleteProduct(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "provide product id"})
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "id not valid"})
		return
	}

	errDelete := h.service.DeleteProduct(c.Request.Context(), id)
	if errors.Is(errDelete, ErrProductsNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": fmt.Sprintf("%s not found", id)})
		return
	}
	if errDelete != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "an error occurred while deleting product"})
		log.Printf("error while deleting product: %s", errDelete.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "deleted"})
}

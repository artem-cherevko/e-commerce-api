package cart

import (
	"e-commerce-api/internal/modules/models"
	"e-commerce-api/internal/modules/products"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	service *Service
}

func NewCartHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) AddProductToCart(c *gin.Context) {
	userAny, exist := c.Get("id")
	if !exist {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "can't get id or product id not provided in params"})
		return
	}
	user := userAny.(*models.User)
	userID := user.ID

	productIDStr := c.Param("id")
	if productIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "can't get id or product id not provided in params"})
		return
	}
	productID, err := uuid.Parse(productIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "product id not valid uuid"})
		return
	}
	quantityStr := c.Query("quantity")
	quantity, err := strconv.ParseUint(quantityStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "quantity not valid (not be above zero)"})
		return
	}

	cart, err := h.service.AddProductToCart(c.Request.Context(), userID, productID, uint(quantity))
	if errors.Is(err, products.ErrProductNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": fmt.Sprintf("product (%s) not found", productIDStr)})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "an error occurred while trying to add product to cart"})
		log.Printf("error while adding product to cart: %s", err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok", "cart": cart})
}

func (h *Handler) RemoveProductFromCart(c *gin.Context) {
	userAny, exist := c.Get("id")
	if !exist {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "can't get id or product id not provided in params"})
		return
	}
	user := userAny.(*models.User)
	userID := user.ID

	productIDStr := c.Param("id")
	if productIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "can't get id or product id not provided in params"})
		return
	}
	productID, err := uuid.Parse(productIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "product id not valid uuid"})
		return
	}

	errCart := h.service.RemoveProductFromCart(c.Request.Context(), userID, productID)
	if errors.Is(errCart, ErrProductNotFoundInCart) {
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": ErrProductNotFoundInCart.Error()})
		return
	}
	if errCart != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "an error occurred while trying to delete product"})
		log.Printf("error while deleting product from card: %s", errCart.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "product deleted from cart"})
}

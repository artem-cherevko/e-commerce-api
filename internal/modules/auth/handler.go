package auth

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Handler struct {
	service *Service
}

func NewAuthHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Login(c *gin.Context) {
	var input LoginInput
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "provide valid input"})
		return
	}

	tokenPair, err := h.service.Login(c.Request.Context(), input.Email, input.Password)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "user not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "can't create token pair"})
		log.Printf("error while trying to create token pair: %s", err.Error())
		return
	}

	c.SetCookie("access_token", tokenPair.AccessToken, int((time.Hour * 24 * 7).Seconds()), "/", "", false, true)
	c.SetCookie("refresh_token", tokenPair.RefreshToken, int(h.service.jwt.MaxRefresh), "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "user is logged in"})
}

func (h *Handler) Register(c *gin.Context) {
	var input RegisterInput
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "provide valid input"})
		log.Println(err)
		return
	}

	result, err := h.service.Register(c.Request.Context(), input.Name, input.UserName, input.Email, input.Password)
	if errors.Is(err, ErrUserAlreadyExists) {
		c.JSON(http.StatusConflict, gin.H{"status": "error", "message": "user already exists"})
		return
	}
	if errors.Is(err, ErrUserNameTaken) {
		c.JSON(http.StatusConflict, gin.H{"status": "error", "message": "username already taken"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "error while creating user"})
		log.Printf("error while creating user: %s", err.Error())
		return
	}

	c.SetCookie("access_token", result.tokenPair.AccessToken, int((time.Hour * 24 * 7).Seconds()), "/", "", false, true)
	c.SetCookie("refresh_token", result.tokenPair.RefreshToken, int(h.service.jwt.MaxRefresh), "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{"status": "ok", "user": result.user})
}

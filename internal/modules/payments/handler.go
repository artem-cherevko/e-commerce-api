package payments

import (
	"e-commerce-api/internal/modules/models"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stripe/stripe-go/v84"
	"github.com/stripe/stripe-go/v84/webhook"
	"gorm.io/gorm"
)

type Handler struct {
	s                   *Service
	stripeWebhookSecret string
}

func NewPaymentsHandler(s *Service, stripeWebhookSecret string) *Handler {
	return &Handler{s: s, stripeWebhookSecret: stripeWebhookSecret}
}

func (h *Handler) CreateCheckout(c *gin.Context) {
	userAny, exist := c.Get("id")
	if !exist {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "can't get id or product id not provided in params"})
		return
	}
	user := userAny.(*models.User)
	userID := user.ID

	cartIDStr := c.Param("id")
	if cartIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "provide cart id in params"})
		return
	}

	cartID, err := uuid.Parse(cartIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "id is not valid uuid"})
		return
	}

	link, err := h.s.CreateCheckout(c.Request.Context(), userID, cartID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "cart not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "can't create checkout session"})
		log.Printf("error while creating checkout: %v", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok", "checkout": link})
}

func (h *Handler) StripeWebhook(c *gin.Context) {
	const MaxBodyBytes = int64(65536)
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, MaxBodyBytes)

	payload, err := io.ReadAll(c.Request.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading request body: %v\n", err)
		c.Status(http.StatusServiceUnavailable)
		return
	}

	event := stripe.Event{}

	if err := json.Unmarshal(payload, &event); err != nil {
		fmt.Fprintf(os.Stderr, "⚠️  Webhook error while parsing basic request. %v\n", err.Error())
		c.Status(http.StatusBadRequest)
		return
	}

	signatureHeader := c.Request.Header.Get("Stripe-Signature")
	event, err = webhook.ConstructEvent(payload, signatureHeader, h.stripeWebhookSecret)
	if err != nil {
		fmt.Fprintf(os.Stderr, "⚠️  Webhook signature verification failed. %v\n", err)
		c.Status(http.StatusBadRequest)
		return
	}

	switch event.Type {
	case "payment_intent.succeeded":
		var paymentIntent stripe.PaymentIntent
		err := json.Unmarshal(event.Data.Raw, &paymentIntent)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing webhook JSON: %v\n", err)
			c.Status(http.StatusBadRequest)
			return
		}
		log.Printf("Successful payment for %d.", paymentIntent.Amount)
	}
	c.Status(http.StatusOK)
}

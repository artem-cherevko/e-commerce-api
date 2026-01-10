package payments

import (
	jwt "github.com/appleboy/gin-jwt/v3"
	"github.com/gin-gonic/gin"
)

func NewPaymentsRouter(r *gin.RouterGroup, h *Handler, mv *jwt.GinJWTMiddleware) {
	payments := r.Group("payments")
	payments.Use(mv.MiddlewareFunc())
	{
		payments.POST("/checkout/:id/create", h.CreateCheckout)
	}
	paymentsNoMV := r.Group("payments")
	{
		paymentsNoMV.POST("/webhook", h.StripeWebhook)
	}
}

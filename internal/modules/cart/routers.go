package cart

import (
	jwt "github.com/appleboy/gin-jwt/v3"
	"github.com/gin-gonic/gin"
)

func RegisterCartRouters(r *gin.RouterGroup, h *Handler, mv *jwt.GinJWTMiddleware) {
	cart := r.Group("/cart")
	cart.Use(mv.MiddlewareFunc())
	{
		cart.POST("/add-product/:id", h.AddProductToCart)
		cart.DELETE("/remove-product/:id", h.RemoveProductFromCart)
	}
}

package products

import (
	jwt "github.com/appleboy/gin-jwt/v3"
	"github.com/gin-gonic/gin"
)

func RegisterProductsRouters(r *gin.RouterGroup, h *Handler, mv *jwt.GinJWTMiddleware) {
	products := r.Group("/products")
	products.Use(mv.MiddlewareFunc())
	{
		products.GET("/", h.GetProducts)
		products.POST("/add", h.AddProduct)
		products.DELETE("/delete/:id", h.DeleteProduct)
	}
}

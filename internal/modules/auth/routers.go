package auth

import "github.com/gin-gonic/gin"

func RegisterAuthRouters(r *gin.RouterGroup, h *Handler) {
	auth := r.Group("/auth")
	{
		auth.POST("/login", h.Login)
		auth.POST("/register", h.Register)
	}
}

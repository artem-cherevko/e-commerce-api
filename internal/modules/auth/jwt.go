package auth

import (
	"e-commerce-api/internal/modules/models"
	"time"

	jwt "github.com/appleboy/gin-jwt/v3"
	"github.com/gin-gonic/gin"
	gojwt "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
	identityKey = "id"
)

func NewJWTMiddleware(secret string) (*jwt.GinJWTMiddleware, error) {
	authMiddleware, err := jwt.New(initParams(secret))
	if err != nil {
		return nil, err
	}
	errInit := authMiddleware.MiddlewareInit()
	if errInit != nil {
		return nil, err
	}

	return authMiddleware, nil
}

func initParams(secret string) *jwt.GinJWTMiddleware {
	return &jwt.GinJWTMiddleware{
		Realm:       "e-commerce",
		Key:         []byte(secret),
		Timeout:     time.Minute * 5,
		MaxRefresh:  time.Hour * 24 * 7,
		TokenLookup: "cookie: access_token",

		IdentityKey:     identityKey,
		IdentityHandler: identityHandler(),
		PayloadFunc:     payloadFunc(),
		Unauthorized:    unauthorized(),
	}
}

func identityHandler() func(c *gin.Context) any {
	return func(c *gin.Context) any {
		claims := jwt.ExtractClaims(c)
		idStr, ok := claims[identityKey].(string)
		if !ok {
			return nil
		}

		id, err := uuid.Parse(idStr)
		if err != nil {
			return nil
		}
		return &models.User{
			ID: id,
		}
	}
}

func payloadFunc() func(data any) gojwt.MapClaims {
	return func(data any) gojwt.MapClaims {
		if v, ok := data.(*models.User); ok {
			return gojwt.MapClaims{
				identityKey: v.ID.String(),
			}
		}
		return gojwt.MapClaims{}
	}
}

func unauthorized() func(c *gin.Context, code int, message string) {
	return func(c *gin.Context, code int, message string) {
		c.JSON(code, gin.H{
			"code":    code,
			"message": message,
		})
	}
}

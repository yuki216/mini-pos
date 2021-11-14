package middleware

import (
	"github.com/golang-jwt/jwt"
	"mini_pos/config"

	echo "github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)


type JwtCustomClaims struct {
	Name     string `json:"name"`
	ID       int    `json:"id"`
	RoleID int    `json:"role_id"`
	UserName string `json:"username"`
	jwt.StandardClaims
}

var tokenCtxKey = &contextKey{"token"}

type contextKey struct {
	name string
}


func JWTMiddleware(cfg config.Config ) echo.MiddlewareFunc {
	return middleware.JWTWithConfig(middleware.JWTConfig{
		SigningMethod: "HS256",
		Claims: &JwtCustomClaims{},
		SigningKey:    []byte(cfg.JWTConfig.Secret),
	})
}

// GetTokenFromContext return Token Object inside context data
func GetTokenFromContext(c echo.Context) *JwtCustomClaims {
	raw, _ := c.Get("user").(*jwt.Token)
	return raw.Claims.(*JwtCustomClaims)
}


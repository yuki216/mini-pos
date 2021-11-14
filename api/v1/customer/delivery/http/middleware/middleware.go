package middleware

import (
	"github.com/labstack/echo/v4"

	"github.com/labstack/echo/v4/middleware"
	"mini_pos/config"
)

// GoMiddleware represent the data-struct for middleware
type GoMiddleware struct {
	// another stuff , may be needed by middleware
}

// CORS will handle the CORS middleware
func (m *GoMiddleware) CORS(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set("Access-Control-Allow-Origin", "*")
		return next(c)
	}
}

// InitMiddleware initialize the middleware
func InitMiddleware(cfg config.Config ) echo.MiddlewareFunc {
	return middleware.JWTWithConfig(middleware.JWTConfig{
		SigningMethod: "HS256",
		Claims: &JwtCustomClaims{},
		SigningKey:    []byte(cfg.JWTConfig.Secret),
	})
}

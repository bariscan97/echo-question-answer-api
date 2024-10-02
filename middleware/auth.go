package middleware

import (
	model "articles-api/models/token"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")

		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			if c.Request().Method == "GET" {
				return next(c)
			}
			return c.JSON(http.StatusUnauthorized, echo.Map{
				"error": "Missing or invalid Authorization header",
			})
		}

		acces_token := strings.Split(authHeader, " ")[1]

		claims := &model.Claim{}

		secretKey := os.Getenv("JWT_SECRET")

		if _, err := jwt.ParseWithClaims(acces_token, claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); ok {
				return []byte(secretKey), nil
			}
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}); err != nil {
			return c.JSON(http.StatusUnauthorized, echo.Map{
				"error": err.Error(),
			})
		}

		if err := claims.Valid(); err != nil {
			return c.JSON(http.StatusUnauthorized, echo.Map{
				"error": "Invalid or expired token",
			})
		}

	c.Set("user", claims)

		return next(c)
	}
}

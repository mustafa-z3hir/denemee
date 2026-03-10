package middleware

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func Auth(jwtSecret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(401).JSON(fiber.Map{"error": "missing authorization header"})
		}
		
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(401).JSON(fiber.Map{"error": "invalid authorization header"})
		}
		
		tokenStr := parts[1]
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method")
			}
			return []byte(jwtSecret), nil
		})
		
		if err != nil || !token.Valid {
			return c.Status(401).JSON(fiber.Map{"error": "invalid token"})
		}
		
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.Status(401).JSON(fiber.Map{"error": "invalid token claims"})
		}
		
		userID, ok := claims["user_id"].(string)
		if !ok {
			return c.Status(401).JSON(fiber.Map{"error": "invalid user_id in token"})
		}
		
		c.Locals("user_id", userID)
		return c.Next()
	}
}

func AdminOnly() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// TODO: Check from database if user is admin
		// For now, allow all authenticated users
		return c.Next()
	}
}
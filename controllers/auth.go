package controllers

import (
	"encoding/base64"
	"goFactory/models"
	"strings"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func BasicAuthMiddleware(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		auth := c.Get("Authorization")
		if auth == "" {
			c.Response().Header.Add("WWW-Authenticate", `Basic realm="Restricted"`)
			return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
		}

		parts := strings.SplitN(auth, " ", 2)
		if len(parts) != 2 || parts[0] != "Basic" {
			return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
		}

		payload, err := base64.StdEncoding.DecodeString(parts[1])
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
		}

		pair := strings.SplitN(string(payload), ":", 2)
		if len(pair) != 2 {
			return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
		}

		username, password := pair[0], pair[1]

		var user models.User
		if err := db.Where("username = ? AND password = ?", username, password).First(&user).Error; err != nil {
			return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
		}

		return c.Next()
	}
}

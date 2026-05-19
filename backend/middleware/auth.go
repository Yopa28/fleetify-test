package middleware

import (
	"fleetify-backend/database"
	"fleetify-backend/models"

	"github.com/gofiber/fiber/v2"
)

func RequireRole(allowedRoles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID := c.Get("X-User-ID")

		if userID == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "X-User-ID header is required",
			})
		}

		var user models.User
		if err := database.DB.First(&user, userID).Error; err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Invalid user",
			})
		}

		for _, role := range allowedRoles {
			if user.Role == role {
				c.Locals("user", user)
				return c.Next()
			}
		}

		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "You are not allowed to access this resource",
		})
	}
}

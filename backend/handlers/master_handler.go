package handlers

import (
	"fleetify-backend/database"
	"fleetify-backend/models"

	"github.com/gofiber/fiber/v2"
)

func GetUsers(c *fiber.Ctx) error {
	var users []models.User

	if err := database.DB.Find(&users).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to get users",
		})
	}

	return c.JSON(fiber.Map{
		"data": users,
	})
}

func GetVehicles(c *fiber.Ctx) error {
	var vehicles []models.Vehicle

	if err := database.DB.Find(&vehicles).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to get vehicles",
		})
	}

	return c.JSON(fiber.Map{
		"data": vehicles,
	})
}

func GetMasterItems(c *fiber.Ctx) error {
	var items []models.MasterItem

	if err := database.DB.Find(&items).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to get master items",
		})
	}

	return c.JSON(fiber.Map{
		"data": items,
	})
}
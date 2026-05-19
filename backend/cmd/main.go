package main

import (
	"log"

	"fleetify-backend/database"
	"fleetify-backend/models"
	"fleetify-backend/seeders"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	database.ConnectDatabase()

	database.DB.AutoMigrate(
		&models.User{},
		&models.Vehicle{},
		&models.MasterItem{},
		&models.MaintenanceReport{},
		&models.ReportItem{},
	)

	seeders.SeedData()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Fleetify API Running",
		})
	})

	log.Fatal(app.Listen(":8080"))
}

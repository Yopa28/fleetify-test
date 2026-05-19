package main

import (
	"log"

	"fleetify-backend/database"
	"fleetify-backend/models"
	"fleetify-backend/routes"
	"fleetify-backend/seeders"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	app := fiber.New()
	app.Use(cors.New())

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

	routes.ReportRoutes(app)

	log.Fatal(app.Listen(":8080"))
}

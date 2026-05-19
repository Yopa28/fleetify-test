package routes

import (
	"fleetify-backend/handlers"
	"fleetify-backend/middleware"

	"github.com/gofiber/fiber/v2"
)

func ReportRoutes(app *fiber.App) {
	api := app.Group("/api")

	api.Post("/reports", middleware.RequireRole("SA"), handlers.CreateReport)
}

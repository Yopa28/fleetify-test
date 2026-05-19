package routes

import (
	"fleetify-backend/handlers"
	"fleetify-backend/middleware"

	"github.com/gofiber/fiber/v2"
)

func ReportRoutes(app *fiber.App) {
	api := app.Group("/api")

	api.Get("/reports", handlers.GetReports)
	api.Post("/reports", middleware.RequireRole("SA"), handlers.CreateReport)
	api.Patch("/reports/:id/approve", middleware.RequireRole("APPROVAL"), handlers.ApproveReport)
	api.Patch("/reports/:id/complete", middleware.RequireRole("SA"), handlers.CompleteReport)
}	

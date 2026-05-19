package handlers

import (
	"fleetify-backend/database"
	"fleetify-backend/dto"
	"fleetify-backend/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func CreateReport(c *fiber.Ctx) error {
	var request dto.CreateReportRequest

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	if request.VehicleID == 0 || request.Odometer <= 0 || request.Complaint == "" || len(request.Items) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "vehicle_id, odometer, complaint, and items are required",
		})
	}

	user := c.Locals("user").(models.User)

	err := database.DB.Transaction(func(tx *gorm.DB) error {
		report := models.MaintenanceReport{
			VehicleID:    request.VehicleID,
			CreatedBy:    user.ID,
			Odometer:     request.Odometer,
			Complaint:    request.Complaint,
			Status:       "PENDING_APPROVAL",
			InitialPhoto: request.InitialPhoto,
		}

		if err := tx.Create(&report).Error; err != nil {
			return err
		}

		for _, itemReq := range request.Items {
			var masterItem models.MasterItem

			if err := tx.First(&masterItem, itemReq.ItemID).Error; err != nil {
				return err
			}

			reportItem := models.ReportItem{
				ReportID:      report.ID,
				ItemID:        masterItem.ID,
				Quantity:      itemReq.Quantity,
				PriceSnapshot: masterItem.Price,
			}

			if err := tx.Create(&reportItem).Error; err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create report",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Report created successfully",
	})
}

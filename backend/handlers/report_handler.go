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

func GetReports(c *fiber.Ctx) error {
	var reports []models.MaintenanceReport

	if err := database.DB.
		Preload("Vehicle").
		Preload("User").
		Preload("ReportItems.MasterItem").
		Order("created_at DESC").
		Find(&reports).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to get reports",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Reports retrieved successfully",
		"data":    reports,
	})
}

func ApproveReport(c *fiber.Ctx) error {
	id := c.Params("id")

	var report models.MaintenanceReport

	if err := database.DB.First(&report, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Report not found",
		})
	}

	if report.Status != "PENDING_APPROVAL" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Only pending reports can be approved",
		})
	}

	report.Status = "APPROVED"

	if err := database.DB.Save(&report).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to approve report",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Report approved successfully",
	})
}

func CompleteReport(c *fiber.Ctx) error {
	id := c.Params("id")

	var request dto.CompleteReportRequest

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	if request.ProofPhoto == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "proof_photo is required",
		})
	}

	var report models.MaintenanceReport

	if err := database.DB.First(&report, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Report not found",
		})
	}

	if report.Status != "APPROVED" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Only approved reports can be completed",
		})
	}

	report.Status = "COMPLETED"
	report.ProofPhoto = request.ProofPhoto

	if err := database.DB.Save(&report).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to complete report",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Report completed successfully",
	})
}

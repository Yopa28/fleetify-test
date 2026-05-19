package models

import "time"

type MaintenanceReport struct {
	ID           uint `gorm:"primaryKey"`
	VehicleID    uint
	CreatedBy    uint
	Odometer     int
	Complaint    string
	Status       string `gorm:"default:PENDING_APPROVAL"`
	InitialPhoto string
	ProofPhoto   string
	CreatedAt    time.Time

	Vehicle Vehicle `gorm:"foreignKey:VehicleID"`
	User    User    `gorm:"foreignKey:CreatedBy"`

	ReportItems []ReportItem `gorm:"foreignKey:ReportID"`
}

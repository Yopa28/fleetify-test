package models

type Vehicle struct {
	ID           uint   `gorm:"primaryKey"`
	LicensePlate string `gorm:"not null"`
	Model        string `gorm:"not null"`
}
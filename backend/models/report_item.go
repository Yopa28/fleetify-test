package models

type ReportItem struct {
	ID            uint `gorm:"primaryKey"`
	ReportID      uint
	ItemID        uint
	Quantity      int
	PriceSnapshot float64

	MasterItem MasterItem `gorm:"foreignKey:ItemID"`
}

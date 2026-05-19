package models

type MasterItem struct {
	ID       uint    `gorm:"primaryKey"`
	ItemName string  `gorm:"not null"`
	Type     string  `gorm:"type:enum('PART','SERVICE')"`
	Price    float64 `gorm:"not null"`
}

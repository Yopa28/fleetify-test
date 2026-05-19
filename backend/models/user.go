package models

type User struct {
	ID       uint   `gorm:"primaryKey"`
	Username string `gorm:"not null"`
	Role     string `gorm:"type:enum('SA','APPROVAL')"`
}

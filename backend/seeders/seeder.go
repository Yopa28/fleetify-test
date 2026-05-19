package seeders

import (
	"fleetify-backend/database"
	"fleetify-backend/models"
)

func SeedData() {
	var count int64

	database.DB.Model(&models.User{}).Count(&count)
	if count == 0 {
		users := []models.User{
			{Username: "sandy_sa", Role: "SA"},
			{Username: "manager_approval", Role: "APPROVAL"},
		}
		database.DB.Create(&users)
	}

	database.DB.Model(&models.Vehicle{}).Count(&count)
	if count == 0 {
		vehicles := []models.Vehicle{
			{LicensePlate: "B 1234 TIA", Model: "Toyota Avanza"},
			{LicensePlate: "B 5678 TIA", Model: "Daihatsu Gran Max"},
			{LicensePlate: "B 9012 TIA", Model: "Mitsubishi L300"},
		}
		database.DB.Create(&vehicles)
	}

	database.DB.Model(&models.MasterItem{}).Count(&count)
	if count == 0 {
		items := []models.MasterItem{
			{ItemName: "Oli Mesin", Type: "PART", Price: 150000},
			{ItemName: "Filter Oli", Type: "PART", Price: 75000},
			{ItemName: "Ban", Type: "PART", Price: 650000},
			{ItemName: "Service Rutin", Type: "SERVICE", Price: 200000},
			{ItemName: "Spooring Balancing", Type: "SERVICE", Price: 300000},
		}
		database.DB.Create(&items)
	}
}
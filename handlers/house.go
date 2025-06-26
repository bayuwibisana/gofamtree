package handlers

import (
	"encoding/json"
	"gofamtree/config"
	"gofamtree/models"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type CreateHouseInput struct {
	Name      string `json:"name"`
	CreatedBy uint   `json:"created_by"` // Admin ID
}

type UpdateHouseInput struct {
	Name string `json:"name"`
}

func CreateHouse(w http.ResponseWriter, r *http.Request) {
	var input CreateHouseInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Validate that admin exists
	var admin models.Admin
	if err := config.DB.First(&admin, input.CreatedBy).Error; err != nil {
		http.Error(w, "Admin not found", http.StatusBadRequest)
		return
	}

	house := models.House{
		Name:      input.Name,
		CreatedBy: input.CreatedBy,
		CreatedAt: time.Now(),
	}

	if err := config.DB.Create(&house).Error; err != nil {
		http.Error(w, "Failed to create house", http.StatusInternalServerError)
		return
	}

	// Load the admin relationship
	config.DB.Preload("Admin").First(&house, house.ID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(house)
}

func GetHouses(w http.ResponseWriter, r *http.Request) {
	var houses []models.House
	
	if err := config.DB.Preload("Admin").Find(&houses).Error; err != nil {
		http.Error(w, "Failed to fetch houses", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(houses)
}

func GetHouse(w http.ResponseWriter, r *http.Request) {
	// Extract ID from URL path
	path := strings.TrimPrefix(r.URL.Path, "/houses/")
	id, err := strconv.ParseUint(path, 10, 32)
	if err != nil {
		http.Error(w, "Invalid house ID", http.StatusBadRequest)
		return
	}

	var house models.House
	if err := config.DB.Preload("Admin").Preload("Persons").First(&house, uint(id)).Error; err != nil {
		http.Error(w, "House not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(house)
}

func UpdateHouse(w http.ResponseWriter, r *http.Request) {
	// Extract ID from URL path
	path := strings.TrimPrefix(r.URL.Path, "/houses/")
	id, err := strconv.ParseUint(path, 10, 32)
	if err != nil {
		http.Error(w, "Invalid house ID", http.StatusBadRequest)
		return
	}

	var input UpdateHouseInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	var house models.House
	if err := config.DB.First(&house, uint(id)).Error; err != nil {
		http.Error(w, "House not found", http.StatusNotFound)
		return
	}

	house.Name = input.Name
	if err := config.DB.Save(&house).Error; err != nil {
		http.Error(w, "Failed to update house", http.StatusInternalServerError)
		return
	}

	// Load relationships
	config.DB.Preload("Admin").First(&house, house.ID)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(house)
}

func DeleteHouse(w http.ResponseWriter, r *http.Request) {
	// Extract ID from URL path
	path := strings.TrimPrefix(r.URL.Path, "/houses/")
	id, err := strconv.ParseUint(path, 10, 32)
	if err != nil {
		http.Error(w, "Invalid house ID", http.StatusBadRequest)
		return
	}

	var house models.House
	if err := config.DB.First(&house, uint(id)).Error; err != nil {
		http.Error(w, "House not found", http.StatusNotFound)
		return
	}

	// Delete all relations in this house first
	config.DB.Where("house_id = ?", uint(id)).Delete(&models.Relation{})
	// Delete all persons in this house
	config.DB.Where("house_id = ?", uint(id)).Delete(&models.Person{})
	// Delete the house
	config.DB.Delete(&house)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "House deleted successfully",
	})
}

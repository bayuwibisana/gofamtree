package handlers

import (
	"encoding/json"
	"fmt"
	"gofamtree/config"
	"gofamtree/models"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type CreatePersonInput struct {
	HouseID     uint   `json:"house_id"`
	Name        string `json:"name"`
	Contact     string `json:"contact"`
	Description string `json:"description"`
	Gender      string `json:"gender"` // male/female
	DOB         string `json:"dob"`    // in format YYYY-MM-DD (optional)
}

type UpdatePersonInput struct {
	Name        string `json:"name"`
	Contact     string `json:"contact"`
	Description string `json:"description"`
	Gender      string `json:"gender"`
	DOB         string `json:"dob"`
}

func CreatePerson(w http.ResponseWriter, r *http.Request) {
	var input CreatePersonInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Validate that house exists
	var house models.House
	if err := config.DB.First(&house, input.HouseID).Error; err != nil {
		http.Error(w, "House not found", http.StatusBadRequest)
		return
	}
	


	person := models.Person{
		HouseID:     input.HouseID,
		Name:        input.Name,
		Contact:     input.Contact,
		Description: input.Description,
		Gender:      input.Gender,
		CreatedAt:   time.Now(),
	}

	// Parse DOB if provided
	if input.DOB != "" {
		parsedDOB, err := time.Parse("2006-01-02", input.DOB)
		if err != nil {
			http.Error(w, "Invalid date format. Use YYYY-MM-DD", http.StatusBadRequest)
			return
		}
		person.DOB = &parsedDOB
	}

	if err := config.DB.Create(&person).Error; err != nil {
		http.Error(w, "Failed to create person", http.StatusInternalServerError)
		return
	}

	// Load relationships
	config.DB.Preload("House").First(&person, person.ID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(person)
}

func GetPersons(w http.ResponseWriter, r *http.Request) {
	var persons []models.Person
	// Validate request
	
	// Optional: filter by house_id if provided as query parameter
	houseID := r.URL.Query().Get("house_id")
	
	if houseID != "" {
		if err := config.DB.Where("house_id = ?", houseID).Preload("House").Find(&persons).Error; err != nil {
			http.Error(w, "Failed to fetch persons", http.StatusInternalServerError)
			return
		}
	} else {
		if err := config.DB.Preload("House").Find(&persons).Error; err != nil {
			http.Error(w, "Failed to fetch persons", http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(persons)
}

func GetPerson(w http.ResponseWriter, r *http.Request) {
	
	
	// Extract ID from URL path
	path := strings.TrimPrefix(r.URL.Path, "/persons/")
	id, err := strconv.ParseUint(path, 10, 32)
	if err != nil {
		http.Error(w, "Invalid person ID", http.StatusBadRequest)
		return
	}
	
	var person models.Person
	
	// Debug: Add logging to see what GORM is doing
	fmt.Printf("Looking for person with ID: %d\n", uint(id))
	
	// Debug: Check what table name GORM is using
	tableName := config.DB.NamingStrategy.TableName("Person")
	fmt.Printf("GORM is looking in table: %s\n", tableName)
	
	// Debug: Try a raw SQL query first to verify data exists
	var count int64
	config.DB.Table("persons").Where("id = ?", uint(id)).Count(&count)
	fmt.Printf("Raw SQL count for id %d: %d\n", uint(id), count)
	
	if err := config.DB.Preload("House").First(&person, uint(id)).Error; err != nil {
		fmt.Printf("GORM Error: %v\n", err)
		http.Error(w, "Person not found", http.StatusNotFound)
		return
	}
	
	fmt.Printf("Found person: %+v\n", person)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(person)
}

func UpdatePerson(w http.ResponseWriter, r *http.Request) {
	// Extract ID from URL path
	path := strings.TrimPrefix(r.URL.Path, "/persons/")
	id, err := strconv.ParseUint(path, 10, 32)
	if err != nil {
		http.Error(w, "Invalid person ID", http.StatusBadRequest)
		return
	}

	var input UpdatePersonInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	var person models.Person
	if err := config.DB.First(&person, uint(id)).Error; err != nil {
		http.Error(w, "Person not found", http.StatusNotFound)
		return
	}

	// Update fields
	person.Name = input.Name
	person.Contact = input.Contact
	person.Description = input.Description
	person.Gender = input.Gender

	// Parse DOB if provided
	if input.DOB != "" {
		parsedDOB, err := time.Parse("2006-01-02", input.DOB)
		if err != nil {
			http.Error(w, "Invalid date format. Use YYYY-MM-DD", http.StatusBadRequest)
			return
		}
		person.DOB = &parsedDOB
	} else {
		person.DOB = nil
	}

	if err := config.DB.Save(&person).Error; err != nil {
		http.Error(w, "Failed to update person", http.StatusInternalServerError)
		return
	}

	// Load relationships
	config.DB.Preload("House").First(&person, person.ID)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(person)
}

func DeletePerson(w http.ResponseWriter, r *http.Request) {
	// Extract ID from URL path
	path := strings.TrimPrefix(r.URL.Path, "/persons/")
	id, err := strconv.ParseUint(path, 10, 32)
	if err != nil {
		http.Error(w, "Invalid person ID", http.StatusBadRequest)
		return
	}

	
	var person models.Person
	if err := config.DB.First(&person, uint(id)).Error; err != nil {
		http.Error(w, "Person not found", http.StatusNotFound)
		return
	}

	// Delete all relations involving this person
	config.DB.Where("person_id = ? OR related_to_id = ?", uint(id), uint(id)).Delete(&models.Relation{})
	// Delete the person
	config.DB.Delete(&person)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Person deleted successfully",
	})
}

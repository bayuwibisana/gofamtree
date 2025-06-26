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

type CreateRelationInput struct {
	HouseID      uint   `json:"house_id"`
	PersonID     uint   `json:"person_id"`
	RelatedToID  uint   `json:"related_to_id"`
	RelationType string `json:"relation_type"` // parent/spouse/sibling
}

type UpdateRelationInput struct {
	RelationType string `json:"relation_type"`
}

type FamilyTreeResponse struct {
	House     models.House     `json:"house"`
	Persons   []models.Person  `json:"persons"`
	Relations []models.Relation `json:"relations"`
}

func CreateRelation(w http.ResponseWriter, r *http.Request) {
	var input CreateRelationInput
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

	// Validate that both persons exist and belong to the same house
	var person, relatedTo models.Person
	if err := config.DB.First(&person, input.PersonID).Error; err != nil {
		http.Error(w, "Person not found", http.StatusBadRequest)
		return
	}
	if err := config.DB.First(&relatedTo, input.RelatedToID).Error; err != nil {
		http.Error(w, "Related person not found", http.StatusBadRequest)
		return
	}

	if person.HouseID != input.HouseID || relatedTo.HouseID != input.HouseID {
		http.Error(w, "Both persons must belong to the specified house", http.StatusBadRequest)
		return
	}

	// Check for duplicate relations
	var existingRelation models.Relation
	if err := config.DB.Where("person_id = ? AND related_to_id = ? AND relation_type = ?", 
		input.PersonID, input.RelatedToID, input.RelationType).First(&existingRelation).Error; err == nil {
		http.Error(w, "Relation already exists", http.StatusConflict)
		return
	}

	// Prevent self-relation
	if input.PersonID == input.RelatedToID {
		http.Error(w, "Cannot create relation with self", http.StatusBadRequest)
		return
	}

	relation := models.Relation{
		HouseID:      input.HouseID,
		PersonID:     input.PersonID,
		RelatedToID:  input.RelatedToID,
		RelationType: input.RelationType,
		CreatedAt:    time.Now(),
	}

	if err := config.DB.Create(&relation).Error; err != nil {
		http.Error(w, "Failed to create relation", http.StatusInternalServerError)
		return
	}

	// Load relationships
	config.DB.Preload("House").Preload("Person").Preload("RelatedTo").First(&relation, relation.ID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(relation)
}

func GetRelations(w http.ResponseWriter, r *http.Request) {
	var relations []models.Relation
	
	// Optional: filter by house_id if provided as query parameter
	houseID := r.URL.Query().Get("house_id")
	if houseID != "" {
		if err := config.DB.Where("house_id = ?", houseID).
			Preload("House").Preload("Person").Preload("RelatedTo").
			Find(&relations).Error; err != nil {
			http.Error(w, "Failed to fetch relations", http.StatusInternalServerError)
			return
		}
	} else {
		if err := config.DB.Preload("House").Preload("Person").Preload("RelatedTo").
			Find(&relations).Error; err != nil {
			http.Error(w, "Failed to fetch relations", http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(relations)
}

func GetRelation(w http.ResponseWriter, r *http.Request) {
	// Extract ID from URL path
	path := strings.TrimPrefix(r.URL.Path, "/relations/")
	id, err := strconv.ParseUint(path, 10, 32)
	if err != nil {
		http.Error(w, "Invalid relation ID", http.StatusBadRequest)
		return
	}

	var relation models.Relation
	if err := config.DB.Preload("House").Preload("Person").Preload("RelatedTo").
		First(&relation, uint(id)).Error; err != nil {
		http.Error(w, "Relation not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(relation)
}

func UpdateRelation(w http.ResponseWriter, r *http.Request) {
	// Extract ID from URL path
	path := strings.TrimPrefix(r.URL.Path, "/relations/")
	id, err := strconv.ParseUint(path, 10, 32)
	if err != nil {
		http.Error(w, "Invalid relation ID", http.StatusBadRequest)
		return
	}

	var input UpdateRelationInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	var relation models.Relation
	if err := config.DB.First(&relation, uint(id)).Error; err != nil {
		http.Error(w, "Relation not found", http.StatusNotFound)
		return
	}

	// Update the relation type
	relation.RelationType = input.RelationType

	if err := config.DB.Save(&relation).Error; err != nil {
		http.Error(w, "Failed to update relation", http.StatusInternalServerError)
		return
	}

	// Load relationships
	config.DB.Preload("House").Preload("Person").Preload("RelatedTo").First(&relation, relation.ID)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(relation)
}

func DeleteRelation(w http.ResponseWriter, r *http.Request) {
	// Extract ID from URL path
	path := strings.TrimPrefix(r.URL.Path, "/relations/")
	id, err := strconv.ParseUint(path, 10, 32)
	if err != nil {
		http.Error(w, "Invalid relation ID", http.StatusBadRequest)
		return
	}

	var relation models.Relation
	if err := config.DB.First(&relation, uint(id)).Error; err != nil {
		http.Error(w, "Relation not found", http.StatusNotFound)
		return
	}

	if err := config.DB.Delete(&relation).Error; err != nil {
		http.Error(w, "Failed to delete relation", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Relation deleted successfully",
	})
}

func GetFamilyTree(w http.ResponseWriter, r *http.Request) {
	// Extract house ID from URL path
	path := strings.TrimPrefix(r.URL.Path, "/family-tree/")
	houseID, err := strconv.ParseUint(path, 10, 32)
	if err != nil {
		http.Error(w, "Invalid house ID", http.StatusBadRequest)
		return
	}

	// Get the house
	var house models.House
	if err := config.DB.Preload("Admin").First(&house, uint(houseID)).Error; err != nil {
		http.Error(w, "House not found", http.StatusNotFound)
		return
	}

	// Get all persons in the house
	var persons []models.Person
	if err := config.DB.Where("house_id = ?", uint(houseID)).Find(&persons).Error; err != nil {
		http.Error(w, "Failed to fetch persons", http.StatusInternalServerError)
		return
	}

	// Get all relations in the house
	var relations []models.Relation
	if err := config.DB.Where("house_id = ?", uint(houseID)).
		Preload("Person").Preload("RelatedTo").Find(&relations).Error; err != nil {
		http.Error(w, "Failed to fetch relations", http.StatusInternalServerError)
		return
	}

	response := FamilyTreeResponse{
		House:     house,
		Persons:   persons,
		Relations: relations,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

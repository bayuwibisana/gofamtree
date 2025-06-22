package handlers

import (
	"net/http"
	"strconv"

	"github.com/bayuwibisana/gofamtree/models"
	"github.com/bayuwibisana/gofamtree/services"
	"github.com/gin-gonic/gin"
)

type PersonHandler struct {
	service *services.PersonService
}

func NewPersonHandler(service *services.PersonService) *PersonHandler {
	return &PersonHandler{service: service}
}

func (h *PersonHandler) CreatePerson(c *gin.Context) {
	var req models.PersonRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	person, err := h.service.CreatePerson(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": person})
}

func (h *PersonHandler) GetAllPersons(c *gin.Context) {
	format := c.Query("format")
	
	persons, err := h.service.GetAllPersons()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// If format=relationships-table, return the relationship table view
	if format == "relationships-table" {
		// Create relationship table
		relationships := []map[string]interface{}{}
		houses := map[string][]map[string]interface{}{}

		for _, person := range persons {
			personData := map[string]interface{}{
				"id":         person.ID,
				"name":       person.FirstName + " " + person.LastName,
				"gender":     person.Gender,
				"birth_year": nil,
				"death_year": nil,
				"father":     nil,
				"mother":     nil,
				"house":      person.LastName,
			}

			if person.BirthDate != nil {
				personData["birth_year"] = person.BirthDate.Year()
			}
			if person.DeathDate != nil {
				personData["death_year"] = person.DeathDate.Year()
			}

			// Add parent relationships
			if person.FatherID != nil {
				father, _ := h.service.GetPersonByID(*person.FatherID)
				if father != nil {
					personData["father"] = father.FirstName + " " + father.LastName
				}
			}
			if person.MotherID != nil {
				mother, _ := h.service.GetPersonByID(*person.MotherID)
				if mother != nil {
					personData["mother"] = mother.FirstName + " " + mother.LastName
				}
			}

			relationships = append(relationships, personData)

			// Group by house (last name)
			houseName := person.LastName
			if houses[houseName] == nil {
				houses[houseName] = []map[string]interface{}{}
			}
			houses[houseName] = append(houses[houseName], personData)
		}

		// Create response
		response := map[string]interface{}{
			"relationships_table": relationships,
			"houses": houses,
			"summary": map[string]interface{}{
				"total_people": len(persons),
				"total_houses": len(houses),
			},
		}

		c.JSON(http.StatusOK, gin.H{"data": response})
		return
	}

	// Default format: return regular person list
	c.JSON(http.StatusOK, gin.H{"data": persons})
}

func (h *PersonHandler) GetPersonByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	person, err := h.service.GetPersonByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": person})
}

func (h *PersonHandler) UpdatePerson(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var req models.PersonRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	person, err := h.service.UpdatePerson(uint(id), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": person})
}

func (h *PersonHandler) DeletePerson(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := h.service.DeletePerson(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Person deleted successfully"})
}

func (h *PersonHandler) SearchPersons(c *gin.Context) {
	query := c.Query("q")
	format := c.Query("format")
	
	// Special case: relationships table format
	if format == "relationships-table" {
		persons, err := h.service.GetAllPersons()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Create relationships table
		relationships := []map[string]interface{}{}
		houses := make(map[string][]map[string]interface{})

		for _, person := range persons {
			var father, mother string = "None", "None"
			
			// Get parent names
			if person.FatherID != nil {
				if f, err := h.service.GetPersonByID(*person.FatherID); err == nil {
					father = f.FirstName + " " + f.LastName
				}
			}
			if person.MotherID != nil {
				if m, err := h.service.GetPersonByID(*person.MotherID); err == nil {
					mother = m.FirstName + " " + m.LastName
				}
			}

			// Create relationship entry
			relationshipEntry := map[string]interface{}{
				"id":         person.ID,
				"name":       person.FirstName + " " + person.LastName,
				"gender":     person.Gender,
				"father":     father,
				"mother":     mother,
				"house":      person.LastName,
				"birth_year": 0,
			}

			if person.BirthDate != nil {
				relationshipEntry["birth_year"] = person.BirthDate.Year()
			}

			relationships = append(relationships, relationshipEntry)

			// Group by house
			if houses[person.LastName] == nil {
				houses[person.LastName] = []map[string]interface{}{}
			}
			houses[person.LastName] = append(houses[person.LastName], relationshipEntry)
		}

		response := map[string]interface{}{
			"relationships_table": relationships,
			"houses":             houses,
			"summary": map[string]interface{}{
				"total_people": len(persons),
				"total_houses": len(houses),
			},
		}

		c.JSON(http.StatusOK, gin.H{"data": response})
		return
	}
	
	// Regular search functionality  
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Query parameter 'q' is required"})
		return
	}

	persons, err := h.service.SearchPersons(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": persons})
}

func (h *PersonHandler) GetFamilyTree(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	tree, err := h.service.GetFamilyTree(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": tree})
}

// NEW: Enhanced Family Tree Handlers

func (h *PersonHandler) GetComprehensiveFamilyTree(c *gin.Context) {
	tree, err := h.service.GetComprehensiveFamilyTree()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": tree})
}

func (h *PersonHandler) GetFamilyTreeByHouse(c *gin.Context) {
	houseName := c.Param("house")
	if houseName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "House name is required"})
		return
	}

	house, err := h.service.GetFamilyTreeByHouse(houseName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": house})
}

func (h *PersonHandler) GetPersonWithAllRelationships(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	personWithRelationships, err := h.service.GetPersonWithAllRelationships(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": personWithRelationships})
}

// NEW: Family Relationships Table - what you asked for!
func (h *PersonHandler) GetFamilyRelationshipsTable(c *gin.Context) {
	// Get all persons
	persons, err := h.service.GetAllPersons()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Create relationship table
	relationships := []map[string]interface{}{}
	houses := map[string][]map[string]interface{}{}

	for _, person := range persons {
		personData := map[string]interface{}{
			"id":         person.ID,
			"name":       person.FirstName + " " + person.LastName,
			"gender":     person.Gender,
			"birth_year": nil,
			"death_year": nil,
			"father":     nil,
			"mother":     nil,
			"house":      person.LastName,
		}

		if person.BirthDate != nil {
			personData["birth_year"] = person.BirthDate.Year()
		}
		if person.DeathDate != nil {
			personData["death_year"] = person.DeathDate.Year()
		}

		// Add parent relationships
		if person.FatherID != nil {
			father, _ := h.service.GetPersonByID(*person.FatherID)
			if father != nil {
				personData["father"] = father.FirstName + " " + father.LastName
			}
		}
		if person.MotherID != nil {
			mother, _ := h.service.GetPersonByID(*person.MotherID)
			if mother != nil {
				personData["mother"] = mother.FirstName + " " + mother.LastName
			}
		}

		relationships = append(relationships, personData)

		// Group by house (last name)
		houseName := person.LastName
		if houses[houseName] == nil {
			houses[houseName] = []map[string]interface{}{}
		}
		houses[houseName] = append(houses[houseName], personData)
	}

	// Create response
	response := map[string]interface{}{
		"relationships_table": relationships,
		"houses": houses,
		"summary": map[string]interface{}{
			"total_people": len(persons),
			"total_houses": len(houses),
		},
	}

	c.JSON(http.StatusOK, gin.H{"data": response})
}

// NEW: Family Tree by House (what you suggested!)
func (h *PersonHandler) GetFamilyTreeByHouseSimple(c *gin.Context) {
	houseName := c.Param("house")
	if houseName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "House name is required"})
		return
	}

	// Get all people with this last name
	persons, err := h.service.SearchPersons(houseName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Filter by exact last name match
	houseMembers := []map[string]interface{}{}
	var patriarch, matriarch map[string]interface{}

	for _, person := range persons {
		if person.LastName == houseName {
			memberData := map[string]interface{}{
				"id":         person.ID,
				"name":       person.FirstName + " " + person.LastName,
				"gender":     person.Gender,
				"birth_year": nil,
				"death_year": nil,
				"generation": 1, // Will calculate properly
				"parents":    []string{},
				"children":   []string{},
			}

			if person.BirthDate != nil {
				memberData["birth_year"] = person.BirthDate.Year()
			}
			if person.DeathDate != nil {
				memberData["death_year"] = person.DeathDate.Year()
			}

			// Identify patriarch/matriarch (no parents)
			if person.FatherID == nil && person.MotherID == nil {
				if person.Gender == "male" && patriarch == nil {
					patriarch = memberData
				} else if person.Gender == "female" && matriarch == nil {
					matriarch = memberData
				}
			}

			houseMembers = append(houseMembers, memberData)
		}
	}

	response := map[string]interface{}{
		"house_name": "House " + houseName,
		"patriarch":  patriarch,
		"matriarch":  matriarch,
		"members":    houseMembers,
		"total_members": len(houseMembers),
	}

	c.JSON(http.StatusOK, gin.H{"data": response})
}

// NEW: Simple Family Relationships Table endpoint 
func (h *PersonHandler) GetRelationshipsTable(c *gin.Context) {
	persons, err := h.service.GetAllPersons()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Create relationships table
	relationships := []map[string]interface{}{}
	houses := make(map[string][]map[string]interface{})

	for _, person := range persons {
		var father, mother string = "None", "None"
		
		// Get parent names
		if person.FatherID != nil {
			if f, err := h.service.GetPersonByID(*person.FatherID); err == nil {
				father = f.FirstName + " " + f.LastName
			}
		}
		if person.MotherID != nil {
			if m, err := h.service.GetPersonByID(*person.MotherID); err == nil {
				mother = m.FirstName + " " + m.LastName
			}
		}

		// Create relationship entry
		relationshipEntry := map[string]interface{}{
			"id":         person.ID,
			"name":       person.FirstName + " " + person.LastName,
			"gender":     person.Gender,
			"father":     father,
			"mother":     mother,
			"house":      person.LastName,
			"birth_year": 0,
		}

		if person.BirthDate != nil {
			relationshipEntry["birth_year"] = person.BirthDate.Year()
		}

		relationships = append(relationships, relationshipEntry)

		// Group by house
		if houses[person.LastName] == nil {
			houses[person.LastName] = []map[string]interface{}{}
		}
		houses[person.LastName] = append(houses[person.LastName], relationshipEntry)
	}

	response := map[string]interface{}{
		"relationships_table": relationships,
		"houses":             houses,
		"summary": map[string]interface{}{
			"total_people": len(persons),
			"total_houses": len(houses),
		},
	}

	c.JSON(http.StatusOK, gin.H{"data": response})
} 
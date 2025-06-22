package services

import (
	"errors"
	"fmt"
	"time"
	"github.com/bayuwibisana/gofamtree/models"
	"gorm.io/gorm"
)

type PersonService struct {
	db *gorm.DB
}

func NewPersonService(db *gorm.DB) *PersonService {
	return &PersonService{db: db}
}

func (s *PersonService) CreatePerson(req *models.PersonRequest) (*models.Person, error) {
	person := &models.Person{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Gender:    req.Gender,
		BirthDate: req.BirthDate,
		DeathDate: req.DeathDate,
		PhotoURL:  req.PhotoURL,
		Biography: req.Biography,
		FatherID:  req.FatherID,
		MotherID:  req.MotherID,
	}

	if err := s.db.Create(person).Error; err != nil {
		return nil, err
	}

	return person, nil
}

func (s *PersonService) GetAllPersons() ([]models.Person, error) {
	var persons []models.Person
	if err := s.db.Find(&persons).Error; err != nil {
		return nil, err
	}
	return persons, nil
}

func (s *PersonService) GetPersonByID(id uint) (*models.Person, error) {
	var person models.Person
	if err := s.db.First(&person, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("person not found")
		}
		return nil, err
	}
	return &person, nil
}

func (s *PersonService) UpdatePerson(id uint, req *models.PersonRequest) (*models.Person, error) {
	var person models.Person
	if err := s.db.First(&person, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("person not found")
		}
		return nil, err
	}

	person.FirstName = req.FirstName
	person.LastName = req.LastName
	person.Gender = req.Gender
	person.BirthDate = req.BirthDate
	person.DeathDate = req.DeathDate
	person.PhotoURL = req.PhotoURL
	person.Biography = req.Biography
	person.FatherID = req.FatherID
	person.MotherID = req.MotherID

	if err := s.db.Save(&person).Error; err != nil {
		return nil, err
	}

	return &person, nil
}

func (s *PersonService) DeletePerson(id uint) error {
	if err := s.db.Delete(&models.Person{}, id).Error; err != nil {
		return err
	}
	return nil
}

func (s *PersonService) SearchPersons(query string) ([]models.Person, error) {
	var persons []models.Person
	searchQuery := "%" + query + "%"
	if err := s.db.Where("first_name ILIKE ? OR last_name ILIKE ?", searchQuery, searchQuery).Find(&persons).Error; err != nil {
		return nil, err
	}
	return persons, nil
}

func (s *PersonService) GetFamilyTree(personID uint) (*models.FamilyTreeNode, error) {
	person, err := s.GetPersonByID(personID)
	if err != nil {
		return nil, err
	}

	tree := &models.FamilyTreeNode{
		Person: *person,
	}

	children, err := s.getChildren(personID)
	if err != nil {
		return nil, err
	}

	for _, child := range children {
		childTree, err := s.GetFamilyTree(child.ID)
		if err != nil {
			continue
		}
		tree.Children = append(tree.Children, *childTree)
	}

	return tree, nil
}

func (s *PersonService) getChildren(personID uint) ([]models.Person, error) {
	var children []models.Person
	if err := s.db.Where("father_id = ? OR mother_id = ?", personID, personID).Find(&children).Error; err != nil {
		return nil, err
	}
	return children, nil
}

// NEW: Comprehensive Family Tree Functions

func (s *PersonService) GetComprehensiveFamilyTree() (*models.ComprehensiveFamilyTree, error) {
	houses, err := s.getAllHousesWithMembers()
	if err != nil {
		return nil, err
	}

	relationships, err := s.getAllRelationships()
	if err != nil {
		return nil, err
	}

	generations, err := s.getGenerations()
	if err != nil {
		return nil, err
	}

	statistics, err := s.getFamilyStatistics()
	if err != nil {
		return nil, err
	}

	return &models.ComprehensiveFamilyTree{
		Houses:        houses,
		Relationships: relationships,
		Generations:   generations,
		Statistics:    statistics,
	}, nil
}

func (s *PersonService) GetFamilyTreeByHouse(houseName string) (*models.HouseWithMembers, error) {
	var members []models.Person
	// Get all people with the same last name (assuming house name matches surname)
	if err := s.db.Where("last_name = ?", houseName).Find(&members).Error; err != nil {
		return nil, err
	}

	house := models.House{
		Name:        "House " + houseName,
		Description: fmt.Sprintf("Members of the %s family", houseName),
	}

	return &models.HouseWithMembers{
		House:   house,
		Members: members,
	}, nil
}

func (s *PersonService) GetPersonWithAllRelationships(personID uint) (*models.PersonWithRelationships, error) {
	person, err := s.GetPersonByID(personID)
	if err != nil {
		return nil, err
	}

	result := &models.PersonWithRelationships{
		Person: *person,
	}

	// Get parents
	if person.FatherID != nil {
		father, _ := s.GetPersonByID(*person.FatherID)
		result.Father = father
	}
	if person.MotherID != nil {
		mother, _ := s.GetPersonByID(*person.MotherID)
		result.Mother = mother
	}

	// Get children
	children, _ := s.getChildren(personID)
	result.Children = children

	// Get siblings
	siblings, _ := s.getSiblings(personID)
	result.Siblings = siblings

	// Get grandparents
	grandparents, _ := s.getGrandparents(personID)
	result.Grandparents = grandparents

	// Get grandchildren
	grandchildren, _ := s.getGrandchildren(personID)
	result.Grandchildren = grandchildren

	// Get all relationships (if we implement the relationship table later)
	relationships, _ := s.getPersonRelationships(personID)
	result.Relationships = relationships

	return result, nil
}

func (s *PersonService) getAllHousesWithMembers() ([]models.HouseWithMembers, error) {
	var houses []models.HouseWithMembers
	
	// Get unique last names to create houses
	var lastNames []string
	if err := s.db.Model(&models.Person{}).Distinct("last_name").Pluck("last_name", &lastNames).Error; err != nil {
		return nil, err
	}

	for _, lastName := range lastNames {
		houseWithMembers, err := s.GetFamilyTreeByHouse(lastName)
		if err != nil {
			continue
		}
		houses = append(houses, *houseWithMembers)
	}

	return houses, nil
}

func (s *PersonService) getAllRelationships() ([]models.FamilyRelationship, error) {
	var relationships []models.FamilyRelationship
	var persons []models.Person
	
	if err := s.db.Find(&persons).Error; err != nil {
		return nil, err
	}

	// Generate relationships from person data
	for _, person := range persons {
		// Father relationship
		if person.FatherID != nil {
			relationships = append(relationships, models.FamilyRelationship{
				PersonID:     person.ID,
				RelatedToID:  *person.FatherID,
				Relationship: "child",
			})
			relationships = append(relationships, models.FamilyRelationship{
				PersonID:     *person.FatherID,
				RelatedToID:  person.ID,
				Relationship: "father",
			})
		}

		// Mother relationship
		if person.MotherID != nil {
			relationships = append(relationships, models.FamilyRelationship{
				PersonID:     person.ID,
				RelatedToID:  *person.MotherID,
				Relationship: "child",
			})
			relationships = append(relationships, models.FamilyRelationship{
				PersonID:     *person.MotherID,
				RelatedToID:  person.ID,
				Relationship: "mother",
			})
		}
	}

	return relationships, nil
}

func (s *PersonService) getGenerations() ([]models.Generation, error) {
	var generations []models.Generation
	
	// Find founders (people with no parents)
	var founders []models.Person
	if err := s.db.Where("father_id IS NULL AND mother_id IS NULL").Find(&founders).Error; err != nil {
		return nil, err
	}
	
	if len(founders) > 0 {
		generations = append(generations, models.Generation{
			Level:  1,
			Name:   "Founders",
			People: founders,
		})
	}

	// Find second generation (children of founders)
	if len(founders) > 0 {
		var secondGen []models.Person
		for _, founder := range founders {
			children, _ := s.getChildren(founder.ID)
			secondGen = append(secondGen, children...)
		}
		if len(secondGen) > 0 {
			generations = append(generations, models.Generation{
				Level:  2,
				Name:   "Second Generation",
				People: secondGen,
			})
		}

		// Find third generation (grandchildren)
		var thirdGen []models.Person
		for _, person := range secondGen {
			children, _ := s.getChildren(person.ID)
			thirdGen = append(thirdGen, children...)
		}
		if len(thirdGen) > 0 {
			generations = append(generations, models.Generation{
				Level:  3,
				Name:   "Third Generation",
				People: thirdGen,
			})
		}
	}

	return generations, nil
}

func (s *PersonService) getFamilyStatistics() (models.FamilyStatistics, error) {
	var totalPeople int64
	s.db.Model(&models.Person{}).Count(&totalPeople)

	var livingMembers int64
	s.db.Model(&models.Person{}).Where("death_date IS NULL").Count(&livingMembers)
	
	deceasedMembers := totalPeople - livingMembers

	// Get unique houses count
	var lastNames []string
	s.db.Model(&models.Person{}).Distinct("last_name").Pluck("last_name", &lastNames)

	// Count generations
	generations, _ := s.getGenerations()

	// Gender distribution
	genderCount := make(map[string]int)
	var males, females int64
	s.db.Model(&models.Person{}).Where("gender = ?", "male").Count(&males)
	s.db.Model(&models.Person{}).Where("gender = ?", "female").Count(&females)
	genderCount["male"] = int(males)
	genderCount["female"] = int(females)

	return models.FamilyStatistics{
		TotalPeople:     int(totalPeople),
		TotalHouses:     len(lastNames),
		Generations:     len(generations),
		LivingMembers:   int(livingMembers),
		DeceasedMembers: int(deceasedMembers),
		GenderCount:     genderCount,
	}, nil
}

func (s *PersonService) getSiblings(personID uint) ([]models.Person, error) {
	person, err := s.GetPersonByID(personID)
	if err != nil {
		return nil, err
	}

	var siblings []models.Person
	if person.FatherID != nil || person.MotherID != nil {
		query := s.db.Where("id != ?", personID)
		if person.FatherID != nil {
			query = query.Where("father_id = ?", *person.FatherID)
		}
		if person.MotherID != nil {
			query = query.Or("mother_id = ?", *person.MotherID)
		}
		query.Find(&siblings)
	}

	return siblings, nil
}

func (s *PersonService) getGrandparents(personID uint) ([]models.Person, error) {
	var grandparents []models.Person
	person, err := s.GetPersonByID(personID)
	if err != nil {
		return grandparents, err
	}

	// Get paternal grandparents
	if person.FatherID != nil {
		father, _ := s.GetPersonByID(*person.FatherID)
		if father != nil {
			if father.FatherID != nil {
				paternalGrandfather, _ := s.GetPersonByID(*father.FatherID)
				if paternalGrandfather != nil {
					grandparents = append(grandparents, *paternalGrandfather)
				}
			}
			if father.MotherID != nil {
				paternalGrandmother, _ := s.GetPersonByID(*father.MotherID)
				if paternalGrandmother != nil {
					grandparents = append(grandparents, *paternalGrandmother)
				}
			}
		}
	}

	// Get maternal grandparents
	if person.MotherID != nil {
		mother, _ := s.GetPersonByID(*person.MotherID)
		if mother != nil {
			if mother.FatherID != nil {
				maternalGrandfather, _ := s.GetPersonByID(*mother.FatherID)
				if maternalGrandfather != nil {
					grandparents = append(grandparents, *maternalGrandfather)
				}
			}
			if mother.MotherID != nil {
				maternalGrandmother, _ := s.GetPersonByID(*mother.MotherID)
				if maternalGrandmother != nil {
					grandparents = append(grandparents, *maternalGrandmother)
				}
			}
		}
	}

	return grandparents, nil
}

func (s *PersonService) getGrandchildren(personID uint) ([]models.Person, error) {
	var grandchildren []models.Person
	
	// Get children first
	children, err := s.getChildren(personID)
	if err != nil {
		return grandchildren, err
	}

	// Get grandchildren from each child
	for _, child := range children {
		childChildren, _ := s.getChildren(child.ID)
		grandchildren = append(grandchildren, childChildren...)
	}

	return grandchildren, nil
}

func (s *PersonService) getPersonRelationships(personID uint) ([]models.FamilyRelationship, error) {
	// This would be used if we implement a proper relationships table
	// For now, return empty slice
	return []models.FamilyRelationship{}, nil
} 
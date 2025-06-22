package models

import (
	"time"
	"gorm.io/gorm"
)

type Person struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	FirstName   string    `json:"first_name" gorm:"not null"`
	LastName    string    `json:"last_name" gorm:"not null"`
	Gender      string    `json:"gender" gorm:"type:varchar(10)"`
	BirthDate   *time.Time `json:"birth_date"`
	DeathDate   *time.Time `json:"death_date"`
	PhotoURL    string    `json:"photo_url"`
	Biography   string    `json:"biography" gorm:"type:text"`
	
	// Family relationships
	FatherID    *uint     `json:"father_id"`
	MotherID    *uint     `json:"mother_id"`
	Father      *Person   `json:"father,omitempty" gorm:"foreignKey:FatherID"`
	Mother      *Person   `json:"mother,omitempty" gorm:"foreignKey:MotherID"`
	Children    []Person  `json:"children,omitempty" gorm:"foreignKey:FatherID;foreignKey:MotherID"`
	
	// Metadata
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

type PersonRequest struct {
	FirstName string     `json:"first_name" binding:"required"`
	LastName  string     `json:"last_name" binding:"required"`
	Gender    string     `json:"gender"`
	BirthDate *time.Time `json:"birth_date"`
	DeathDate *time.Time `json:"death_date"`
	PhotoURL  string     `json:"photo_url"`
	Biography string     `json:"biography"`
	FatherID  *uint      `json:"father_id"`
	MotherID  *uint      `json:"mother_id"`
}

type FamilyTreeNode struct {
	Person   Person           `json:"person"`
	Children []FamilyTreeNode `json:"children,omitempty"`
}

// Enhanced family tree structures
type House struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"not null"` // e.g., "House Johnson", "House Williams"
	Description string    `json:"description"`
	FoundedYear int       `json:"founded_year"`
	CoatOfArms  string    `json:"coat_of_arms_url"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type FamilyRelationship struct {
	ID           uint   `json:"id" gorm:"primaryKey"`
	PersonID     uint   `json:"person_id"`
	RelatedToID  uint   `json:"related_to_id"`
	Relationship string `json:"relationship"` // father, mother, spouse, child, sibling
	Person       Person `json:"person" gorm:"foreignKey:PersonID"`
	RelatedTo    Person `json:"related_to" gorm:"foreignKey:RelatedToID"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type ComprehensiveFamilyTree struct {
	Houses        []HouseWithMembers    `json:"houses"`
	Relationships []FamilyRelationship  `json:"relationships"`
	Generations   []Generation          `json:"generations"`
	Statistics    FamilyStatistics      `json:"statistics"`
}

type HouseWithMembers struct {
	House   House    `json:"house"`
	Members []Person `json:"members"`
}

type Generation struct {
	Level   int      `json:"level"` // 1 = founders, 2 = children, 3 = grandchildren
	Name    string   `json:"name"`  // "Founders", "Second Generation", etc.
	People  []Person `json:"people"`
}

type FamilyStatistics struct {
	TotalPeople    int            `json:"total_people"`
	TotalHouses    int            `json:"total_houses"`
	Generations    int            `json:"generations"`
	LivingMembers  int            `json:"living_members"`
	DeceasedMembers int           `json:"deceased_members"`
	GenderCount    map[string]int `json:"gender_count"`
	AgeDistribution map[string]int `json:"age_distribution"`
}

type PersonWithRelationships struct {
	Person        Person               `json:"person"`
	Father        *Person              `json:"father,omitempty"`
	Mother        *Person              `json:"mother,omitempty"`
	Spouse        *Person              `json:"spouse,omitempty"`
	Children      []Person             `json:"children"`
	Siblings      []Person             `json:"siblings"`
	Grandparents  []Person             `json:"grandparents,omitempty"`
	Grandchildren []Person             `json:"grandchildren,omitempty"`
	Relationships []FamilyRelationship `json:"all_relationships"`
} 
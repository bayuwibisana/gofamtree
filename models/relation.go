package models

import "time"

type Relation struct {
	ID           uint      `json:"id" gorm:"primaryKey;column:id"`
	HouseID      uint      `json:"house_id" gorm:"not null;column:house_id"`
	PersonID     uint      `json:"person_id" gorm:"not null;column:person_id"`
	RelatedToID  uint      `json:"related_to_id" gorm:"not null;column:related_to_id"`
	RelationType string    `json:"relation_type" gorm:"type:text;column:relation_type"`
	CreatedAt    time.Time `json:"created_at" gorm:"column:created_at"`
	
	// Relationships (without foreign key constraints in GORM since we manage them manually)
	House     House  `json:"house" gorm:"foreignKey:HouseID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Person    Person `json:"person" gorm:"foreignKey:PersonID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	RelatedTo Person `json:"related_to" gorm:"foreignKey:RelatedToID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

// Add unique constraint at the table level
// This will be handled in the migration

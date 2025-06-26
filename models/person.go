package models

import (
	"time"
)

type Person struct {
	ID          uint       `json:"id" gorm:"primaryKey;column:id"`
	HouseID     uint       `json:"house_id" gorm:"not null;column:house_id"`
	Name        string     `json:"name" gorm:"not null;column:name"`
	Contact     string     `json:"contact" gorm:"column:contact"`
	Description string     `json:"description" gorm:"column:description"`
	Gender      string     `json:"gender" gorm:"type:text;column:gender"`
	DOB         *time.Time `json:"dob" gorm:"type:date;column:dob"`
	CreatedAt   time.Time  `json:"created_at" gorm:"column:created_at"`
	
	// Relationships
	House House `json:"house" gorm:"foreignKey:HouseID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

// TableName explicitly sets the table name for GORM
func (Person) TableName() string {
	return "persons"
}

package models

import "time"

type House struct {
	ID        uint      `json:"id" gorm:"primaryKey;column:id"`
	Name      string    `json:"name" gorm:"not null;column:name"`
	CreatedBy uint      `json:"created_by" gorm:"not null;column:created_by"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at"`
	
	// Relationships
	Admin   Admin    `json:"admin" gorm:"foreignKey:CreatedBy;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Persons []Person `json:"persons,omitempty" gorm:"foreignKey:HouseID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

package models

import (
	"time"
)

type Admin struct {
	ID        uint      `json:"id" gorm:"primaryKey;column:id"`
	Username  string    `json:"username" gorm:"unique;not null;column:username"`
	Password  string    `json:"-" gorm:"not null;column:password"` // Hashed, hidden from JSON
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at"`
}

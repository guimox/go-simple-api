package models

import "time"

type Locker struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Number    string    `json:"number"`  // E.g., "123A"
	Status    string    `json:"status"`  // "available", "in-use", etc.
	UserID    *uint     `json:"user_id"` // Foreign key to User
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

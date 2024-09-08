package models

import "time"

type User struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`                 // Never include password in JSON responses
	Lockers   []Locker  `json:"lockers,omitempty"` // One-to-many relation with lockers
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

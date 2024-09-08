package repository

import "go-simple-api/internal/models"

type UserRepository interface {
	AddUser(user *models.User) error
	GetUserByID(id int) (*models.User, error)
	UpdateUser(user *models.User) error
	DeleteUser(id int) error
}

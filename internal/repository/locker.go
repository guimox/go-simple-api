package repository

import "go-simple-api/internal/models"

type LockerRepository interface {
	AddLocker(locker *models.Locker) error
	GetLockerByID(id int) (*models.Locker, error)
	UpdateLocker(locker *models.Locker) error
	DeleteLocker(id int) error
}

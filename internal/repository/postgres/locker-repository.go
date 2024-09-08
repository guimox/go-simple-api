package postgres

import (
	"database/sql"
	"go-simple-api/internal/models"
	"go-simple-api/internal/repository"
)

type LockerRepositorySQL struct {
	db *sql.DB
}

func NewLockerRepositorySQL(db *sql.DB) repository.LockerRepository {
	return &LockerRepositorySQL{db: db}
}

func (r *LockerRepositorySQL) AddLocker(locker *models.Locker) error {
	_, err := r.db.Exec("INSERT INTO lockers (number, status, user_id) VALUES ($1, $2, $3)", locker.Number, locker.Status, locker.UserID)
	return err
}

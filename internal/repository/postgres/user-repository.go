// repository/postgres/user_repository.go
package postgres

import (
	"database/sql"
	"myapp/models"
	"myapp/repository"
)

type UserRepositorySQL struct {
	db *sql.DB
}

func NewUserRepositorySQL(db *sql.DB) repository.UserRepository {
	return &UserRepositorySQL{db: db}
}

func (r *UserRepositorySQL) AddUser(user *models.User) error {
	_, err := r.db.Exec("INSERT INTO users (email, first_name, last_name, password) VALUES ($1, $2, $3, $4)", user.Email, user.FirstName, user.LastName, user.Password)
	return err
}

func (r *UserRepositorySQL) GetUserByID(id int) (*models.User, error) {
	user := &models.User{}
	err := r.db.QueryRow("SELECT id, email, first_name, last_name FROM users WHERE id = $1", id).Scan(&user.ID, &user.Email, &user.FirstName, &user.LastName)
	if err != nil {
		return nil, err
	}
	return user, nil
}

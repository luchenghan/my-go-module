package repository

import (
	"database/sql"
	"mymodule/gin/internal/db"
	"mymodule/gin/internal/models"
	"mymodule/gin/pkg/logger"
)

type UsersRepository struct {
	db *sql.DB
}

func NewUsersRepository() *UsersRepository {
	repo := new(UsersRepository)
	repo.db = db.GetDB()
	return repo
}

func (r *UsersRepository) GetUserByEmail(email string) (*models.User, error) {
	query := `
        SELECT id, name, email, password
        FROM users
        WHERE email = ?
    `

	user := &models.User{}
	err := r.db.QueryRow(query, email).Scan(&user.ID, &user.Name, &user.Email, &user.Password)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UsersRepository) GetUserByID(id string) (*models.User, error) {
	query := `
		SELECT id, name, email
		FROM users
		WHERE id = ?
	`

	user := &models.User{}
	err := r.db.QueryRow(query, id).Scan(&user.ID, &user.Name, &user.Email)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UsersRepository) GetAllUsers() ([]models.User, error) {
	var users []models.User
	query := `
		SELECT id, name, email
		FROM users
	`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user models.User
		if err = rows.Scan(&user.ID, &user.Name, &user.Email); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (r *UsersRepository) SaveUser(user models.User) error {
	// Save user to the database
	query := `
		INSERT INTO users (name, email, password)
		VALUES (?, ?, ?)
	`
	_, err := r.db.Exec(query, user.Name, user.Email, user.Password)
	if err != nil {
		logger.Errorf("Failed to save user: %v", err)
		return err
	}

	return nil
}

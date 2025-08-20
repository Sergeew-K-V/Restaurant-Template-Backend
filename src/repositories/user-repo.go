package repositories

import (
	"database/sql"
	"fmt"
	"log"
	"restaurant-backend/src/models"
	"time"

	"github.com/google/uuid"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db}
}

func (ur *UserRepository) CreateUser(user *models.User) error {
	user.Id = uuid.New()

	query := `
		INSERT INTO users (name, email, password, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id`

	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	err := ur.db.QueryRow(query, user.Name, user.Email, user.Password, user.CreatedAt, user.UpdatedAt).Scan(&user.Id)

	if err != nil {
		log.Printf("ERROR: Failed to create user: %v", err)
		return fmt.Errorf("error creating user: %v", err)
	}

	return nil
}

func (ur *UserRepository) GetUserById(id string) (*models.User, error) {
	query := `SELECT id, name, email, password, created_at, updated_at FROM users WHERE email = $1`

	user := &models.User{}

	err := ur.db.QueryRow(query, id).Scan(&user.Id, &user.Name, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		log.Printf("ERROR: Failed to get user by id: %v", err)
		return nil, fmt.Errorf("error getting user by id: %v", err)
	}

	return user, nil

}

func (ur *UserRepository) UserExists(email string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`

	var exists bool

	err := ur.db.QueryRow(query, email).Scan(&exists)
	if err != nil {
		log.Printf("ERROR: Failed to check if user exists: %v", err)
		return false, fmt.Errorf("error checking if user exists: %v", err)
	}

	return exists, nil
}

func (ur *UserRepository) GetUserByEmail(email string) (*models.User, error) {
	query := `SELECT * FROM users WHERE email = $1`

	user := &models.User{}

	err := ur.db.QueryRow(query, email).Scan(&user.Id, &user.Name, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		log.Printf("ERROR: Failed to get user by email: %v", err)
		return nil, fmt.Errorf("error getting user by email: %v", err)
	}

	return user, nil
}

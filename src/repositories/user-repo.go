package repositories

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
)

type UserRepository struct {
	db *sql.DB
}

type User struct {
	Id        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db}
}

func (ur *UserRepository) CreateUser(user *User) error {
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

func (ur *UserRepository) GetUserById(id string) (*User, error) {
	query := `SELECT id, name, email, password, created_at, updated_at FROM users WHERE email = $1`

	user := &User{}

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

func (ur *UserRepository) GetUserByEmail(email string) (*User, error) {
	query := `SELECT * FROM users WHERE email = $1`

	user := &User{}

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

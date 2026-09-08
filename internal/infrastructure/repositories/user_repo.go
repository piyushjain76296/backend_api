package repositories

import (
	"database/sql"
	"errors"

	"ticket-system/internal/core/models"
)

var ErrUserNotFound = errors.New("user not found")
var ErrDuplicateEmail = errors.New("user with this email already exists")

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *models.User) error {
	query := `INSERT INTO users (email, password_hash) VALUES (?, ?)`
	result, err := r.db.Exec(query, user.Email, user.PasswordHash)
	if err != nil {
		if err.Error() == "UNIQUE constraint failed: users.email" {
			return ErrDuplicateEmail
		}
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	user.ID = int(id)
	return nil
}

func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	query := `SELECT id, email, password_hash, created_at FROM users WHERE email = ?`
	row := r.db.QueryRow(query, email)

	var user models.User
	err := row.Scan(&user.ID, &user.Email, &user.PasswordHash, &user.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

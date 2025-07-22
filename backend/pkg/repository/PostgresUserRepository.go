package repository

import (
	"backend/pkg/models"
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"strings"
)

type PostgresUserRepository struct {
	db *sql.DB
}

func NewPostgresUserRepository(db *sql.DB) (*PostgresUserRepository, error) {
	repo := &PostgresUserRepository{db: db}
	return repo, nil
}

func (r *PostgresUserRepository) GetById(id string) (*models.User, error) {
	user := &models.User{}
	err := r.db.QueryRow(`
		SELECT id, username, password_hash 
		FROM users 
		WHERE id = $1`, id).Scan(
		&user.ID,
		&user.Username,
		&user.Password,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("User not found")
		}
		return nil, err
	}

	return user, nil
}

func (r *PostgresUserRepository) GetByUsername(username, password string) (*models.User, error) {
	user := &models.User{}
	err := r.db.QueryRow(`
		SELECT id, username, password_hash
		FROM users 
		WHERE username = $1 AND password_hash = $2`, username, password).Scan(
		&user.ID,
		&user.Username,
		&user.Password,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("User not found")
		}
		return nil, err
	}

	return user, nil
}

func (r *PostgresUserRepository) Add(user *models.User) error {

	//err := r.db.QueryRow(
	//	`SELECT id, username FROM users WHERE username=$1`, user.Username).Scan(
	//		&created.ID,
	//		&created.Username,
	//	)
	//if err != nil {
	//	return err
	//}
	if strings.TrimSpace(user.ID) == "" {
		id, err := uuid.NewRandom()
		if err != nil {
			return err
		}
		user.ID = id.String()
	}
	_, err := r.db.Exec(`INSERT into users (id, username, password_hash) 
VALUES ($1, $2, $3)`, &user.ID,
		&user.Username,
		&user.Password,
	)
	if err != nil {
		return err
	}
	return nil
}

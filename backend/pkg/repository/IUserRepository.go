package repository

import "backend/pkg/models"

type IUserRepository interface {
	GetById(id string) (*models.User, error)
	GetByUsername(username, password string) (*models.User, error)
	Add(user *models.User) error
}

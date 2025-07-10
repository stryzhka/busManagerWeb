package repository

import "backend/pkg/models"

type IBusRepository interface {
	GetById(id string) (*models.Bus, error)
	GetByNumber(number string) (*models.Bus, error)
	Add(bus *models.Bus) error
	DeleteById(id string) error
	GetAll() ([]models.Bus, error)
	UpdateById(bus *models.Bus) error
}

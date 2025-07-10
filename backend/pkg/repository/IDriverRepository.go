package repository

import "backend/pkg/models"

type IDriverRepository interface {
	GetById(id string) (*models.Driver, error)
	GetByPassportSeries(passportSeries string) (*models.Driver, error)
	Add(driver *models.Driver) error
	DeleteById(id string) error
	GetAll() ([]models.Driver, error)
	UpdateById(driver *models.Driver) error
}

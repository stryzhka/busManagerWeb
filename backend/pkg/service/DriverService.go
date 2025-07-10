package service

import (
	"backend/pkg/models"
	"backend/pkg/repository"
	"errors"
)

type DriverService struct {
	repo repository.IDriverRepository
}

func NewDriverService(r repository.IDriverRepository) *DriverService {
	b := &DriverService{r}
	return b
}

func (ds DriverService) GetById(id string) (*models.Driver, error) {

	driver, err := ds.repo.GetById(id)
	if err != nil {
		return nil, err
	}
	if driver == nil {
		return nil, errors.New("Driver not found")
	}
	return driver, nil
}

func (ds DriverService) GetByPassportSeries(series string) (*models.Driver, error) {
	driver, err := ds.repo.GetByPassportSeries(series)
	if err != nil {
		return nil, err
	}
	if driver == nil {
		return nil, errors.New("Driver not found")
	}
	return driver, nil
}

func (ds DriverService) Add(driver *models.Driver) error {
	err := ds.repo.Add(driver)
	return err
}

func (ds DriverService) GetAll() []models.Driver {
	var m []models.Driver
	m, _ = ds.repo.GetAll()
	return m

}

func (ds DriverService) DeleteById(id string) error {
	err := ds.repo.DeleteById(id)
	return err
}

func (ds DriverService) UpdateById(driver *models.Driver) error {
	err := ds.repo.UpdateById(driver)
	return err
}

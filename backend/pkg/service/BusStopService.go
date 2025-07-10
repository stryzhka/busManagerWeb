package service

import (
	"backend/pkg/models"
	"backend/pkg/repository"
	"errors"
)

type BusStopService struct {
	repo repository.IBusStopRepository
}

func NewBusStopService(r repository.IBusStopRepository) *BusStopService {
	b := &BusStopService{r}
	return b
}

func (ds BusStopService) GetById(id string) (*models.BusStop, error) {

	busStop, err := ds.repo.GetById(id)
	if err != nil {
		return nil, err
	}
	if busStop == nil {
		return nil, errors.New("Bus stop not found")
	}
	return busStop, nil
}

func (ds BusStopService) GetByName(name string) (*models.BusStop, error) {
	busStop, err := ds.repo.GetByName(name)
	if err != nil {
		return nil, err
	}
	if busStop == nil {
		return nil, errors.New("Bus stop not found")
	}
	return busStop, nil
}

func (ds BusStopService) Add(busStop *models.BusStop) error {
	err := ds.repo.Add(busStop)
	return err
}

func (ds BusStopService) GetAll() ([]models.BusStop, error) {
	var m []models.BusStop
	m, err := ds.repo.GetAll()
	if err != nil {
		return nil, err
	}
	return m, nil

}

func (ds BusStopService) DeleteById(id string) error {
	err := ds.repo.DeleteById(id)
	return err
}

func (ds BusStopService) UpdateById(busStop *models.BusStop) error {
	err := ds.repo.UpdateById(busStop)
	return err
}

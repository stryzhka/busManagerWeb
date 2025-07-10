package service

import (
	"backend/pkg/models"
	"backend/pkg/repository"
	"errors"
)

type BusService struct {
	repo repository.IBusRepository
}

func NewBusService(r repository.IBusRepository) *BusService {
	b := &BusService{r}
	return b
}

func (bs BusService) GetById(id string) (*models.Bus, error) {

	bus, err := bs.repo.GetById(id)
	if err != nil {
		return nil, err
	}
	if bus == nil {
		return nil, errors.New("Bus not found")
	}
	return bus, nil
}

func (bs BusService) GetByNumber(number string) (*models.Bus, error) {
	bus, err := bs.repo.GetByNumber(number)
	if err != nil {
		return nil, err
	}
	if bus == nil {
		return nil, errors.New("Bus not found")
	}
	return bus, nil
}

func (bs BusService) Add(bus *models.Bus) error {
	err := bs.repo.Add(bus)
	return err
}

func (bs BusService) GetAll() []models.Bus {
	var m []models.Bus
	m, _ = bs.repo.GetAll()
	return m

}

func (bs BusService) DeleteById(id string) error {
	err := bs.repo.DeleteById(id)
	return err
}

func (bs BusService) UpdateById(bus *models.Bus) error {
	err := bs.repo.UpdateById(bus)
	return err
}

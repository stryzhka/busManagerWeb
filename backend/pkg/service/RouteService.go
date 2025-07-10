package service

import (
	"backend/pkg/models"
	"backend/pkg/repository"

	"errors"
)

type RouteService struct {
	repo        repository.IRouteRepository
	driverRepo  repository.IDriverRepository
	busRepo     repository.IBusRepository
	busStopRepo repository.IBusStopRepository
}

func NewRouteService(
	r repository.IRouteRepository,
	driverRepo repository.IDriverRepository,
	busRepo repository.IBusRepository,
	busStopRepo repository.IBusStopRepository,
) *RouteService {
	b := &RouteService{r, driverRepo, busRepo, busStopRepo}
	return b
}

func (rs RouteService) GetById(id string) (*models.Route, error) {

	route, err := rs.repo.GetById(id)
	if err != nil {
		return nil, err
	}
	if route == nil {
		return nil, errors.New("Route not found")
	}
	return route, nil
}

func (rs RouteService) GetByNumber(number string) (*models.Route, error) {
	route, err := rs.repo.GetByNumber(number)
	if err != nil {
		return nil, err
	}
	if route == nil {
		return nil, errors.New("Route not found")
	}
	return route, nil
}

func (rs RouteService) Add(route *models.Route) error {
	err := rs.repo.Add(route)
	return err
}

func (rs RouteService) GetAll() ([]models.Route, error) {
	var m []models.Route
	m, err := rs.repo.GetAll()
	if err != nil {
		return nil, err
	}
	return m, nil

}

func (rs RouteService) DeleteById(id string) error {
	err := rs.repo.DeleteById(id)
	return err
}

func (rs RouteService) UpdateById(route *models.Route) error {
	err := rs.repo.UpdateById(route)
	return err
}

func (rs RouteService) AssignDriver(routeId, driverId string) error {
	route, err := rs.GetById(routeId)
	if route == nil {
		return errors.New("Route not found")
	}
	if err != nil {
		return err
	}

	driver, err := rs.driverRepo.GetById(driverId)
	if driver == nil {
		return errors.New("Driver not found")
	}
	if err != nil {
		return err
	}
	err = rs.repo.AssignDriver(routeId, driverId)
	if err != nil {
		return err
	}
	return nil
}

func (rs RouteService) AssignBusStop(routeId, busStopId string) error {
	route, err := rs.GetById(routeId)
	if route == nil {
		return errors.New("Route not found")
	}
	if err != nil {
		return err
	}

	busStop, err := rs.busStopRepo.GetById(busStopId)
	if busStop == nil {
		return errors.New("Bus stop not found")
	}
	if err != nil {
		return err
	}
	err = rs.repo.AssignBusStop(routeId, busStopId)
	if err != nil {
		return err
	}
	return nil
}

func (rs RouteService) AssignBus(routeId, busId string) error {
	route, err := rs.GetById(routeId)
	if route == nil {
		return errors.New("Route not found")
	}
	if err != nil {
		return err
	}

	bus, err := rs.busRepo.GetById(busId)
	if bus == nil {
		return errors.New("Bus not found")
	}
	if err != nil {
		return err
	}
	err = rs.repo.AssignBus(routeId, busId)
	if err != nil {
		return err
	}
	return nil
}

func (rs RouteService) UnassignDriver(routeId, driverId string) error {
	route, err := rs.GetById(routeId)
	if route == nil {
		return errors.New("Route not found")
	}
	if err != nil {
		return err
	}

	driver, err := rs.driverRepo.GetById(driverId)
	if driver == nil {
		return errors.New("Driver not found")
	}
	if err != nil {
		return err
	}
	err = rs.repo.UnassignDriver(routeId, driverId)
	if err != nil {
		return err
	}
	return nil
}

func (rs RouteService) UnassignBusStop(routeId, busStopId string) error {
	route, err := rs.GetById(routeId)
	if route == nil {
		return errors.New("Route not found")
	}
	if err != nil {
		return err
	}

	busStop, err := rs.busStopRepo.GetById(busStopId)
	if busStop == nil {
		return errors.New("Bus stop not found")
	}
	if err != nil {
		return err
	}
	err = rs.repo.UnassignBusStop(routeId, busStopId)
	if err != nil {
		return err
	}
	return nil
}

func (rs RouteService) UnassignBus(routeId, busId string) error {
	route, err := rs.GetById(routeId)
	if route == nil {
		return errors.New("Route not found")
	}
	if err != nil {
		return err
	}

	bus, err := rs.busRepo.GetById(busId)
	if bus == nil {
		return errors.New("Bus not found")
	}
	if err != nil {
		return err
	}
	err = rs.repo.UnassignBus(routeId, busId)
	if err != nil {
		return err
	}
	return nil
}

func (rs RouteService) GetAllDriversById(routeId string) ([]models.Driver, error) {
	route, err := rs.GetById(routeId)
	if route == nil {
		return nil, errors.New("Route not found")
	}
	if err != nil {
		return nil, err
	}

	drivers, err := rs.repo.GetAllDriversById(routeId)
	if err != nil {
		return nil, err
	}
	if drivers == nil {
		return drivers, errors.New("Drivers not found")
	}
	return drivers, nil
}

func (rs RouteService) GetAllBusStopsById(routeId string) ([]models.BusStop, error) {
	route, err := rs.GetById(routeId)
	if route == nil {
		return nil, errors.New("Route not found")
	}
	if err != nil {
		return nil, err
	}

	busStops, err := rs.repo.GetAllBusStopsById(routeId)
	if err != nil {
		return nil, err
	}
	if busStops == nil {
		return busStops, errors.New("Bus stops not found")
	}
	return busStops, nil
}

func (rs RouteService) GetAllBusesById(routeId string) ([]models.Bus, error) {
	route, err := rs.GetById(routeId)
	if route == nil {
		return nil, errors.New("Route not found")
	}
	if err != nil {
		return nil, err
	}

	buses, err := rs.repo.GetAllBusesById(routeId)
	if err != nil {
		return nil, err
	}
	if buses == nil {
		return buses, errors.New("Buses not found")
	}
	return buses, nil
}

package service

import "backend/pkg/models"

type IRouteService interface {
	GetById(id string) (*models.Route, error)
	GetByNumber(number string) (*models.Route, error)
	Add(route *models.Route) error
	DeleteById(id string) error
	GetAll() ([]models.Route, error)
	UpdateById(route *models.Route) error
	AssignDriver(routeId, driverId string) error
	AssignBusStop(routeId, busStopId string) error
	AssignBus(routeId, busId string) error
	UnassignDriver(routeId, driverId string) error
	UnassignBusStop(routeId, busStopId string) error
	UnassignBus(routeId, busId string) error
	GetAllDriversById(routeId string) ([]models.Driver, error)
	GetAllBusStopsById(routeId string) ([]models.BusStop, error)
	GetAllBusesById(routeId string) ([]models.Bus, error)
	// TODO: getall for all models, unassign
}

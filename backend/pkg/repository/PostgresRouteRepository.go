package repository

import (
	"backend/pkg/models"
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"strings"
)

type PostgresRouteRepository struct {
	db *sql.DB
}

func NewPostgresRouteRepository(db *sql.DB) (*PostgresRouteRepository, error) {
	repo := &PostgresRouteRepository{db: db}
	return repo, nil
}

func (r *PostgresRouteRepository) GetById(id string) (*models.Route, error) {
	route := &models.Route{}
	err := r.db.QueryRow(`
		SELECT id, number 
		FROM routes 
		WHERE id = $1`, id).Scan(
		&route.ID,
		&route.Number,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("Route not found")
		}
		return nil, err
	}

	return route, nil
}

func (r *PostgresRouteRepository) GetByNumber(number string) (*models.Route, error) {
	route := &models.Route{}
	err := r.db.QueryRow(`
		SELECT id, number
		FROM routes 
		WHERE number = $1`, number).Scan(
		&route.ID,
		&route.Number,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("Route not found")
		}
		return nil, err
	}

	return route, nil
}

func (r *PostgresRouteRepository) Add(route *models.Route) error {
	exist, err := r.GetByNumber(route.Number)
	if exist != nil {
		return errors.New("Route already exists")
	}
	if strings.TrimSpace(route.ID) == "" {
		id, err := uuid.NewRandom()
		if err != nil {
			return err
		}
		route.ID = id.String()
	}
	_, err = r.db.Exec(`INSERT into routes (id, number) 
VALUES ($1, $2)`, &route.ID,
		&route.Number,
	)
	if err != nil {
		return err
	}
	return nil
}

func (r *PostgresRouteRepository) GetAll() ([]models.Route, error) {
	var routes []models.Route
	rows, err := r.db.Query(`
		SELECT id, number
		FROM routes 
		`)
	if err != nil {
		//if err == sql.ErrNoRows {
		//	return nil, errors.New("Empty DB")
		//}
		return nil, err
	}
	for rows.Next() {
		route := &models.Route{}
		err := rows.Scan(
			&route.ID,
			&route.Number,
		)
		if err != nil {
			return nil, err
		}
		routes = append(routes, *route)
	}
	return routes, nil
}

func (r *PostgresRouteRepository) DeleteById(id string) error {
	exist, err := r.GetById(id)
	if exist == nil {
		return errors.New("Route not found")
	}
	if err != nil {
		return err
	}
	_, err = r.db.Exec("DELETE FROM routes WHERE id = $1", id)
	if err != nil {
		return err
	}
	return nil
}

func (r *PostgresRouteRepository) UpdateById(route *models.Route) error {

	exist, err := r.GetById(route.ID)
	if exist == nil {
		return errors.New("Route not found")
	}
	if err != nil {
		return err
	}
	_, err = r.db.Exec("UPDATE routes SET number = $1 WHERE id = $2",
		route.Number, route.ID)
	if err != nil {
		return err
	}
	return nil
}

func (r *PostgresRouteRepository) AssignDriver(routeId, driverId string) error {
	exist, err := r.GetById(routeId)
	if exist == nil {
		return errors.New("Route not found")
	}
	var count int
	err = r.db.QueryRow(`SELECT COUNT(*) FROM routes_drivers WHERE route_id = $1 AND driver_id = $2`, routeId, driverId).Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("Pair route_id and driver_id already exists")
	}
	_, err = r.db.Exec(`INSERT into routes_drivers (route_id, driver_id) 
VALUES ($1, $2)`, routeId,
		driverId,
	)
	if err != nil {
		return err
	}
	return nil
}

func (r *PostgresRouteRepository) AssignBusStop(routeId, busStopId string) error {
	exist, err := r.GetById(routeId)
	if exist == nil {
		return errors.New("Route not found")
	}
	var count int
	err = r.db.QueryRow(`SELECT COUNT(*) FROM routes_bus_stops WHERE route_id = $1 AND bus_stop_id = $2`, routeId, busStopId).Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("Pair route_id and bus_stop_id already exists")
	}
	_, err = r.db.Exec(`INSERT into routes_bus_stops (route_id, bus_stop_id) 
VALUES ($1, $2)`, routeId,
		busStopId,
	)
	if err != nil {
		return err
	}
	return nil
}

func (r *PostgresRouteRepository) AssignBus(routeId, busId string) error {
	exist, err := r.GetById(routeId)
	if exist == nil {
		return errors.New("Route not found")
	}
	var count int
	err = r.db.QueryRow(`SELECT COUNT(*) FROM routes_buses WHERE route_id = $1 AND bus_id = $2`, routeId, busId).Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("Pair route_id and bus_id already exists")
	}
	_, err = r.db.Exec(`INSERT into routes_buses (route_id, bus_id) 
VALUES ($1, $2)`, routeId,
		busId,
	)
	if err != nil {
		return err
	}
	return nil
}

func (r *PostgresRouteRepository) UnassignBusStop(routeId, busStopId string) error {
	exist, err := r.GetById(routeId)
	if exist == nil {
		return errors.New("Route not found")
	}
	_, err = r.db.Exec(`DELETE FROM routes_bus_stops WHERE route_id = $1 AND bus_stop_id = $2`, routeId, busStopId)
	if err != nil {
		return err
	}
	return nil
}

func (r *PostgresRouteRepository) UnassignBus(routeId, busId string) error {
	exist, err := r.GetById(routeId)
	if exist == nil {
		return errors.New("Route not found")
	}
	_, err = r.db.Exec(`DELETE FROM routes_buses WHERE route_id = $1 AND bus_id = $2`, routeId, busId)
	if err != nil {
		return err
	}
	return nil
}

func (r *PostgresRouteRepository) UnassignDriver(routeId, driverId string) error {
	exist, err := r.GetById(routeId)
	if exist == nil {
		return errors.New("Route not found")
	}
	_, err = r.db.Exec(`DELETE FROM routes_drivers WHERE route_id = $1 AND driver_id = $2`, routeId, driverId)
	if err != nil {
		return err
	}
	return nil
}

func (r *PostgresRouteRepository) GetAllDriversById(routeId string) ([]models.Driver, error) {
	var drivers []models.Driver
	exist, err := r.GetById(routeId)
	if exist == nil {
		return nil, errors.New("Route not found")
	}
	if err != nil {
		return nil, err
	}
	rows, err := r.db.Query(`
		SELECT d.id, d.name, d.surname, d.patronymic, d.birth_date, d.passport_series, d.snils, d.license_series
		FROM drivers d 
		JOIN routes_drivers rd ON d.id = rd.driver_id
		WHERE rd.route_id=$1
	`, routeId)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		driver := &models.Driver{}
		err := rows.Scan(
			&driver.ID,
			&driver.Name,
			&driver.Surname,
			&driver.Patronymic,
			&driver.BirthDate,
			&driver.PassportSeries,
			&driver.Snils,
			&driver.LicenseSeries,
		)
		if err != nil {
			return nil, err
		}
		drivers = append(drivers, *driver)
	}
	return drivers, nil
}

func (r *PostgresRouteRepository) GetAllBusStopsById(routeId string) ([]models.BusStop, error) {
	var busStops []models.BusStop
	exist, err := r.GetById(routeId)
	if exist == nil {
		return nil, errors.New("Route not found")
	}
	if err != nil {
		return nil, err
	}
	rows, err := r.db.Query(`
		SELECT d.id, d.lat, d.long, d.name
		FROM bus_stops d 
		JOIN routes_bus_stops rd ON d.id = rd.bus_stop_id
		WHERE rd.route_id=$1
	`, routeId)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		busStop := &models.BusStop{}
		err := rows.Scan(
			&busStop.ID,
			&busStop.Lat,
			&busStop.Long,
			&busStop.Name,
		)
		if err != nil {
			return nil, err
		}
		busStops = append(busStops, *busStop)
	}
	return busStops, nil
}

func (r *PostgresRouteRepository) GetAllBusesById(routeId string) ([]models.Bus, error) {
	var buses []models.Bus
	exist, err := r.GetById(routeId)
	if exist == nil {
		return nil, errors.New("Route not found")
	}
	if err != nil {
		return nil, err
	}
	rows, err := r.db.Query(`
		SELECT d.id, d.brand, d.bus_model, d.register_number, d.assembly_date, d.last_repair_date
		FROM buses d 
		JOIN routes_buses rd ON d.id = rd.bus_id
		WHERE rd.route_id=$1
	`, routeId)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		bus := &models.Bus{}
		err := rows.Scan(
			&bus.ID,
			&bus.Brand,
			&bus.BusModel,
			&bus.RegisterNumber,
			&bus.AssemblyDate,
			&bus.LastRepairDate,
		)
		if err != nil {
			return nil, err
		}
		buses = append(buses, *bus)
	}
	return buses, nil
}

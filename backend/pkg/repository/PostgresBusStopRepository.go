package repository

import (
	"backend/pkg/models"
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"strings"
)

type PostgresBusStopRepository struct {
	db *sql.DB
}

func NewPostgresBusStopRepository(db *sql.DB) (*PostgresBusStopRepository, error) {
	repo := &PostgresBusStopRepository{db: db}
	return repo, nil
}

func (r *PostgresBusStopRepository) GetById(id string) (*models.BusStop, error) {
	stop := &models.BusStop{}
	err := r.db.QueryRow(`
		SELECT id, lat, long, name 
		FROM bus_stops 
		WHERE id = $1`, id).Scan(
		&stop.ID,
		&stop.Lat,
		&stop.Long,
		&stop.Name,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("Bus stop not found")
		}
		return nil, err
	}

	return stop, nil
}

func (r *PostgresBusStopRepository) GetByName(name string) (*models.BusStop, error) {
	stop := &models.BusStop{}
	err := r.db.QueryRow(`
		SELECT id, lat, long, name 
		FROM bus_stops 
		WHERE name = $4`, name).Scan(
		&stop.ID,
		&stop.Lat,
		&stop.Long,
		&stop.Name,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("Bus stop not found")
		}
		return nil, err
	}

	return stop, nil
}

func (r *PostgresBusStopRepository) Add(busStop *models.BusStop) error {
	exist, err := r.GetByName(busStop.Name)
	if exist != nil {
		return errors.New("Bus stop already exists")
	}
	if strings.TrimSpace(busStop.ID) == "" {
		id, err := uuid.NewRandom()
		if err != nil {
			return err
		}
		busStop.ID = id.String()
	}
	_, err = r.db.Exec(`INSERT into bus_stops
    (id, lat, long, name ) 
VALUES ($1, $2, $3, $4)`,
		&busStop.ID,
		&busStop.Lat,
		&busStop.Long,
		&busStop.Name)
	if err != nil {
		return err
	}
	return nil
}

func (r *PostgresBusStopRepository) GetAll() ([]models.BusStop, error) {
	var busStops []models.BusStop
	rows, err := r.db.Query(`
		SELECT id, lat, long, name
		FROM bus_stops 
		`)
	if err != nil {
		//if err == sql.ErrNoRows {
		//	return nil, errors.New("Empty DB")
		//}
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

func (r *PostgresBusStopRepository) DeleteById(id string) error {
	exist, err := r.GetById(id)
	if exist == nil {
		return errors.New("Bus stop not found")
	}
	if err != nil {
		return err
	}
	_, err = r.db.Exec("DELETE FROM bus_stops WHERE id = $1", id)
	if err != nil {
		return err
	}
	return nil
}

func (r *PostgresBusStopRepository) UpdateById(busStop *models.BusStop) error {
	exist, err := r.GetById(busStop.ID)
	if exist == nil {
		return errors.New("Bus stop not found")
	}
	if err != nil {
		return err
	}
	_, err = r.db.Exec(`UPDATE bus_stops SET lat = $1, long = $2, name = $3 WHERE id = $4`,
		busStop.Lat,
		busStop.Long,
		busStop.Name,
		busStop.ID,
	)
	if err != nil {
		return err
	}
	return nil
}

package repository

import (
	"backend/pkg/models"
	"database/sql"
	"errors"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"strings"
)

type PostgresBusRepository struct {
	db *sql.DB
}

func NewPostgresBusRepository(db *sql.DB) (*PostgresBusRepository, error) {

	repo := &PostgresBusRepository{db: db}
	return repo, nil
}

func (r *PostgresBusRepository) GetById(id string) (*models.Bus, error) {
	bus := &models.Bus{}
	err := r.db.QueryRow(`
		SELECT id, brand, bus_model, register_number, assembly_date, last_repair_date 
		FROM buses 
		WHERE id = $1`, id).Scan(
		&bus.ID,
		&bus.Brand,
		&bus.BusModel,
		&bus.RegisterNumber,
		&bus.AssemblyDate,
		&bus.LastRepairDate,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("Bus not found")
		}
		return nil, err
	}

	return bus, nil
}

func (r *PostgresBusRepository) GetByNumber(number string) (*models.Bus, error) {
	bus := &models.Bus{}
	err := r.db.QueryRow(`
		SELECT id, brand, bus_model, register_number, assembly_date, last_repair_date 
		FROM buses 
		WHERE register_number = $1`, number).Scan(
		&bus.ID,
		&bus.Brand,
		&bus.BusModel,
		&bus.RegisterNumber,
		&bus.AssemblyDate,
		&bus.LastRepairDate,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("Bus not found")
		}
		return nil, err
	}

	return bus, nil
}

func (r *PostgresBusRepository) Add(bus *models.Bus) error {
	exist, err := r.GetByNumber(bus.RegisterNumber)
	if exist != nil {
		return errors.New("Bus already exists")
	}
	if strings.TrimSpace(bus.ID) == "" {
		id, err := uuid.NewRandom()
		if err != nil {
			return err
		}
		bus.ID = id.String()
	}
	_, err = r.db.Exec(`INSERT into buses (id, brand, bus_model, register_number, assembly_date, last_repair_date ) 
VALUES ($1, $2, $3, $4, $5, $6)`, &bus.ID,
		&bus.Brand,
		&bus.BusModel,
		&bus.RegisterNumber,
		&bus.AssemblyDate,
		&bus.LastRepairDate)
	if err != nil {
		return err
	}
	return nil

}

func (r *PostgresBusRepository) GetAll() ([]models.Bus, error) {

	var buses []models.Bus
	rows, err := r.db.Query(`
		SELECT id, brand, bus_model, register_number, assembly_date, last_repair_date 
		FROM buses 
		`)
	if err != nil {
		//if err == sql.ErrNoRows {
		//	return nil, errors.New("Empty DB")
		//}
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

func (r *PostgresBusRepository) DeleteById(id string) error {
	exist, err := r.GetById(id)
	if exist == nil {
		return errors.New("Bus not found")
	}
	if err != nil {
		return err
	}
	_, err = r.db.Exec("DELETE FROM buses WHERE id = $1", id)
	if err != nil {
		return err
	}
	return nil
}

func (r *PostgresBusRepository) UpdateById(bus *models.Bus) error {
	exist, err := r.GetById(bus.ID)
	if exist == nil {
		return errors.New("Bus not found")
	}
	if err != nil {
		return err
	}
	_, err = r.db.Exec("UPDATE buses SET brand = $1, bus_model = $2, register_number = $3, assembly_date = $4, last_repair_date = $5 WHERE id = $6", bus.Brand, bus.BusModel, bus.RegisterNumber, bus.AssemblyDate, bus.LastRepairDate, bus.ID)

	if err != nil {
		return err
	}
	return nil
}

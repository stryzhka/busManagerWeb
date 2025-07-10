package repository

import (
	"backend/pkg/models"
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"strings"
)

type PostgresDriverRepository struct {
	db *sql.DB
}

func NewPostgresDriverRepository(db *sql.DB) (*PostgresDriverRepository, error) {
	repo := &PostgresDriverRepository{db: db}
	return repo, nil
}

func (r *PostgresDriverRepository) GetById(id string) (*models.Driver, error) {
	driver := &models.Driver{}
	err := r.db.QueryRow(`
		SELECT id, name, surname, patronymic, birth_date, passport_series, snils, license_series 
		FROM drivers 
		WHERE id = $1`, id).Scan(
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
		if err == sql.ErrNoRows {
			return nil, errors.New("Driver not found")
		}
		return nil, err
	}

	return driver, nil
}

func (r *PostgresDriverRepository) GetByPassportSeries(series string) (*models.Driver, error) {
	driver := &models.Driver{}
	err := r.db.QueryRow(`
		SELECT id, name, surname, patronymic, birth_date, passport_series, snils, license_series
		FROM drivers 
		WHERE passport_series = $1`, series).Scan(
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
		if err == sql.ErrNoRows {
			return nil, errors.New("Driver not found")
		}
		return nil, err
	}

	return driver, nil
}

func (r *PostgresDriverRepository) Add(driver *models.Driver) error {
	exist, err := r.GetByPassportSeries(driver.PassportSeries)
	if exist != nil {
		return errors.New("Driver already exists")
	}
	if strings.TrimSpace(driver.ID) == "" {
		id, err := uuid.NewRandom()
		if err != nil {
			return err
		}
		driver.ID = id.String()
	}
	_, err = r.db.Exec(`INSERT into drivers (id, name, surname, patronymic, birth_date, passport_series, snils, license_series ) 
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`, &driver.ID,
		&driver.Name,
		&driver.Surname,
		&driver.Patronymic,
		&driver.BirthDate,
		&driver.PassportSeries,
		&driver.Snils,
		&driver.LicenseSeries)
	if err != nil {
		return err
	}
	return nil

}

func (r *PostgresDriverRepository) GetAll() ([]models.Driver, error) {
	var drivers []models.Driver
	rows, err := r.db.Query(`
		SELECT id, name, surname, patronymic, birth_date, passport_series, snils, license_series
		FROM drivers 
		`)
	if err != nil {
		//if err == sql.ErrNoRows {
		//	return nil, errors.New("Empty DB")
		//}
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

func (r *PostgresDriverRepository) DeleteById(id string) error {
	exist, err := r.GetById(id)
	if exist == nil {
		return errors.New("Driver not found")
	}
	if err != nil {
		return err
	}
	_, err = r.db.Exec("DELETE FROM drivers WHERE id = $1", id)
	if err != nil {
		return err
	}
	return nil
}

func (r *PostgresDriverRepository) UpdateById(driver *models.Driver) error {
	exist, err := r.GetById(driver.ID)
	if exist == nil {
		return errors.New("Driver not found")
	}
	if err != nil {
		return err
	}
	_, err = r.db.Exec("UPDATE drivers SET name = $1, surname = $2, patronymic = $3, birth_date = $4, passport_series = $5, snils = $6, license_series = $7 WHERE id = $8",
		driver.Name, driver.Surname, driver.Patronymic, driver.BirthDate, driver.PassportSeries, driver.Snils, driver.LicenseSeries, driver.ID)
	if err != nil {
		return err
	}
	return nil
}

package repository

import (
	"backend/pkg/models"
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"reflect"
	"testing"
	"time"
)

func setupMockDriver(t *testing.T) (*sql.DB, sqlmock.Sqlmock, *PostgresDriverRepository) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Ошибка создания mock базы данных: %v", err)
	}
	repo := &PostgresDriverRepository{db: db}
	return db, mock, repo
}

func TestPostgresDriverRepository(t *testing.T) {
	t.Run("NewPostgresDriverRepository", func(t *testing.T) {
		db, _, _ := setupMockDriver(t)
		defer db.Close()

		repo, err := NewPostgresDriverRepository(db)
		if err != nil {
			t.Errorf("Ошибка при создании репозитория: %v", err)
		}
		if repo == nil {
			t.Error("Репозиторий не должен быть nil")
		}
	})

	t.Run("GetById", func(t *testing.T) {
		db, mock, repo := setupMockDriver(t)
		defer db.Close()

		driverID := uuid.New().String()
		driver := &models.Driver{
			ID:             driverID,
			Name:           "Иван",
			Surname:        "Иванов",
			Patronymic:     "Иванович",
			BirthDate:      time.Date(1980, 1, 1, 0, 0, 0, 0, time.UTC),
			PassportSeries: "1234 567890",
			Snils:          "123-456-789 01",
			LicenseSeries:  "AB1234567",
		}

		rows := sqlmock.NewRows([]string{"id", "name", "surname", "patronymic", "birth_date", "passport_series", "snils", "license_series"}).
			AddRow(driver.ID, driver.Name, driver.Surname, driver.Patronymic, driver.BirthDate, driver.PassportSeries, driver.Snils, driver.LicenseSeries)
		mock.ExpectQuery(`SELECT id, name, surname, patronymic, birth_date, passport_series, snils, license_series FROM drivers WHERE id = \$1`).
			WithArgs(driverID).
			WillReturnRows(rows)

		retrievedDriver, err := repo.GetById(driverID)
		if err != nil {
			t.Errorf("Ошибка при получении водителя по ID: %v", err)
		}
		if !reflect.DeepEqual(driver, retrievedDriver) {
			t.Errorf("Полученный водитель не совпадает: ожидался %v, получен %v", driver, retrievedDriver)
		}

		mock.ExpectQuery(`SELECT id, name, surname, patronymic, birth_date, passport_series, snils, license_series FROM drivers WHERE id = \$1`).
			WithArgs("nonexistent").
			WillReturnError(sql.ErrNoRows)

		_, err = repo.GetById("nonexistent")
		if err == nil || err.Error() != "Driver not found" {
			t.Errorf("Ожидалась ошибка 'Driver not found', получена: %v", err)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Не все ожидаемые SQL запросы были выполнены: %v", err)
		}
	})

	t.Run("GetByPassportSeries", func(t *testing.T) {
		db, mock, repo := setupMockDriver(t)
		defer db.Close()

		driver := &models.Driver{
			ID:             uuid.New().String(),
			Name:           "Петр",
			Surname:        "Петров",
			Patronymic:     "Петрович",
			BirthDate:      time.Date(1985, 2, 2, 0, 0, 0, 0, time.UTC),
			PassportSeries: "9876 543210",
			Snils:          "987-654-321 02",
			LicenseSeries:  "CD9876543",
		}

		rows := sqlmock.NewRows([]string{"id", "name", "surname", "patronymic", "birth_date", "passport_series", "snils", "license_series"}).
			AddRow(driver.ID, driver.Name, driver.Surname, driver.Patronymic, driver.BirthDate, driver.PassportSeries, driver.Snils, driver.LicenseSeries)
		mock.ExpectQuery(`SELECT id, name, surname, patronymic, birth_date, passport_series, snils, license_series FROM drivers WHERE passport_series = \$1`).
			WithArgs("9876 543210").
			WillReturnRows(rows)

		retrievedDriver, err := repo.GetByPassportSeries("9876 543210")
		if err != nil {
			t.Errorf("Ошибка при получении водителя по серии паспорта: %v", err)
		}
		if !reflect.DeepEqual(driver, retrievedDriver) {
			t.Errorf("Полученный водитель не совпадает: ожидался %v, получен %v", driver, retrievedDriver)
		}

		mock.ExpectQuery(`SELECT id, name, surname, patronymic, birth_date, passport_series, snils, license_series FROM drivers WHERE passport_series = \$1`).
			WithArgs("nonexistent").
			WillReturnError(sql.ErrNoRows)

		_, err = repo.GetByPassportSeries("nonexistent")
		if err == nil || err.Error() != "Driver not found" {
			t.Errorf("Ожидалась ошибка 'Driver not found', получена: %v", err)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Не все ожидаемые SQL запросы были выполнены: %v", err)
		}
	})

	t.Run("Add", func(t *testing.T) {
		db, mock, repo := setupMockDriver(t)
		defer db.Close()

		driver := &models.Driver{
			ID:             uuid.New().String(),
			Name:           "Алексей",
			Surname:        "Сидоров",
			Patronymic:     "Алексеевич",
			BirthDate:      time.Date(1990, 3, 3, 0, 0, 0, 0, time.UTC),
			PassportSeries: "1111 222222",
			Snils:          "111-222-333 03",
			LicenseSeries:  "EF1112223",
		}

		mock.ExpectQuery(`SELECT id, name, surname, patronymic, birth_date, passport_series, snils, license_series FROM drivers WHERE passport_series = \$1`).
			WithArgs(driver.PassportSeries).
			WillReturnError(sql.ErrNoRows)

		mock.ExpectExec(`INSERT into drivers \(id, name, surname, patronymic, birth_date, passport_series, snils, license_series \) VALUES \(\$1, \$2, \$3, \$4, \$5, \$6, \$7, \$8\)`).
			WithArgs(driver.ID, driver.Name, driver.Surname, driver.Patronymic, driver.BirthDate, driver.PassportSeries, driver.Snils, driver.LicenseSeries).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.Add(driver)
		if err != nil {
			t.Errorf("Ошибка при добавлении водителя: %v", err)
		}

		mock.ExpectQuery(`SELECT id, name, surname, patronymic, birth_date, passport_series, snils, license_series FROM drivers WHERE passport_series = \$1`).
			WithArgs(driver.PassportSeries).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "surname", "patronymic", "birth_date", "passport_series", "snils", "license_series"}).
				AddRow(driver.ID, driver.Name, driver.Surname, driver.Patronymic, driver.BirthDate, driver.PassportSeries, driver.Snils, driver.LicenseSeries))

		err = repo.Add(driver)
		if err == nil || err.Error() != "Driver already exists" {
			t.Errorf("Ожидалась ошибка 'Driver already exists', получена: %v", err)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Не все ожидаемые SQL запросы были выполнены: %v", err)
		}
	})

	t.Run("GetAll", func(t *testing.T) {
		db, mock, repo := setupMockDriver(t)
		defer db.Close()

		driver1 := models.Driver{
			ID:             uuid.New().String(),
			Name:           "Сергей",
			Surname:        "Козлов",
			Patronymic:     "Викторович",
			BirthDate:      time.Date(1975, 4, 4, 0, 0, 0, 0, time.UTC),
			PassportSeries: "3333 444444",
			Snils:          "333-444-555 04",
			LicenseSeries:  "GH3334445",
		}
		driver2 := models.Driver{
			ID:             uuid.New().String(),
			Name:           "Михаил",
			Surname:        "Смирнов",
			Patronymic:     "Александрович",
			BirthDate:      time.Date(1988, 5, 5, 0, 0, 0, 0, time.UTC),
			PassportSeries: "5555 666666",
			Snils:          "555-666-777 05",
			LicenseSeries:  "IJ5556667",
		}

		rows := sqlmock.NewRows([]string{"id", "name", "surname", "patronymic", "birth_date", "passport_series", "snils", "license_series"}).
			AddRow(driver1.ID, driver1.Name, driver1.Surname, driver1.Patronymic, driver1.BirthDate, driver1.PassportSeries, driver1.Snils, driver1.LicenseSeries).
			AddRow(driver2.ID, driver2.Name, driver2.Surname, driver2.Patronymic, driver2.BirthDate, driver2.PassportSeries, driver2.Snils, driver2.LicenseSeries)
		mock.ExpectQuery(`SELECT id, name, surname, patronymic, birth_date, passport_series, snils, license_series FROM drivers`).
			WillReturnRows(rows)

		drivers, err := repo.GetAll()
		if err != nil {
			t.Errorf("Ошибка при получении всех водителей: %v", err)
		}
		if len(drivers) != 2 {
			t.Errorf("Ожидалось 2 водителя, получено: %d", len(drivers))
		}
		foundDriver1, foundDriver2 := false, false
		for _, d := range drivers {
			if reflect.DeepEqual(d, driver1) {
				foundDriver1 = true
			}
			if reflect.DeepEqual(d, driver2) {
				foundDriver2 = true
			}
		}
		if !foundDriver1 || !foundDriver2 {
			t.Errorf("Не все водители найдены в списке: driver1=%v, driver2=%v", foundDriver1, foundDriver2)
		}

		mock.ExpectQuery(`SELECT id, name, surname, patronymic, birth_date, passport_series, snils, license_series FROM drivers`).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "surname", "patronymic", "birth_date", "passport_series", "snils", "license_series"}))

		drivers, err = repo.GetAll()
		if err != nil {
			t.Errorf("Ошибка при получении пустого списка водителей: %v", err)
		}
		if len(drivers) != 0 {
			t.Errorf("Ожидался пустой список, получено: %d элементов", len(drivers))
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Не все ожидаемые SQL запросы были выполнены: %v", err)
		}
	})

	t.Run("DeleteById", func(t *testing.T) {
		db, mock, repo := setupMockDriver(t)
		defer db.Close()

		driverID := uuid.New().String()

		mock.ExpectQuery(`SELECT id, name, surname, patronymic, birth_date, passport_series, snils, license_series FROM drivers WHERE id = \$1`).
			WithArgs(driverID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "surname", "patronymic", "birth_date", "passport_series", "snils", "license_series"}).
				AddRow(driverID, "Иван", "Иванов", "Иванович", time.Date(1980, 1, 1, 0, 0, 0, 0, time.UTC), "1234 567890", "123-456-789 01", "AB1234567"))

		mock.ExpectExec(`DELETE FROM drivers WHERE id = \$1`).
			WithArgs(driverID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.DeleteById(driverID)
		if err != nil {
			t.Errorf("Ошибка при удалении водителя: %v", err)
		}

		mock.ExpectQuery(`SELECT id, name, surname, patronymic, birth_date, passport_series, snils, license_series FROM drivers WHERE id = \$1`).
			WithArgs("nonexistent").
			WillReturnError(sql.ErrNoRows)

		err = repo.DeleteById("nonexistent")
		if err == nil || err.Error() != "Driver not found" {
			t.Errorf("Ожидалась ошибка 'Driver not found', получена: %v", err)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Не все ожидаемые SQL запросы были выполнены: %v", err)
		}
	})

	t.Run("UpdateById", func(t *testing.T) {
		db, mock, repo := setupMockDriver(t)
		defer db.Close()

		driver := &models.Driver{
			ID:             uuid.New().String(),
			Name:           "Алексей",
			Surname:        "Сидоров",
			Patronymic:     "Алексеевич",
			BirthDate:      time.Date(1990, 3, 3, 0, 0, 0, 0, time.UTC),
			PassportSeries: "1111 222222",
			Snils:          "111-222-333 03",
			LicenseSeries:  "EF1112223",
		}

		mock.ExpectQuery(`SELECT id, name, surname, patronymic, birth_date, passport_series, snils, license_series FROM drivers WHERE id = \$1`).
			WithArgs(driver.ID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "surname", "patronymic", "birth_date", "passport_series", "snils", "license_series"}).
				AddRow(driver.ID, driver.Name, driver.Surname, driver.Patronymic, driver.BirthDate, driver.PassportSeries, driver.Snils, driver.LicenseSeries))

		mock.ExpectExec(`UPDATE drivers SET name = \$1, surname = \$2, patronymic = \$3, birth_date = \$4, passport_series = \$5, snils = \$6, license_series = \$7 WHERE id = \$8`).
			WithArgs(driver.Name, driver.Surname, driver.Patronymic, driver.BirthDate, driver.PassportSeries, driver.Snils, driver.LicenseSeries, driver.ID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.UpdateById(driver)
		if err != nil {
			t.Errorf("Ошибка при обновлении водителя: %v", err)
		}

		mock.ExpectQuery(`SELECT id, name, surname, patronymic, birth_date, passport_series, snils, license_series FROM drivers WHERE id = \$1`).
			WithArgs("nonexistent").
			WillReturnError(sql.ErrNoRows)

		err = repo.UpdateById(&models.Driver{ID: "nonexistent"})
		if err == nil || err.Error() != "Driver not found" {
			t.Errorf("Ожидалась ошибка 'Driver not found', получена: %v", err)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Не все ожидаемые SQL запросы были выполнены: %v", err)
		}
	})
}

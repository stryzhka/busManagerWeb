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

func setupMockBus(t *testing.T) (*sql.DB, sqlmock.Sqlmock, *PostgresBusRepository) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Ошибка создания mock базы данных: %v", err)
	}
	repo := &PostgresBusRepository{db: db}
	return db, mock, repo
}

func TestPostgresBusRepository(t *testing.T) {
	t.Run("NewPostgresBusRepository", func(t *testing.T) {
		db, _, _ := setupMockBus(t)
		defer db.Close()

		repo, err := NewPostgresBusRepository(db)
		if err != nil {
			t.Errorf("Ошибка при создании репозитория: %v", err)
		}
		if repo == nil {
			t.Error("Репозиторий не должен быть nil")
		}
	})

	t.Run("GetById", func(t *testing.T) {
		db, mock, repo := setupMockBus(t)
		defer db.Close()

		busID := uuid.New().String()
		bus := &models.Bus{
			ID:             busID,
			Brand:          "Volvo",
			BusModel:       "B9R",
			RegisterNumber: "ABC123",
			AssemblyDate:   time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			LastRepairDate: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
		}

		rows := sqlmock.NewRows([]string{"id", "brand", "bus_model", "register_number", "assembly_date", "last_repair_date"}).
			AddRow(bus.ID, bus.Brand, bus.BusModel, bus.RegisterNumber, bus.AssemblyDate, bus.LastRepairDate)
		mock.ExpectQuery(`SELECT id, brand, bus_model, register_number, assembly_date, last_repair_date FROM buses WHERE id = \$1`).
			WithArgs(busID).
			WillReturnRows(rows)

		retrievedBus, err := repo.GetById(busID)
		if err != nil {
			t.Errorf("Ошибка при получении автобуса по ID: %v", err)
		}
		if !reflect.DeepEqual(bus, retrievedBus) {
			t.Errorf("Полученный автобус не совпадает: ожидался %v, получен %v", bus, retrievedBus)
		}

		mock.ExpectQuery(`SELECT id, brand, bus_model, register_number, assembly_date, last_repair_date FROM buses WHERE id = \$1`).
			WithArgs("nonexistent").
			WillReturnError(sql.ErrNoRows)

		_, err = repo.GetById("nonexistent")
		if err == nil || err.Error() != "Bus not found" {
			t.Errorf("Ожидалась ошибка 'Bus not found', получена: %v", err)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Не все ожидаемые SQL запросы были выполнены: %v", err)
		}
	})

	t.Run("GetByNumber", func(t *testing.T) {
		db, mock, repo := setupMockBus(t)
		defer db.Close()

		bus := &models.Bus{
			ID:             uuid.New().String(),
			Brand:          "Mercedes",
			BusModel:       "Travego",
			RegisterNumber: "DEF456",
			AssemblyDate:   time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
		}

		rows := sqlmock.NewRows([]string{"id", "brand", "bus_model", "register_number", "assembly_date", "last_repair_date"}).
			AddRow(bus.ID, bus.Brand, bus.BusModel, bus.RegisterNumber, bus.AssemblyDate, bus.LastRepairDate)
		mock.ExpectQuery(`SELECT id, brand, bus_model, register_number, assembly_date, last_repair_date FROM buses WHERE register_number = \$1`).
			WithArgs("DEF456").
			WillReturnRows(rows)

		retrievedBus, err := repo.GetByNumber("DEF456")
		if err != nil {
			t.Errorf("Ошибка при получении автобуса по номеру: %v", err)
		}
		if !reflect.DeepEqual(bus, retrievedBus) {
			t.Errorf("Полученный автобус не совпадает: ожидался %v, получен %v", bus, retrievedBus)
		}

		mock.ExpectQuery(`SELECT id, brand, bus_model, register_number, assembly_date, last_repair_date FROM buses WHERE register_number = \$1`).
			WithArgs("NONEXISTENT").
			WillReturnError(sql.ErrNoRows)

		_, err = repo.GetByNumber("NONEXISTENT")
		if err == nil || err.Error() != "Bus not found" {
			t.Errorf("Ожидалась ошибка 'Bus not found', получена: %v", err)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Не все ожидаемые SQL запросы были выполнены: %v", err)
		}
	})

	t.Run("Add", func(t *testing.T) {
		db, mock, repo := setupMockBus(t)
		defer db.Close()

		bus := &models.Bus{
			ID:             uuid.New().String(),
			Brand:          "Volvo",
			BusModel:       "B9R",
			RegisterNumber: "GHI789",
			AssemblyDate:   time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		}

		mock.ExpectQuery(`SELECT id, brand, bus_model, register_number, assembly_date, last_repair_date FROM buses WHERE register_number = \$1`).
			WithArgs(bus.RegisterNumber).
			WillReturnError(sql.ErrNoRows)

		// Точный SQL-запрос из кода репозитория
		mock.ExpectExec(`INSERT into buses \(id, brand, bus_model, register_number, assembly_date, last_repair_date \) VALUES \(\$1, \$2, \$3, \$4, \$5, \$6\)`).
			WithArgs(bus.ID, bus.Brand, bus.BusModel, bus.RegisterNumber, bus.AssemblyDate, bus.LastRepairDate).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.Add(bus)
		if err != nil {
			t.Errorf("Ошибка при добавлении автобуса: %v", err)
		}

		mock.ExpectQuery(`SELECT id, brand, bus_model, register_number, assembly_date, last_repair_date FROM buses WHERE register_number = \$1`).
			WithArgs(bus.RegisterNumber).
			WillReturnRows(sqlmock.NewRows([]string{"id", "brand", "bus_model", "register_number", "assembly_date", "last_repair_date"}).
				AddRow(bus.ID, bus.Brand, bus.BusModel, bus.RegisterNumber, bus.AssemblyDate, bus.LastRepairDate))

		err = repo.Add(bus)
		if err == nil || err.Error() != "Bus already exists" {
			t.Errorf("Ожидалась ошибка 'Bus already exists', получена: %v", err)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Не все ожидаемые SQL запросы были выполнены: %v", err)
		}
	})

	t.Run("GetAll", func(t *testing.T) {
		db, mock, repo := setupMockBus(t)
		defer db.Close()

		bus1 := models.Bus{
			ID:             uuid.New().String(),
			Brand:          "Volvo",
			BusModel:       "B9R",
			RegisterNumber: "JKL012",
			AssemblyDate:   time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		}
		bus2 := models.Bus{
			ID:             uuid.New().String(),
			Brand:          "Mercedes",
			BusModel:       "Travego",
			RegisterNumber: "MNO345",
			AssemblyDate:   time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
		}

		rows := sqlmock.NewRows([]string{"id", "brand", "bus_model", "register_number", "assembly_date", "last_repair_date"}).
			AddRow(bus1.ID, bus1.Brand, bus1.BusModel, bus1.RegisterNumber, bus1.AssemblyDate, bus1.LastRepairDate).
			AddRow(bus2.ID, bus2.Brand, bus2.BusModel, bus2.RegisterNumber, bus2.AssemblyDate, bus2.LastRepairDate)
		mock.ExpectQuery(`SELECT id, brand, bus_model, register_number, assembly_date, last_repair_date FROM buses`).
			WillReturnRows(rows)

		buses, err := repo.GetAll()
		if err != nil {
			t.Errorf("Ошибка при получении всех автобусов: %v", err)
		}
		if len(buses) != 2 {
			t.Errorf("Ожидалось 2 автобуса, получено: %d", len(buses))
		}
		foundBus1, foundBus2 := false, false
		for _, b := range buses {
			if reflect.DeepEqual(b, bus1) {
				foundBus1 = true
			}
			if reflect.DeepEqual(b, bus2) {
				foundBus2 = true
			}
		}
		if !foundBus1 || !foundBus2 {
			t.Errorf("Не все автобусы найдены в списке: bus1=%v, bus2=%v", foundBus1, foundBus2)
		}

		mock.ExpectQuery(`SELECT id, brand, bus_model, register_number, assembly_date, last_repair_date FROM buses`).
			WillReturnRows(sqlmock.NewRows([]string{"id", "brand", "bus_model", "register_number", "assembly_date", "last_repair_date"}))

		buses, err = repo.GetAll()
		if err != nil {
			t.Errorf("Ошибка при получении пустого списка автобусов: %v", err)
		}
		if len(buses) != 0 {
			t.Errorf("Ожидался пустой список, получено: %d элементов", len(buses))
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Не все ожидаемые SQL запросы были выполнены: %v", err)
		}
	})

	t.Run("DeleteById", func(t *testing.T) {
		db, mock, repo := setupMockBus(t)
		defer db.Close()

		busID := uuid.New().String()

		mock.ExpectQuery(`SELECT id, brand, bus_model, register_number, assembly_date, last_repair_date FROM buses WHERE id = \$1`).
			WithArgs(busID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "brand", "bus_model", "register_number", "assembly_date", "last_repair_date"}).
				AddRow(busID, "Volvo", "B9R", "PQR678", time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC), time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)))

		mock.ExpectExec(`DELETE FROM buses WHERE id = \$1`).
			WithArgs(busID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.DeleteById(busID)
		if err != nil {
			t.Errorf("Ошибка при удалении автобуса: %v", err)
		}

		mock.ExpectQuery(`SELECT id, brand, bus_model, register_number, assembly_date, last_repair_date FROM buses WHERE id = \$1`).
			WithArgs("nonexistent").
			WillReturnError(sql.ErrNoRows)

		err = repo.DeleteById("nonexistent")
		if err == nil || err.Error() != "Bus not found" {
			t.Errorf("Ожидалась ошибка 'Bus not found', получена: %v", err)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Не все ожидаемые SQL запросы были выполнены: %v", err)
		}
	})

	t.Run("UpdateById", func(t *testing.T) {
		db, mock, repo := setupMockBus(t)
		defer db.Close()

		bus := &models.Bus{
			ID:             uuid.New().String(),
			Brand:          "Volvo",
			BusModel:       "B9R",
			RegisterNumber: "STU901",
			AssemblyDate:   time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			LastRepairDate: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
		}

		mock.ExpectQuery(`SELECT id, brand, bus_model, register_number, assembly_date, last_repair_date FROM buses WHERE id = \$1`).
			WithArgs(bus.ID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "brand", "bus_model", "register_number", "assembly_date", "last_repair_date"}).
				AddRow(bus.ID, bus.Brand, bus.BusModel, bus.RegisterNumber, bus.AssemblyDate, bus.LastRepairDate))

		mock.ExpectExec(`UPDATE buses SET brand = \$1, bus_model = \$2, register_number = \$3, assembly_date = \$4, last_repair_date = \$5 WHERE id = \$6`).
			WithArgs(bus.Brand, bus.BusModel, bus.RegisterNumber, bus.AssemblyDate, bus.LastRepairDate, bus.ID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.UpdateById(bus)
		if err != nil {
			t.Errorf("Ошибка при обновлении автобуса: %v", err)
		}

		mock.ExpectQuery(`SELECT id, brand, bus_model, register_number, assembly_date, last_repair_date FROM buses WHERE id = \$1`).
			WithArgs("nonexistent").
			WillReturnError(sql.ErrNoRows)

		err = repo.UpdateById(&models.Bus{ID: "nonexistent"})
		if err == nil || err.Error() != "Bus not found" {
			t.Errorf("Ожидалась ошибка 'Bus not found', получена: %v", err)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Не все ожидаемые SQL запросы были выполнены: %v", err)
		}
	})
}

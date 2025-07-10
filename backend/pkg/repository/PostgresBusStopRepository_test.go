package repository

import (
	"backend/pkg/models"
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"reflect"
	"testing"
)

func setupMock(t *testing.T) (*sql.DB, sqlmock.Sqlmock, *PostgresBusStopRepository) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Ошибка создания mock базы данных: %v", err)
	}
	repo := &PostgresBusStopRepository{db: db}
	return db, mock, repo
}

func TestPostgresBusStopRepository(t *testing.T) {
	t.Run("NewPostgresBusStopRepository", func(t *testing.T) {
		db, _, _ := setupMock(t)
		defer db.Close()

		repo, err := NewPostgresBusStopRepository(db)
		if err != nil {
			t.Errorf("Ошибка при создании репозитория: %v", err)
		}
		if repo == nil {
			t.Error("Репозиторий не должен быть nil")
		}
	})

	t.Run("GetById", func(t *testing.T) {
		db, mock, repo := setupMock(t)
		defer db.Close()

		stopID := uuid.New().String()
		stop := &models.BusStop{
			ID:   stopID,
			Lat:  55.7558,
			Long: 37.6173,
			Name: "Центральная",
		}

		rows := sqlmock.NewRows([]string{"id", "lat", "long", "name"}).
			AddRow(stop.ID, stop.Lat, stop.Long, stop.Name)
		mock.ExpectQuery(`SELECT id, lat, long, name FROM bus_stops WHERE id = \$1`).
			WithArgs(stopID).
			WillReturnRows(rows)

		retrievedStop, err := repo.GetById(stopID)
		if err != nil {
			t.Errorf("Ошибка при получении остановки по ID: %v", err)
		}
		if !reflect.DeepEqual(stop, retrievedStop) {
			t.Errorf("Полученная остановка не совпадает: ожидалась %v, получена %v", stop, retrievedStop)
		}

		mock.ExpectQuery(`SELECT id, lat, long, name FROM bus_stops WHERE id = \$1`).
			WithArgs("nonexistent").
			WillReturnError(sql.ErrNoRows)

		_, err = repo.GetById("nonexistent")
		if err == nil || err.Error() != "Bus stop not found" {
			t.Errorf("Ожидалась ошибка 'Bus stop not found', получена: %v", err)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Не все ожидаемые SQL запросы были выполнены: %v", err)
		}
	})

	t.Run("GetByName", func(t *testing.T) {
		db, mock, repo := setupMock(t)
		defer db.Close()

		stop := &models.BusStop{
			ID:   uuid.New().String(),
			Lat:  59.9343,
			Long: 30.3351,
			Name: "Невский проспект",
		}

		rows := sqlmock.NewRows([]string{"id", "lat", "long", "name"}).
			AddRow(stop.ID, stop.Lat, stop.Long, stop.Name)
		mock.ExpectQuery(`SELECT id, lat, long, name FROM bus_stops WHERE name = \$4`).
			WithArgs("Невский проспект").
			WillReturnRows(rows)

		retrievedStop, err := repo.GetByName("Невский проспект")
		if err != nil {
			t.Errorf("Ошибка при получении остановки по имени: %v", err)
		}
		if !reflect.DeepEqual(stop, retrievedStop) {
			t.Errorf("Полученная остановка не совпадает: ожидалась %v, получена %v", stop, retrievedStop)
		}

		mock.ExpectQuery(`SELECT id, lat, long, name FROM bus_stops WHERE name = \$4`).
			WithArgs("nonexistent").
			WillReturnError(sql.ErrNoRows)

		_, err = repo.GetByName("nonexistent")
		if err == nil || err.Error() != "Bus stop not found" {
			t.Errorf("Ожидалась ошибка 'Bus stop not found', получена: %v", err)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Не все ожидаемые SQL запросы были выполнены: %v", err)
		}
	})

	t.Run("Add", func(t *testing.T) {
		db, mock, repo := setupMock(t)
		defer db.Close()

		stop := &models.BusStop{
			ID:   uuid.New().String(),
			Lat:  48.1351,
			Long: 11.5820,
			Name: "Мариенплац",
		}

		// Проверка успешного добавления
		mock.ExpectQuery(`SELECT id, lat, long, name FROM bus_stops WHERE name = \$4`).
			WithArgs(stop.Name).
			WillReturnError(sql.ErrNoRows)

		mock.ExpectExec(`INSERT into bus_stops \(id, lat, long, name \) VALUES \(\$1, \$2, \$3, \$4\)`).
			WithArgs(stop.ID, stop.Lat, stop.Long, stop.Name).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.Add(stop)
		if err != nil {
			t.Errorf("Ошибка при добавлении остановки: %v", err)
		}

		// Проверка добавления с дублирующимся именем
		mock.ExpectQuery(`SELECT id, lat, long, name FROM bus_stops WHERE name = \$4`).
			WithArgs(stop.Name).
			WillReturnRows(sqlmock.NewRows([]string{"id", "lat", "long", "name"}).
				AddRow(stop.ID, stop.Lat, stop.Long, stop.Name))

		err = repo.Add(stop)
		if err == nil || err.Error() != "Bus stop already exists" {
			t.Errorf("Ожидалась ошибка 'Bus stop already exists', получена: %v", err)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Не все ожидаемые SQL запросы были выполнены: %v", err)
		}
	})

	t.Run("GetAll", func(t *testing.T) {
		db, mock, repo := setupMock(t)
		defer db.Close()

		stop1 := models.BusStop{
			ID:   uuid.New().String(),
			Lat:  51.5074,
			Long: -0.1278,
			Name: "Трафальгарская площадь",
		}
		stop2 := models.BusStop{
			ID:   uuid.New().String(),
			Lat:  48.8566,
			Long: 2.3522,
			Name: "Эйфелева башня",
		}

		rows := sqlmock.NewRows([]string{"id", "lat", "long", "name"}).
			AddRow(stop1.ID, stop1.Lat, stop1.Long, stop1.Name).
			AddRow(stop2.ID, stop2.Lat, stop2.Long, stop2.Name)
		mock.ExpectQuery(`SELECT id, lat, long, name FROM bus_stops`).
			WillReturnRows(rows)

		stops, err := repo.GetAll()
		if err != nil {
			t.Errorf("Ошибка при получении всех остановок: %v", err)
		}
		if len(stops) != 2 {
			t.Errorf("Ожидалось 2 остановки, получено: %d", len(stops))
		}
		foundStop1, foundStop2 := false, false
		for _, s := range stops {
			if reflect.DeepEqual(s, stop1) {
				foundStop1 = true
			}
			if reflect.DeepEqual(s, stop2) {
				foundStop2 = true
			}
		}
		if !foundStop1 || !foundStop2 {
			t.Errorf("Не все остановки найдены в списке: stop1=%v, stop2=%v", foundStop1, foundStop2)
		}

		mock.ExpectQuery(`SELECT id, lat, long, name FROM bus_stops`).
			WillReturnRows(sqlmock.NewRows([]string{"id", "lat", "long", "name"}))

		stops, err = repo.GetAll()
		if err != nil {
			t.Errorf("Ошибка при получении пустого списка остановок: %v", err)
		}
		if len(stops) != 0 {
			t.Errorf("Ожидался пустой список, получено: %d элементов", len(stops))
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Не все ожидаемые SQL запросы были выполнены: %v", err)
		}
	})

	t.Run("DeleteById", func(t *testing.T) {
		db, mock, repo := setupMock(t)
		defer db.Close()

		stopID := uuid.New().String()

		mock.ExpectQuery(`SELECT id, lat, long, name FROM bus_stops WHERE id = \$1`).
			WithArgs(stopID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "lat", "long", "name"}).
				AddRow(stopID, 55.7558, 37.6173, "Центральная"))

		mock.ExpectExec(`DELETE FROM bus_stops WHERE id = \$1`).
			WithArgs(stopID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.DeleteById(stopID)
		if err != nil {
			t.Errorf("Ошибка при удалении остановки: %v", err)
		}

		mock.ExpectQuery(`SELECT id, lat, long, name FROM bus_stops WHERE id = \$1`).
			WithArgs("nonexistent").
			WillReturnError(sql.ErrNoRows)

		err = repo.DeleteById("nonexistent")
		if err == nil || err.Error() != "Bus stop not found" {
			t.Errorf("Ожидалась ошибка 'Bus stop not found', получена: %v", err)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Не все ожидаемые SQL запросы были выполнены: %v", err)
		}
	})

	t.Run("UpdateById", func(t *testing.T) {
		db, mock, repo := setupMock(t)
		defer db.Close()

		stop := &models.BusStop{
			ID:   uuid.New().String(),
			Lat:  48.1351,
			Long: 11.5820,
			Name: "Мариенплац",
		}

		mock.ExpectQuery(`SELECT id, lat, long, name FROM bus_stops WHERE id = \$1`).
			WithArgs(stop.ID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "lat", "long", "name"}).
				AddRow(stop.ID, stop.Lat, stop.Long, stop.Name))

		mock.ExpectExec(`UPDATE bus_stops SET lat = \$1, long = \$2, name = \$3 WHERE id = \$4`).
			WithArgs(stop.Lat, stop.Long, stop.Name, stop.ID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.UpdateById(stop)
		if err != nil {
			t.Errorf("Ошибка при обновлении остановки: %v", err)
		}

		mock.ExpectQuery(`SELECT id, lat, long, name FROM bus_stops WHERE id = \$1`).
			WithArgs("nonexistent").
			WillReturnError(sql.ErrNoRows)

		err = repo.UpdateById(&models.BusStop{ID: "nonexistent"})
		if err == nil || err.Error() != "Bus stop not found" {
			t.Errorf("Ожидалась ошибка 'Bus stop not found', получена: %v", err)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Не все ожидаемые SQL запросы были выполнены: %v", err)
		}
	})
}

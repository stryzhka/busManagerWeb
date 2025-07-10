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

func setupMockRoute(t *testing.T) (*sql.DB, sqlmock.Sqlmock, *PostgresRouteRepository) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Ошибка создания mock базы данных: %v", err)
	}
	repo := &PostgresRouteRepository{db: db}
	return db, mock, repo
}

func TestPostgresRouteRepository(t *testing.T) {
	t.Run("NewPostgresRouteRepository", func(t *testing.T) {
		db, _, _ := setupMockRoute(t)
		defer db.Close()

		repo, err := NewPostgresRouteRepository(db)
		if err != nil {
			t.Errorf("Ошибка при создании репозитория: %v", err)
		}
		if repo == nil {
			t.Error("Репозиторий не должен быть nil")
		}
	})

	t.Run("GetById", func(t *testing.T) {
		db, mock, repo := setupMockRoute(t)
		defer db.Close()

		routeID := uuid.New().String()
		route := &models.Route{
			ID:     routeID,
			Number: "101",
		}

		rows := sqlmock.NewRows([]string{"id", "number"}).
			AddRow(route.ID, route.Number)
		mock.ExpectQuery(`SELECT id, number FROM routes WHERE id = \$1`).
			WithArgs(routeID).
			WillReturnRows(rows)

		retrievedRoute, err := repo.GetById(routeID)
		if err != nil {
			t.Errorf("Ошибка при получении маршрута по ID: %v", err)
		}
		if !reflect.DeepEqual(route, retrievedRoute) {
			t.Errorf("Полученный маршрут не совпадает: ожидался %v, получен %v", route, retrievedRoute)
		}

		mock.ExpectQuery(`SELECT id, number FROM routes WHERE id = \$1`).
			WithArgs("nonexistent").
			WillReturnError(sql.ErrNoRows)

		_, err = repo.GetById("nonexistent")
		if err == nil || err.Error() != "Route not found" {
			t.Errorf("Ожидалась ошибка 'Route not found', получена: %v", err)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Не все ожидаемые SQL запросы были выполнены: %v", err)
		}
	})

	t.Run("GetByNumber", func(t *testing.T) {
		db, mock, repo := setupMockRoute(t)
		defer db.Close()

		route := &models.Route{
			ID:     uuid.New().String(),
			Number: "102",
		}

		rows := sqlmock.NewRows([]string{"id", "number"}).
			AddRow(route.ID, route.Number)
		mock.ExpectQuery(`SELECT id, number FROM routes WHERE number = \$1`).
			WithArgs("102").
			WillReturnRows(rows)

		retrievedRoute, err := repo.GetByNumber("102")
		if err != nil {
			t.Errorf("Ошибка при получении маршрута по номеру: %v", err)
		}
		if !reflect.DeepEqual(route, retrievedRoute) {
			t.Errorf("Полученный маршрут не совпадает: ожидался %v, получен %v", route, retrievedRoute)
		}

		mock.ExpectQuery(`SELECT id, number FROM routes WHERE number = \$1`).
			WithArgs("nonexistent").
			WillReturnError(sql.ErrNoRows)

		_, err = repo.GetByNumber("nonexistent")
		if err == nil || err.Error() != "Route not found" {
			t.Errorf("Ожидалась ошибка 'Route not found', получена: %v", err)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Не все ожидаемые SQL запросы были выполнены: %v", err)
		}
	})

	t.Run("Add", func(t *testing.T) {
		db, mock, repo := setupMockRoute(t)
		defer db.Close()

		route := &models.Route{
			ID:     uuid.New().String(),
			Number: "103",
		}

		mock.ExpectQuery(`SELECT id, number FROM routes WHERE number = \$1`).
			WithArgs(route.Number).
			WillReturnError(sql.ErrNoRows)

		mock.ExpectExec(`INSERT into routes \(id, number\) VALUES \(\$1, \$2\)`).
			WithArgs(route.ID, route.Number).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.Add(route)
		if err != nil {
			t.Errorf("Ошибка при добавлении маршрута: %v", err)
		}

		mock.ExpectQuery(`SELECT id, number FROM routes WHERE number = \$1`).
			WithArgs(route.Number).
			WillReturnRows(sqlmock.NewRows([]string{"id", "number"}).
				AddRow(route.ID, route.Number))

		err = repo.Add(route)
		if err == nil || err.Error() != "Route already exists" {
			t.Errorf("Ожидалась ошибка 'Route already exists', получена: %v", err)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Не все ожидаемые SQL запросы были выполнены: %v", err)
		}
	})

	t.Run("GetAll", func(t *testing.T) {
		db, mock, repo := setupMockRoute(t)
		defer db.Close()

		route1 := models.Route{
			ID:     uuid.New().String(),
			Number: "104",
		}
		route2 := models.Route{
			ID:     uuid.New().String(),
			Number: "105",
		}

		rows := sqlmock.NewRows([]string{"id", "number"}).
			AddRow(route1.ID, route1.Number).
			AddRow(route2.ID, route2.Number)
		mock.ExpectQuery(`SELECT id, number FROM routes`).
			WillReturnRows(rows)

		routes, err := repo.GetAll()
		if err != nil {
			t.Errorf("Ошибка при получении всех маршрутов: %v", err)
		}
		if len(routes) != 2 {
			t.Errorf("Ожидалось 2 маршрута, получено: %d", len(routes))
		}
		foundRoute1, foundRoute2 := false, false
		for _, r := range routes {
			if reflect.DeepEqual(r, route1) {
				foundRoute1 = true
			}
			if reflect.DeepEqual(r, route2) {
				foundRoute2 = true
			}
		}
		if !foundRoute1 || !foundRoute2 {
			t.Errorf("Не все маршруты найдены в списке: route1=%v, route2=%v", foundRoute1, foundRoute2)
		}

		mock.ExpectQuery(`SELECT id, number FROM routes`).
			WillReturnRows(sqlmock.NewRows([]string{"id", "number"}))

		routes, err = repo.GetAll()
		if err != nil {
			t.Errorf("Ошибка при получении пустого списка маршрутов: %v", err)
		}
		if len(routes) != 0 {
			t.Errorf("Ожидался пустой список, получено: %d элементов", len(routes))
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Не все ожидаемые SQL запросы были выполнены: %v", err)
		}
	})

	t.Run("DeleteById", func(t *testing.T) {
		db, mock, repo := setupMockRoute(t)
		defer db.Close()

		routeID := uuid.New().String()

		mock.ExpectQuery(`SELECT id, number FROM routes WHERE id = \$1`).
			WithArgs(routeID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "number"}).
				AddRow(routeID, "106"))

		mock.ExpectExec(`DELETE FROM routes WHERE id = \$1`).
			WithArgs(routeID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.DeleteById(routeID)
		if err != nil {
			t.Errorf("Ошибка при удалении маршрута: %v", err)
		}

		mock.ExpectQuery(`SELECT id, number FROM routes WHERE id = \$1`).
			WithArgs("nonexistent").
			WillReturnError(sql.ErrNoRows)

		err = repo.DeleteById("nonexistent")
		if err == nil || err.Error() != "Route not found" {
			t.Errorf("Ожидалась ошибка 'Route not found', получена: %v", err)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Не все ожидаемые SQL запросы были выполнены: %v", err)
		}
	})

	t.Run("UpdateById", func(t *testing.T) {
		db, mock, repo := setupMockRoute(t)
		defer db.Close()

		route := &models.Route{
			ID:     uuid.New().String(),
			Number: "107",
		}

		mock.ExpectQuery(`SELECT id, number FROM routes WHERE id = \$1`).
			WithArgs(route.ID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "number"}).
				AddRow(route.ID, route.Number))

		mock.ExpectExec(`UPDATE routes SET number = \$1 WHERE id = \$2`).
			WithArgs(route.Number, route.ID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.UpdateById(route)
		if err != nil {
			t.Errorf("Ошибка при обновлении маршрута: %v", err)
		}

		mock.ExpectQuery(`SELECT id, number FROM routes WHERE id = \$1`).
			WithArgs("nonexistent").
			WillReturnError(sql.ErrNoRows)

		err = repo.UpdateById(&models.Route{ID: "nonexistent"})
		if err == nil || err.Error() != "Route not found" {
			t.Errorf("Ожидалась ошибка 'Route not found', получена: %v", err)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Не все ожидаемые SQL запросы были выполнены: %v", err)
		}
	})

	t.Run("AssignDriver", func(t *testing.T) {
		db, mock, repo := setupMockRoute(t)
		defer db.Close()

		routeID := uuid.New().String()
		driverID := uuid.New().String()

		mock.ExpectQuery(`SELECT id, number FROM routes WHERE id = \$1`).
			WithArgs(routeID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "number"}).
				AddRow(routeID, "108"))

		mock.ExpectQuery(`SELECT COUNT\(\*\) FROM routes_drivers WHERE route_id = \$1 AND driver_id = \$2`).
			WithArgs(routeID, driverID).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

		mock.ExpectExec(`INSERT into routes_drivers \(route_id, driver_id\) VALUES \(\$1, \$2\)`).
			WithArgs(routeID, driverID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.AssignDriver(routeID, driverID)
		if err != nil {
			t.Errorf("Ошибка при назначении водителя: %v", err)
		}

		mock.ExpectQuery(`SELECT id, number FROM routes WHERE id = \$1`).
			WithArgs("nonexistent").
			WillReturnError(sql.ErrNoRows)

		err = repo.AssignDriver("nonexistent", driverID)
		if err == nil || err.Error() != "Route not found" {
			t.Errorf("Ожидалась ошибка 'Route not found', получена: %v", err)
		}

		mock.ExpectQuery(`SELECT id, number FROM routes WHERE id = \$1`).
			WithArgs(routeID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "number"}).
				AddRow(routeID, "108"))

		mock.ExpectQuery(`SELECT COUNT\(\*\) FROM routes_drivers WHERE route_id = \$1 AND driver_id = \$2`).
			WithArgs(routeID, driverID).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

		err = repo.AssignDriver(routeID, driverID)
		if err == nil || err.Error() != "Pair route_id and driver_id already exists" {
			t.Errorf("Ожидалась ошибка 'Pair route_id and driver_id already exists', получена: %v", err)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Не все ожидаемые SQL запросы были выполнены: %v", err)
		}
	})

	t.Run("AssignBusStop", func(t *testing.T) {
		db, mock, repo := setupMockRoute(t)
		defer db.Close()

		routeID := uuid.New().String()
		busStopID := uuid.New().String()

		mock.ExpectQuery(`SELECT id, number FROM routes WHERE id = \$1`).
			WithArgs(routeID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "number"}).
				AddRow(routeID, "109"))

		mock.ExpectQuery(`SELECT COUNT\(\*\) FROM routes_bus_stops WHERE route_id = \$1 AND bus_stop_id = \$2`).
			WithArgs(routeID, busStopID).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

		mock.ExpectExec(`INSERT into routes_bus_stops \(route_id, bus_stop_id\) VALUES \(\$1, \$2\)`).
			WithArgs(routeID, busStopID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.AssignBusStop(routeID, busStopID)
		if err != nil {
			t.Errorf("Ошибка при назначении остановки: %v", err)
		}

		mock.ExpectQuery(`SELECT id, number FROM routes WHERE id = \$1`).
			WithArgs("nonexistent").
			WillReturnError(sql.ErrNoRows)

		err = repo.AssignBusStop("nonexistent", busStopID)
		if err == nil || err.Error() != "Route not found" {
			t.Errorf("Ожидалась ошибка 'Route not found', получена: %v", err)
		}

		mock.ExpectQuery(`SELECT id, number FROM routes WHERE id = \$1`).
			WithArgs(routeID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "number"}).
				AddRow(routeID, "109"))

		mock.ExpectQuery(`SELECT COUNT\(\*\) FROM routes_bus_stops WHERE route_id = \$1 AND bus_stop_id = \$2`).
			WithArgs(routeID, busStopID).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

		err = repo.AssignBusStop(routeID, busStopID)
		if err == nil || err.Error() != "Pair route_id and bus_stop_id already exists" {
			t.Errorf("Ожидалась ошибка 'Pair route_id and bus_stop_id already exists', получена: %v", err)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Не все ожидаемые SQL запросы были выполнены: %v", err)
		}
	})

	t.Run("AssignBus", func(t *testing.T) {
		db, mock, repo := setupMockRoute(t)
		defer db.Close()

		routeID := uuid.New().String()
		busID := uuid.New().String()

		mock.ExpectQuery(`SELECT id, number FROM routes WHERE id = \$1`).
			WithArgs(routeID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "number"}).
				AddRow(routeID, "110"))

		mock.ExpectQuery(`SELECT COUNT\(\*\) FROM routes_buses WHERE route_id = \$1 AND bus_id = \$2`).
			WithArgs(routeID, busID).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

		mock.ExpectExec(`INSERT into routes_buses \(route_id, bus_id\) VALUES \(\$1, \$2\)`).
			WithArgs(routeID, busID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.AssignBus(routeID, busID)
		if err != nil {
			t.Errorf("Ошибка при назначении автобуса: %v", err)
		}

		mock.ExpectQuery(`SELECT id, number FROM routes WHERE id = \$1`).
			WithArgs("nonexistent").
			WillReturnError(sql.ErrNoRows)

		err = repo.AssignBus("nonexistent", busID)
		if err == nil || err.Error() != "Route not found" {
			t.Errorf("Ожидалась ошибка 'Route not found', получена: %v", err)
		}

		mock.ExpectQuery(`SELECT id, number FROM routes WHERE id = \$1`).
			WithArgs(routeID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "number"}).
				AddRow(routeID, "110"))

		mock.ExpectQuery(`SELECT COUNT\(\*\) FROM routes_buses WHERE route_id = \$1 AND bus_id = \$2`).
			WithArgs(routeID, busID).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

		err = repo.AssignBus(routeID, busID)
		if err == nil || err.Error() != "Pair route_id and bus_id already exists" {
			t.Errorf("Ожидалась ошибка 'Pair route_id and bus_id already exists', получена: %v", err)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Не все ожидаемые SQL запросы были выполнены: %v", err)
		}
	})

	t.Run("UnassignBusStop", func(t *testing.T) {
		db, mock, repo := setupMockRoute(t)
		defer db.Close()

		routeID := uuid.New().String()
		busStopID := uuid.New().String()

		mock.ExpectQuery(`SELECT id, number FROM routes WHERE id = \$1`).
			WithArgs(routeID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "number"}).
				AddRow(routeID, "111"))

		mock.ExpectExec(`DELETE FROM routes_bus_stops WHERE route_id = \$1 AND bus_stop_id = \$2`).
			WithArgs(routeID, busStopID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.UnassignBusStop(routeID, busStopID)
		if err != nil {
			t.Errorf("Ошибка при снятии назначения остановки: %v", err)
		}

		mock.ExpectQuery(`SELECT id, number FROM routes WHERE id = \$1`).
			WithArgs("nonexistent").
			WillReturnError(sql.ErrNoRows)

		err = repo.UnassignBusStop("nonexistent", busStopID)
		if err == nil || err.Error() != "Route not found" {
			t.Errorf("Ожидалась ошибка 'Route not found', получена: %v", err)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Не все ожидаемые SQL запросы были выполнены: %v", err)
		}
	})

	t.Run("UnassignBus", func(t *testing.T) {
		db, mock, repo := setupMockRoute(t)
		defer db.Close()

		routeID := uuid.New().String()
		busID := uuid.New().String()

		mock.ExpectQuery(`SELECT id, number FROM routes WHERE id = \$1`).
			WithArgs(routeID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "number"}).
				AddRow(routeID, "112"))

		mock.ExpectExec(`DELETE FROM routes_buses WHERE route_id = \$1 AND bus_id = \$2`).
			WithArgs(routeID, busID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.UnassignBus(routeID, busID)
		if err != nil {
			t.Errorf("Ошибка при снятии назначения автобуса: %v", err)
		}

		mock.ExpectQuery(`SELECT id, number FROM routes WHERE id = \$1`).
			WithArgs("nonexistent").
			WillReturnError(sql.ErrNoRows)

		err = repo.UnassignBus("nonexistent", busID)
		if err == nil || err.Error() != "Route not found" {
			t.Errorf("Ожидалась ошибка 'Route not found', получена: %v", err)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Не все ожидаемые SQL запросы были выполнены: %v", err)
		}
	})

	t.Run("UnassignDriver", func(t *testing.T) {
		db, mock, repo := setupMockRoute(t)
		defer db.Close()

		routeID := uuid.New().String()
		driverID := uuid.New().String()

		mock.ExpectQuery(`SELECT id, number FROM routes WHERE id = \$1`).
			WithArgs(routeID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "number"}).
				AddRow(routeID, "113"))

		mock.ExpectExec(`DELETE FROM routes_drivers WHERE route_id = \$1 AND driver_id = \$2`).
			WithArgs(routeID, driverID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.UnassignDriver(routeID, driverID)
		if err != nil {
			t.Errorf("Ошибка при снятии назначения водителя: %v", err)
		}

		mock.ExpectQuery(`SELECT id, number FROM routes WHERE id = \$1`).
			WithArgs("nonexistent").
			WillReturnError(sql.ErrNoRows)

		err = repo.UnassignDriver("nonexistent", driverID)
		if err == nil || err.Error() != "Route not found" {
			t.Errorf("Ожидалась ошибка 'Route not found', получена: %v", err)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Не все ожидаемые SQL запросы были выполнены: %v", err)
		}
	})

	t.Run("GetAllDriversById", func(t *testing.T) {
		db, mock, repo := setupMockRoute(t)
		defer db.Close()

		routeID := uuid.New().String()
		driver1 := models.Driver{
			ID:             uuid.New().String(),
			Name:           "Иван",
			Surname:        "Иванов",
			Patronymic:     "Иванович",
			BirthDate:      time.Date(1980, 1, 1, 0, 0, 0, 0, time.UTC),
			PassportSeries: "1234 567890",
			Snils:          "123-456-789 01",
			LicenseSeries:  "AB1234567",
		}
		driver2 := models.Driver{
			ID:             uuid.New().String(),
			Name:           "Петр",
			Surname:        "Петров",
			Patronymic:     "Петрович",
			BirthDate:      time.Date(1985, 2, 2, 0, 0, 0, 0, time.UTC),
			PassportSeries: "9876 543210",
			Snils:          "987-654-321 02",
			LicenseSeries:  "CD9876543",
		}

		mock.ExpectQuery(`SELECT id, number FROM routes WHERE id = \$1`).
			WithArgs(routeID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "number"}).
				AddRow(routeID, "114"))

		rows := sqlmock.NewRows([]string{"id", "name", "surname", "patronymic", "birth_date", "passport_series", "snils", "license_series"}).
			AddRow(driver1.ID, driver1.Name, driver1.Surname, driver1.Patronymic, driver1.BirthDate, driver1.PassportSeries, driver1.Snils, driver1.LicenseSeries).
			AddRow(driver2.ID, driver2.Name, driver2.Surname, driver2.Patronymic, driver2.BirthDate, driver2.PassportSeries, driver2.Snils, driver2.LicenseSeries)
		mock.ExpectQuery(`SELECT d\.id, d\.name, d\.surname, d\.patronymic, d\.birth_date, d\.passport_series, d\.snils, d\.license_series FROM drivers d JOIN routes_drivers rd ON d\.id = rd\.driver_id WHERE rd\.route_id=\$1`).
			WithArgs(routeID).
			WillReturnRows(rows)

		drivers, err := repo.GetAllDriversById(routeID)
		if err != nil {
			t.Errorf("Ошибка при получении водителей по маршруту: %v", err)
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

		mock.ExpectQuery(`SELECT id, number FROM routes WHERE id = \$1`).
			WithArgs("nonexistent").
			WillReturnError(sql.ErrNoRows)

		_, err = repo.GetAllDriversById("nonexistent")
		if err == nil || err.Error() != "Route not found" {
			t.Errorf("Ожидалась ошибка 'Route not found', получена: %v", err)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Не все ожидаемые SQL запросы были выполнены: %v", err)
		}
	})

	t.Run("GetAllBusStopsById", func(t *testing.T) {
		db, mock, repo := setupMockRoute(t)
		defer db.Close()

		routeID := uuid.New().String()
		busStop1 := models.BusStop{
			ID:   uuid.New().String(),
			Lat:  51.5074,
			Long: -0.1278,
			Name: "Трафальгарская площадь",
		}
		busStop2 := models.BusStop{
			ID:   uuid.New().String(),
			Lat:  48.8566,
			Long: 2.3522,
			Name: "Эйфелева башня",
		}

		mock.ExpectQuery(`SELECT id, number FROM routes WHERE id = \$1`).
			WithArgs(routeID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "number"}).
				AddRow(routeID, "115"))

		rows := sqlmock.NewRows([]string{"id", "lat", "long", "name"}).
			AddRow(busStop1.ID, busStop1.Lat, busStop1.Long, busStop1.Name).
			AddRow(busStop2.ID, busStop2.Lat, busStop2.Long, busStop2.Name)
		mock.ExpectQuery(`SELECT d\.id, d\.lat, d\.long, d\.name FROM bus_stops d JOIN routes_bus_stops rd ON d\.id = rd\.bus_stop_id WHERE rd\.route_id=\$1`).
			WithArgs(routeID).
			WillReturnRows(rows)

		busStops, err := repo.GetAllBusStopsById(routeID)
		if err != nil {
			t.Errorf("Ошибка при получении остановок по маршруту: %v", err)
		}
		if len(busStops) != 2 {
			t.Errorf("Ожидалось 2 остановки, получено: %d", len(busStops))
		}
		foundBusStop1, foundBusStop2 := false, false
		for _, s := range busStops {
			if reflect.DeepEqual(s, busStop1) {
				foundBusStop1 = true
			}
			if reflect.DeepEqual(s, busStop2) {
				foundBusStop2 = true
			}
		}
		if !foundBusStop1 || !foundBusStop2 {
			t.Errorf("Не все остановки найдены в списке: busStop1=%v, busStop2=%v", foundBusStop1, foundBusStop2)
		}

		mock.ExpectQuery(`SELECT id, number FROM routes WHERE id = \$1`).
			WithArgs("nonexistent").
			WillReturnError(sql.ErrNoRows)

		_, err = repo.GetAllBusStopsById("nonexistent")
		if err == nil || err.Error() != "Route not found" {
			t.Errorf("Ожидалась ошибка 'Route not found', получена: %v", err)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Не все ожидаемые SQL запросы были выполнены: %v", err)
		}
	})

	t.Run("GetAllBusesById", func(t *testing.T) {
		db, mock, repo := setupMockRoute(t)
		defer db.Close()

		routeID := uuid.New().String()
		bus1 := models.Bus{
			ID:             uuid.New().String(),
			Brand:          "Volvo",
			BusModel:       "B7R",
			RegisterNumber: "A123BC",
			AssemblyDate:   time.Date(2015, 1, 1, 0, 0, 0, 0, time.UTC),
			LastRepairDate: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		}
		bus2 := models.Bus{
			ID:             uuid.New().String(),
			Brand:          "MAN",
			BusModel:       "Lion's City",
			RegisterNumber: "B456DE",
			AssemblyDate:   time.Date(2018, 2, 2, 0, 0, 0, 0, time.UTC),
			LastRepairDate: time.Date(2024, 3, 3, 0, 0, 0, 0, time.UTC),
		}

		mock.ExpectQuery(`SELECT id, number FROM routes WHERE id = \$1`).
			WithArgs(routeID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "number"}).
				AddRow(routeID, "116"))

		rows := sqlmock.NewRows([]string{"id", "brand", "bus_model", "register_number", "assembly_date", "last_repair_date"}).
			AddRow(bus1.ID, bus1.Brand, bus1.BusModel, bus1.RegisterNumber, bus1.AssemblyDate, bus1.LastRepairDate).
			AddRow(bus2.ID, bus2.Brand, bus2.BusModel, bus2.RegisterNumber, bus2.AssemblyDate, bus2.LastRepairDate)
		mock.ExpectQuery(`SELECT d\.id, d\.brand, d\.bus_model, d\.register_number, d\.assembly_date, d\.last_repair_date FROM buses d JOIN routes_buses rd ON d\.id = rd\.bus_id WHERE rd\.route_id=\$1`).
			WithArgs(routeID).
			WillReturnRows(rows)

		buses, err := repo.GetAllBusesById(routeID)
		if err != nil {
			t.Errorf("Ошибка при получении автобусов по маршруту: %v", err)
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

		mock.ExpectQuery(`SELECT id, number FROM routes WHERE id = \$1`).
			WithArgs("nonexistent").
			WillReturnError(sql.ErrNoRows)

		_, err = repo.GetAllBusesById("nonexistent")
		if err == nil || err.Error() != "Route not found" {
			t.Errorf("Ожидалась ошибка 'Route not found', получена: %v", err)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Не все ожидаемые SQL запросы были выполнены: %v", err)
		}
	})
}

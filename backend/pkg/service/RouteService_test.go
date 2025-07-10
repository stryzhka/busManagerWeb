package service

import (
	"backend/pkg/models"
	"errors"
	"github.com/google/uuid"
	"testing"
	"time"
)

type MockRouteRepository struct {
	getByIdResp            *models.Route
	getByIdErr             error
	getByNumberResp        *models.Route
	getByNumberErr         error
	addErr                 error
	getAllResp             []models.Route
	getAllErr              error
	deleteByIdErr          error
	updateByIdErr          error
	assignDriverErr        error
	assignBusStopErr       error
	assignBusErr           error
	unassignBusStopErr     error
	unassignBusErr         error
	unassignDriverErr      error
	getAllDriversByIdResp  []models.Driver
	getAllDriversByIdErr   error
	getAllBusStopsByIdResp []models.BusStop
	getAllBusStopsByIdErr  error
	getAllBusesByIdResp    []models.Bus
	getAllBusesByIdErr     error
}

func (m *MockRouteRepository) GetById(id string) (*models.Route, error) {
	return m.getByIdResp, m.getByIdErr
}

func (m *MockRouteRepository) GetByNumber(number string) (*models.Route, error) {
	return m.getByNumberResp, m.getByNumberErr
}

func (m *MockRouteRepository) Add(route *models.Route) error {
	return m.addErr
}

func (m *MockRouteRepository) GetAll() ([]models.Route, error) {
	return m.getAllResp, m.getAllErr
}

func (m *MockRouteRepository) DeleteById(id string) error {
	return m.deleteByIdErr
}

func (m *MockRouteRepository) UpdateById(route *models.Route) error {
	return m.updateByIdErr
}

func (m *MockRouteRepository) AssignDriver(routeId, driverId string) error {
	return m.assignDriverErr
}

func (m *MockRouteRepository) AssignBusStop(routeId, busStopId string) error {
	return m.assignBusStopErr
}

func (m *MockRouteRepository) AssignBus(routeId, busId string) error {
	return m.assignBusErr
}

func (m *MockRouteRepository) UnassignBusStop(routeId, busStopId string) error {
	return m.unassignBusStopErr
}

func (m *MockRouteRepository) UnassignBus(routeId, busId string) error {
	return m.unassignBusErr
}

func (m *MockRouteRepository) UnassignDriver(routeId, driverId string) error {
	return m.unassignDriverErr
}

func (m *MockRouteRepository) GetAllDriversById(routeId string) ([]models.Driver, error) {
	return m.getAllDriversByIdResp, m.getAllDriversByIdErr
}

func (m *MockRouteRepository) GetAllBusStopsById(routeId string) ([]models.BusStop, error) {
	return m.getAllBusStopsByIdResp, m.getAllBusStopsByIdErr
}

func (m *MockRouteRepository) GetAllBusesById(routeId string) ([]models.Bus, error) {
	return m.getAllBusesByIdResp, m.getAllBusesByIdErr
}

type MockBusRepository struct {
	getByIdResp     *models.Bus
	getByIdErr      error
	getByNumberResp *models.Bus
	getByNumberErr  error
	addErr          error
	deleteByIdErr   error
	getAllResp      []models.Bus
	updateByIdErr   error
}

func (m *MockBusRepository) GetById(id string) (*models.Bus, error) {
	return m.getByIdResp, m.getByIdErr
}

func (m *MockBusRepository) GetByNumber(number string) (*models.Bus, error) {
	return m.getByNumberResp, m.getByNumberErr
}

func (m *MockBusRepository) Add(bus *models.Bus) error {
	return m.addErr
}

func (m *MockBusRepository) DeleteById(id string) error {
	return m.deleteByIdErr
}

func (m *MockBusRepository) GetAll() ([]models.Bus, error) {
	return m.getAllResp, nil
}

func (m *MockBusRepository) UpdateById(bus *models.Bus) error {
	return m.updateByIdErr
}

func TestRouteService_GetById(t *testing.T) {
	route := &models.Route{ID: uuid.New().String(), Number: "101"}

	t.Run("Success", func(t *testing.T) {
		mockRepo := &MockRouteRepository{getByIdResp: route}
		service := NewRouteService(mockRepo, nil, nil, nil)

		result, err := service.GetById(route.ID)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if result == nil || result.ID != route.ID {
			t.Errorf("Expected route with ID %s, got %v", route.ID, result)
		}
	})

	t.Run("Route not found", func(t *testing.T) {
		mockRepo := &MockRouteRepository{getByIdErr: errors.New("Route not found")}
		service := NewRouteService(mockRepo, nil, nil, nil)

		_, err := service.GetById(uuid.New().String())
		if err == nil || err.Error() != "Route not found" {
			t.Errorf("Expected 'Route not found' error, got %v", err)
		}
	})
}

func TestRouteService_GetByNumber(t *testing.T) {
	route := &models.Route{ID: uuid.New().String(), Number: "101"}

	t.Run("Success", func(t *testing.T) {
		mockRepo := &MockRouteRepository{getByNumberResp: route}
		service := NewRouteService(mockRepo, nil, nil, nil)

		result, err := service.GetByNumber("101")
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if result == nil || result.Number != route.Number {
			t.Errorf("Expected route with number %s, got %v", route.Number, result)
		}
	})

	t.Run("Route not found", func(t *testing.T) {
		mockRepo := &MockRouteRepository{getByNumberErr: errors.New("Route not found")}
		service := NewRouteService(mockRepo, nil, nil, nil)

		_, err := service.GetByNumber("999")
		if err == nil || err.Error() != "Route not found" {
			t.Errorf("Expected 'Route not found' error, got %v", err)
		}
	})
}

func TestRouteService_Add(t *testing.T) {
	route := &models.Route{Number: "101"}

	t.Run("Success", func(t *testing.T) {
		mockRepo := &MockRouteRepository{}
		service := NewRouteService(mockRepo, nil, nil, nil)

		err := service.Add(route)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})

	t.Run("Add with error", func(t *testing.T) {
		mockRepo := &MockRouteRepository{addErr: errors.New("Database error")}
		service := NewRouteService(mockRepo, nil, nil, nil)

		err := service.Add(route)
		if err == nil || err.Error() != "Database error" {
			t.Errorf("Expected 'Database error', got %v", err)
		}
	})
}

func TestRouteService_GetAll(t *testing.T) {
	route1 := models.Route{ID: uuid.New().String(), Number: "101"}
	route2 := models.Route{ID: uuid.New().String(), Number: "102"}

	t.Run("Success", func(t *testing.T) {
		mockRepo := &MockRouteRepository{getAllResp: []models.Route{route1, route2}}
		service := NewRouteService(mockRepo, nil, nil, nil)

		routes, _ := service.GetAll()
		if len(routes) != 2 {
			t.Errorf("Expected 2 routes, got %d", len(routes))
		}
	})

	t.Run("Empty result", func(t *testing.T) {
		mockRepo := &MockRouteRepository{getAllResp: []models.Route{}}
		service := NewRouteService(mockRepo, nil, nil, nil)

		routes, _ := service.GetAll()
		if len(routes) != 0 {
			t.Errorf("Expected 0 routes, got %d", len(routes))
		}
	})
}

func TestRouteService_DeleteById(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockRepo := &MockRouteRepository{}
		service := NewRouteService(mockRepo, nil, nil, nil)

		err := service.DeleteById(uuid.New().String())
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})

	t.Run("Delete with error", func(t *testing.T) {
		mockRepo := &MockRouteRepository{deleteByIdErr: errors.New("Database error")}
		service := NewRouteService(mockRepo, nil, nil, nil)

		err := service.DeleteById(uuid.New().String())
		if err == nil || err.Error() != "Database error" {
			t.Errorf("Expected 'Database error', got %v", err)
		}
	})
}

func TestRouteService_UpdateById(t *testing.T) {
	route := &models.Route{ID: uuid.New().String(), Number: "101"}

	t.Run("Success", func(t *testing.T) {
		mockRepo := &MockRouteRepository{}
		service := NewRouteService(mockRepo, nil, nil, nil)

		err := service.UpdateById(route)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})

	t.Run("Update with error", func(t *testing.T) {
		mockRepo := &MockRouteRepository{updateByIdErr: errors.New("Database error")}
		service := NewRouteService(mockRepo, nil, nil, nil)

		err := service.UpdateById(route)
		if err == nil || err.Error() != "Database error" {
			t.Errorf("Expected 'Database error', got %v", err)
		}
	})
}

func TestRouteService_AssignDriver(t *testing.T) {
	routeID := uuid.New().String()
	driverID := uuid.New().String()
	route := &models.Route{ID: routeID, Number: "101"}
	birthDate := time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC)
	driver := &models.Driver{ID: driverID, Name: "John", Surname: "Doe", BirthDate: birthDate}

	t.Run("Success", func(t *testing.T) {
		mockRouteRepo := &MockRouteRepository{getByIdResp: route}
		mockDriverRepo := &MockDriverRepository{getByIdResp: driver}
		service := NewRouteService(mockRouteRepo, mockDriverRepo, nil, nil)

		err := service.AssignDriver(routeID, driverID)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})

	t.Run("Route not found", func(t *testing.T) {
		mockRouteRepo := &MockRouteRepository{getByIdErr: errors.New("Route not found")}
		mockDriverRepo := &MockDriverRepository{getByIdResp: driver}
		service := NewRouteService(mockRouteRepo, mockDriverRepo, nil, nil)

		err := service.AssignDriver(uuid.New().String(), driverID)
		if err == nil || err.Error() != "Route not found" {
			t.Errorf("Expected 'Route not found' error, got %v", err)
		}
	})

	t.Run("Driver not found", func(t *testing.T) {
		mockRouteRepo := &MockRouteRepository{getByIdResp: route}
		mockDriverRepo := &MockDriverRepository{getByIdErr: errors.New("Driver not found")}
		service := NewRouteService(mockRouteRepo, mockDriverRepo, nil, nil)

		err := service.AssignDriver(routeID, uuid.New().String())
		if err == nil || err.Error() != "Driver not found" {
			t.Errorf("Expected 'Driver not found' error, got %v", err)
		}
	})

	t.Run("Assign driver with repo error", func(t *testing.T) {
		mockRouteRepo := &MockRouteRepository{getByIdResp: route, assignDriverErr: errors.New("Database error")}
		mockDriverRepo := &MockDriverRepository{getByIdResp: driver}
		service := NewRouteService(mockRouteRepo, mockDriverRepo, nil, nil)

		err := service.AssignDriver(routeID, driverID)
		if err == nil || err.Error() != "Database error" {
			t.Errorf("Expected 'Database error', got %v", err)
		}
	})
}

func TestRouteService_AssignBusStop(t *testing.T) {
	routeID := uuid.New().String()
	busStopID := uuid.New().String()
	route := &models.Route{ID: routeID, Number: "101"}
	busStop := &models.BusStop{ID: busStopID, Name: "Stop A", Lat: 55.7558, Long: 37.6173}

	t.Run("Success", func(t *testing.T) {
		mockRouteRepo := &MockRouteRepository{getByIdResp: route}
		mockBusStopRepo := &MockBusStopRepository{getByIdResp: busStop}
		service := NewRouteService(mockRouteRepo, nil, nil, mockBusStopRepo)

		err := service.AssignBusStop(routeID, busStopID)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})

	t.Run("Route not found", func(t *testing.T) {
		mockRouteRepo := &MockRouteRepository{getByIdErr: errors.New("Route not found")}
		mockBusStopRepo := &MockBusStopRepository{getByIdResp: busStop}
		service := NewRouteService(mockRouteRepo, nil, nil, mockBusStopRepo)

		err := service.AssignBusStop(uuid.New().String(), busStopID)
		if err == nil || err.Error() != "Route not found" {
			t.Errorf("Expected 'Route not found' error, got %v", err)
		}
	})

	t.Run("Bus stop not found", func(t *testing.T) {
		mockRouteRepo := &MockRouteRepository{getByIdResp: route}
		mockBusStopRepo := &MockBusStopRepository{getByIdErr: errors.New("Bus stop not found")}
		service := NewRouteService(mockRouteRepo, nil, nil, mockBusStopRepo)

		err := service.AssignBusStop(routeID, uuid.New().String())
		if err == nil || err.Error() != "Bus stop not found" {
			t.Errorf("Expected 'Bus stop not found' error, got %v", err)
		}
	})

	t.Run("Assign bus stop with repo error", func(t *testing.T) {
		mockRouteRepo := &MockRouteRepository{getByIdResp: route, assignBusStopErr: errors.New("Database error")}
		mockBusStopRepo := &MockBusStopRepository{getByIdResp: busStop}
		service := NewRouteService(mockRouteRepo, nil, nil, mockBusStopRepo)

		err := service.AssignBusStop(routeID, busStopID)
		if err == nil || err.Error() != "Database error" {
			t.Errorf("Expected 'Database error', got %v", err)
		}
	})
}

func TestRouteService_AssignBus(t *testing.T) {
	routeID := uuid.New().String()
	busID := uuid.New().String()
	route := &models.Route{ID: routeID, Number: "101"}
	bus := &models.Bus{
		ID:             busID,
		Brand:          "Mercedes",
		BusModel:       "Citaro",
		RegisterNumber: "X123YZ",
		AssemblyDate:   time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		LastRepairDate: time.Date(2023, 6, 1, 0, 0, 0, 0, time.UTC),
	}

	t.Run("Success", func(t *testing.T) {
		mockRouteRepo := &MockRouteRepository{getByIdResp: route}
		mockBusRepo := &MockBusRepository{getByIdResp: bus}
		service := NewRouteService(mockRouteRepo, nil, mockBusRepo, nil)

		err := service.AssignBus(routeID, busID)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})

	t.Run("Route not found", func(t *testing.T) {
		mockRouteRepo := &MockRouteRepository{getByIdErr: errors.New("Route not found")}
		mockBusRepo := &MockBusRepository{getByIdResp: bus}
		service := NewRouteService(mockRouteRepo, nil, mockBusRepo, nil)

		err := service.AssignBus(uuid.New().String(), busID)
		if err == nil || err.Error() != "Route not found" {
			t.Errorf("Expected 'Route not found' error, got %v", err)
		}
	})

	t.Run("Bus not found", func(t *testing.T) {
		mockRouteRepo := &MockRouteRepository{getByIdResp: route}
		mockBusRepo := &MockBusRepository{getByIdErr: errors.New("Bus not found")}
		service := NewRouteService(mockRouteRepo, nil, mockBusRepo, nil)

		err := service.AssignBus(routeID, uuid.New().String())
		if err == nil || err.Error() != "Bus not found" {
			t.Errorf("Expected 'Bus not found' error, got %v", err)
		}
	})

	t.Run("Assign bus with repo error", func(t *testing.T) {
		mockRouteRepo := &MockRouteRepository{getByIdResp: route, assignBusErr: errors.New("Database error")}
		mockBusRepo := &MockBusRepository{getByIdResp: bus}
		service := NewRouteService(mockRouteRepo, nil, mockBusRepo, nil)

		err := service.AssignBus(routeID, busID)
		if err == nil || err.Error() != "Database error" {
			t.Errorf("Expected 'Database error', got %v", err)
		}
	})
}

func TestRouteService_UnassignDriver(t *testing.T) {
	routeID := uuid.New().String()
	driverID := uuid.New().String()
	route := &models.Route{ID: routeID, Number: "101"}
	birthDate := time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC)
	driver := &models.Driver{ID: driverID, Name: "John", Surname: "Doe", BirthDate: birthDate}

	t.Run("Success", func(t *testing.T) {
		mockRouteRepo := &MockRouteRepository{getByIdResp: route}
		mockDriverRepo := &MockDriverRepository{getByIdResp: driver}
		service := NewRouteService(mockRouteRepo, mockDriverRepo, nil, nil)

		err := service.UnassignDriver(routeID, driverID)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})

	t.Run("Route not found", func(t *testing.T) {
		mockRouteRepo := &MockRouteRepository{getByIdErr: errors.New("Route not found")}
		mockDriverRepo := &MockDriverRepository{getByIdResp: driver}
		service := NewRouteService(mockRouteRepo, mockDriverRepo, nil, nil)

		err := service.UnassignDriver(uuid.New().String(), driverID)
		if err == nil || err.Error() != "Route not found" {
			t.Errorf("Expected 'Route not found' error, got %v", err)
		}
	})

	t.Run("Driver not found", func(t *testing.T) {
		mockRouteRepo := &MockRouteRepository{getByIdResp: route}
		mockDriverRepo := &MockDriverRepository{getByIdErr: errors.New("Driver not found")}
		service := NewRouteService(mockRouteRepo, mockDriverRepo, nil, nil)

		err := service.UnassignDriver(routeID, uuid.New().String())
		if err == nil || err.Error() != "Driver not found" {
			t.Errorf("Expected 'Driver not found' error, got %v", err)
		}
	})

	t.Run("Unassign driver with repo error", func(t *testing.T) {
		mockRouteRepo := &MockRouteRepository{getByIdResp: route, unassignDriverErr: errors.New("Database error")}
		mockDriverRepo := &MockDriverRepository{getByIdResp: driver}
		service := NewRouteService(mockRouteRepo, mockDriverRepo, nil, nil)

		err := service.UnassignDriver(routeID, driverID)
		if err == nil || err.Error() != "Database error" {
			t.Errorf("Expected 'Database error', got %v", err)
		}
	})
}

func TestRouteService_UnassignBusStop(t *testing.T) {
	routeID := uuid.New().String()
	busStopID := uuid.New().String()
	route := &models.Route{ID: routeID, Number: "101"}
	busStop := &models.BusStop{ID: busStopID, Name: "Stop A", Lat: 55.7558, Long: 37.6173}

	t.Run("Success", func(t *testing.T) {
		mockRouteRepo := &MockRouteRepository{getByIdResp: route}
		mockBusStopRepo := &MockBusStopRepository{getByIdResp: busStop}
		service := NewRouteService(mockRouteRepo, nil, nil, mockBusStopRepo)

		err := service.UnassignBusStop(routeID, busStopID)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})

	t.Run("Route not found", func(t *testing.T) {
		mockRouteRepo := &MockRouteRepository{getByIdErr: errors.New("Route not found")}
		mockBusStopRepo := &MockBusStopRepository{getByIdResp: busStop}
		service := NewRouteService(mockRouteRepo, nil, nil, mockBusStopRepo)

		err := service.UnassignBusStop(uuid.New().String(), busStopID)
		if err == nil || err.Error() != "Route not found" {
			t.Errorf("Expected 'Route not found' error, got %v", err)
		}
	})

	t.Run("Bus stop not found", func(t *testing.T) {
		mockRouteRepo := &MockRouteRepository{getByIdResp: route}
		mockBusStopRepo := &MockBusStopRepository{getByIdErr: errors.New("Bus stop not found")}
		service := NewRouteService(mockRouteRepo, nil, nil, mockBusStopRepo)

		err := service.UnassignBusStop(routeID, uuid.New().String())
		if err == nil || err.Error() != "Bus stop not found" {
			t.Errorf("Expected 'Bus stop not found' error, got %v", err)
		}
	})

	t.Run("Unassign bus stop with repo error", func(t *testing.T) {
		mockRouteRepo := &MockRouteRepository{getByIdResp: route, unassignBusStopErr: errors.New("Database error")}
		mockBusStopRepo := &MockBusStopRepository{getByIdResp: busStop}
		service := NewRouteService(mockRouteRepo, nil, nil, mockBusStopRepo)

		err := service.UnassignBusStop(routeID, busStopID)
		if err == nil || err.Error() != "Database error" {
			t.Errorf("Expected 'Database error', got %v", err)
		}
	})
}

func TestRouteService_UnassignBus(t *testing.T) {
	routeID := uuid.New().String()
	busID := uuid.New().String()
	route := &models.Route{ID: routeID, Number: "101"}
	bus := &models.Bus{
		ID:             busID,
		Brand:          "Mercedes",
		BusModel:       "Citaro",
		RegisterNumber: "X123YZ",
		AssemblyDate:   time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		LastRepairDate: time.Date(2023, 6, 1, 0, 0, 0, 0, time.UTC),
	}

	t.Run("Success", func(t *testing.T) {
		mockRouteRepo := &MockRouteRepository{getByIdResp: route}
		mockBusRepo := &MockBusRepository{getByIdResp: bus}
		service := NewRouteService(mockRouteRepo, nil, mockBusRepo, nil)

		err := service.UnassignBus(routeID, busID)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})

	t.Run("Route not found", func(t *testing.T) {
		mockRouteRepo := &MockRouteRepository{getByIdErr: errors.New("Route not found")}
		mockBusRepo := &MockBusRepository{getByIdResp: bus}
		service := NewRouteService(mockRouteRepo, nil, mockBusRepo, nil)

		err := service.UnassignBus(uuid.New().String(), busID)
		if err == nil || err.Error() != "Route not found" {
			t.Errorf("Expected 'Route not found' error, got %v", err)
		}
	})

	t.Run("Bus not found", func(t *testing.T) {
		mockRouteRepo := &MockRouteRepository{getByIdResp: route}
		mockBusRepo := &MockBusRepository{getByIdErr: errors.New("Bus not found")}
		service := NewRouteService(mockRouteRepo, nil, mockBusRepo, nil)

		err := service.UnassignBus(routeID, uuid.New().String())
		if err == nil || err.Error() != "Bus not found" {
			t.Errorf("Expected 'Bus not found' error, got %v", err)
		}
	})

	t.Run("Unassign bus with repo error", func(t *testing.T) {
		mockRouteRepo := &MockRouteRepository{getByIdResp: route, unassignBusErr: errors.New("Database error")}
		mockBusRepo := &MockBusRepository{getByIdResp: bus}
		service := NewRouteService(mockRouteRepo, nil, mockBusRepo, nil)

		err := service.UnassignBus(routeID, busID)
		if err == nil || err.Error() != "Database error" {
			t.Errorf("Expected 'Database error', got %v", err)
		}
	})
}

func TestRouteService_GetAllDriversById(t *testing.T) {
	routeID := uuid.New().String()
	birthDate := time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC)
	driver1 := models.Driver{ID: uuid.New().String(), Name: "John", Surname: "Doe", BirthDate: birthDate}
	driver2 := models.Driver{ID: uuid.New().String(), Name: "Jane", Surname: "Smith", BirthDate: birthDate}

	t.Run("Success", func(t *testing.T) {
		mockRepo := &MockRouteRepository{
			getByIdResp:           &models.Route{ID: routeID, Number: "101"},
			getAllDriversByIdResp: []models.Driver{driver1, driver2},
		}
		service := NewRouteService(mockRepo, nil, nil, nil)

		drivers, err := service.GetAllDriversById(routeID)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if len(drivers) != 2 {
			t.Errorf("Expected 2 drivers, got %d", len(drivers))
		}
		if drivers[0].Name != "John" || drivers[1].Name != "Jane" {
			t.Errorf("Expected names John and Jane, got %v", drivers)
		}
	})

	t.Run("Route not found", func(t *testing.T) {
		mockRepo := &MockRouteRepository{
			getByIdErr: errors.New("Route not found"),
		}
		service := NewRouteService(mockRepo, nil, nil, nil)

		_, err := service.GetAllDriversById(uuid.New().String())
		if err == nil || err.Error() != "Route not found" {
			t.Errorf("Expected 'Route not found' error, got %v", err)
		}
	})

	t.Run("Drivers not found", func(t *testing.T) {
		mockRepo := &MockRouteRepository{
			getByIdResp:           &models.Route{ID: routeID, Number: "101"},
			getAllDriversByIdResp: nil,
		}
		service := NewRouteService(mockRepo, nil, nil, nil)

		drivers, err := service.GetAllDriversById(routeID)
		if err == nil || err.Error() != "Drivers not found" {
			t.Errorf("Expected 'Drivers not found' error, got %v", err)
		}
		if drivers != nil {
			t.Errorf("Expected nil drivers, got %v", drivers)
		}
	})

	t.Run("Repo error", func(t *testing.T) {
		mockRepo := &MockRouteRepository{
			getByIdResp:          &models.Route{ID: routeID, Number: "101"},
			getAllDriversByIdErr: errors.New("Database error"),
		}
		service := NewRouteService(mockRepo, nil, nil, nil)

		_, err := service.GetAllDriversById(routeID)
		if err == nil || err.Error() != "Database error" {
			t.Errorf("Expected 'Database error', got %v", err)
		}
	})
}

func TestRouteService_GetAllBusStopsById(t *testing.T) {
	routeID := uuid.New().String()
	busStop1 := models.BusStop{ID: uuid.New().String(), Name: "Stop A", Lat: 55.7558, Long: 37.6173}
	busStop2 := models.BusStop{ID: uuid.New().String(), Name: "Stop B", Lat: 55.7539, Long: 37.6208}

	t.Run("Success", func(t *testing.T) {
		mockRepo := &MockRouteRepository{
			getByIdResp:            &models.Route{ID: routeID, Number: "101"},
			getAllBusStopsByIdResp: []models.BusStop{busStop1, busStop2},
		}
		service := NewRouteService(mockRepo, nil, nil, nil)

		busStops, err := service.GetAllBusStopsById(routeID)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if len(busStops) != 2 {
			t.Errorf("Expected 2 bus stops, got %d", len(busStops))
		}
		if busStops[0].Name != "Stop A" || busStops[1].Name != "Stop B" {
			t.Errorf("Expected names Stop A and Stop B, got %v", busStops)
		}
	})

	t.Run("Route not found", func(t *testing.T) {
		mockRepo := &MockRouteRepository{
			getByIdErr: errors.New("Route not found"),
		}
		service := NewRouteService(mockRepo, nil, nil, nil)

		_, err := service.GetAllBusStopsById(uuid.New().String())
		if err == nil || err.Error() != "Route not found" {
			t.Errorf("Expected 'Route not found' error, got %v", err)
		}
	})

	t.Run("Bus stops not found", func(t *testing.T) {
		mockRepo := &MockRouteRepository{
			getByIdResp:            &models.Route{ID: routeID, Number: "101"},
			getAllBusStopsByIdResp: nil,
		}
		service := NewRouteService(mockRepo, nil, nil, nil)

		busStops, err := service.GetAllBusStopsById(routeID)
		if err == nil || err.Error() != "Bus stops not found" {
			t.Errorf("Expected 'Bus stops not found' error, got %v", err)
		}
		if busStops != nil {
			t.Errorf("Expected nil bus stops, got %v", busStops)
		}
	})

	t.Run("Repo error", func(t *testing.T) {
		mockRepo := &MockRouteRepository{
			getByIdResp:           &models.Route{ID: routeID, Number: "101"},
			getAllBusStopsByIdErr: errors.New("Database error"),
		}
		service := NewRouteService(mockRepo, nil, nil, nil)

		_, err := service.GetAllBusStopsById(routeID)
		if err == nil || err.Error() != "Database error" {
			t.Errorf("Expected 'Database error', got %v", err)
		}
	})
}

func TestRouteService_GetAllBusesById(t *testing.T) {
	routeID := uuid.New().String()
	assemblyDate := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	lastRepairDate := time.Date(2023, 6, 1, 0, 0, 0, 0, time.UTC)
	bus1 := models.Bus{
		ID:             uuid.New().String(),
		Brand:          "Mercedes",
		BusModel:       "Citaro",
		RegisterNumber: "X123YZ",
		AssemblyDate:   assemblyDate,
		LastRepairDate: lastRepairDate,
	}
	bus2 := models.Bus{
		ID:             uuid.New().String(),
		Brand:          "Volvo",
		BusModel:       "B8RLE",
		RegisterNumber: "Y456AB",
		AssemblyDate:   assemblyDate,
		LastRepairDate: lastRepairDate,
	}

	t.Run("Success", func(t *testing.T) {
		mockRepo := &MockRouteRepository{
			getByIdResp:         &models.Route{ID: routeID, Number: "101"},
			getAllBusesByIdResp: []models.Bus{bus1, bus2},
		}
		service := NewRouteService(mockRepo, nil, nil, nil)

		buses, err := service.GetAllBusesById(routeID)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if len(buses) != 2 {
			t.Errorf("Expected 2 buses, got %d", len(buses))
		}
		if buses[0].Brand != "Mercedes" || buses[1].Brand != "Volvo" {
			t.Errorf("Expected brands Mercedes and Volvo, got %v", buses)
		}
	})

	t.Run("Route not found", func(t *testing.T) {
		mockRepo := &MockRouteRepository{
			getByIdErr: errors.New("Route not found"),
		}
		service := NewRouteService(mockRepo, nil, nil, nil)

		_, err := service.GetAllBusesById(uuid.New().String())
		if err == nil || err.Error() != "Route not found" {
			t.Errorf("Expected 'Route not found' error, got %v", err)
		}
	})

	t.Run("Buses not found", func(t *testing.T) {
		mockRepo := &MockRouteRepository{
			getByIdResp:         &models.Route{ID: routeID, Number: "101"},
			getAllBusesByIdResp: nil,
		}
		service := NewRouteService(mockRepo, nil, nil, nil)

		buses, err := service.GetAllBusesById(routeID)
		if err == nil || err.Error() != "Buses not found" {
			t.Errorf("Expected 'Buses not found' error, got %v", err)
		}
		if buses != nil {
			t.Errorf("Expected nil buses, got %v", buses)
		}
	})

	t.Run("Repo error", func(t *testing.T) {
		mockRepo := &MockRouteRepository{
			getByIdResp:        &models.Route{ID: routeID, Number: "101"},
			getAllBusesByIdErr: errors.New("Database error"),
		}
		service := NewRouteService(mockRepo, nil, nil, nil)

		_, err := service.GetAllBusesById(routeID)
		if err == nil || err.Error() != "Database error" {
			t.Errorf("Expected 'Database error', got %v", err)
		}
	})
}

package service

import (
	"backend/pkg/models"
	"errors"
	"testing"
)

type MockBusStopRepository struct {
	getByIdResp   *models.BusStop
	getByIdErr    error
	getByNameResp *models.BusStop
	getByNameErr  error
	addErr        error
	getAllResp    []models.BusStop
	deleteByIdErr error
	updateByIdErr error
}

func (m *MockBusStopRepository) GetById(id string) (*models.BusStop, error) {
	return m.getByIdResp, m.getByIdErr
}

func (m *MockBusStopRepository) GetByName(name string) (*models.BusStop, error) {
	return m.getByNameResp, m.getByNameErr
}

func (m *MockBusStopRepository) Add(busStop *models.BusStop) error {
	return m.addErr
}

func (m *MockBusStopRepository) GetAll() ([]models.BusStop, error) {
	return m.getAllResp, nil
}

func (m *MockBusStopRepository) DeleteById(id string) error {
	return m.deleteByIdErr
}

func (m *MockBusStopRepository) UpdateById(busStop *models.BusStop) error {
	return m.updateByIdErr
}

func TestBusStopService_GetById(t *testing.T) {
	busStop := &models.BusStop{ID: "1", Lat: 55.7558, Long: 37.6173, Name: "Stop A"}

	t.Run("Success", func(t *testing.T) {
		mockRepo := &MockBusStopRepository{getByIdResp: busStop}
		service := NewBusStopService(mockRepo)

		result, err := service.GetById("1")
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if result == nil || result.ID != busStop.ID {
			t.Errorf("Expected bus stop with ID %s, got %v", busStop.ID, result)
		}
	})

	t.Run("Not found", func(t *testing.T) {
		mockRepo := &MockBusStopRepository{getByIdErr: errors.New("Bus stop not found")}
		service := NewBusStopService(mockRepo)

		_, err := service.GetById("2")
		if err == nil || err.Error() != "Bus stop not found" {
			t.Errorf("Expected 'Bus stop not found' error, got %v", err)
		}
	})
}

func TestBusStopService_GetByName(t *testing.T) {
	busStop := &models.BusStop{ID: "1", Lat: 55.7558, Long: 37.6173, Name: "Stop A"}

	t.Run("Success", func(t *testing.T) {
		mockRepo := &MockBusStopRepository{getByNameResp: busStop}
		service := NewBusStopService(mockRepo)

		result, err := service.GetByName("Stop A")
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if result == nil || result.Name != busStop.Name {
			t.Errorf("Expected bus stop with name %s, got %v", busStop.Name, result)
		}
	})

	t.Run("Not found", func(t *testing.T) {
		mockRepo := &MockBusStopRepository{getByNameErr: errors.New("Bus stop not found")}
		service := NewBusStopService(mockRepo)

		_, err := service.GetByName("Unknown")
		if err == nil || err.Error() != "Bus stop not found" {
			t.Errorf("Expected 'Bus stop not found' error, got %v", err)
		}
	})
}

func TestBusStopService_Add(t *testing.T) {
	busStop := &models.BusStop{Lat: 55.7558, Long: 37.6173, Name: "Stop A"}

	t.Run("Success", func(t *testing.T) {
		mockRepo := &MockBusStopRepository{}
		service := NewBusStopService(mockRepo)

		err := service.Add(busStop)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})

	t.Run("Add with error from repo", func(t *testing.T) {
		mockRepo := &MockBusStopRepository{addErr: errors.New("Database error")}
		service := NewBusStopService(mockRepo)

		err := service.Add(busStop)
		if err == nil || err.Error() != "Database error" {
			t.Errorf("Expected 'Database error', got %v", err)
		}
	})
}

func TestBusStopService_GetAll(t *testing.T) {
	busStop1 := models.BusStop{ID: "1", Lat: 55.7558, Long: 37.6173, Name: "Stop A"}
	busStop2 := models.BusStop{ID: "2", Lat: 55.7522, Long: 37.6156, Name: "Stop B"}

	t.Run("Get all bus stops", func(t *testing.T) {
		mockRepo := &MockBusStopRepository{getAllResp: []models.BusStop{busStop1, busStop2}}
		service := NewBusStopService(mockRepo)

		busStops, _ := service.GetAll()
		if len(busStops) != 2 {
			t.Errorf("Expected 2 bus stops, got %d", len(busStops))
		}
	})

	t.Run("Get all from empty repo", func(t *testing.T) {
		mockRepo := &MockBusStopRepository{getAllResp: []models.BusStop{}}
		service := NewBusStopService(mockRepo)

		busStops, _ := service.GetAll()
		if len(busStops) != 0 {
			t.Errorf("Expected 0 bus stops, got %d", len(busStops))
		}
	})
}

func TestBusStopService_DeleteById(t *testing.T) {
	t.Run("Delete existing bus stop", func(t *testing.T) {
		mockRepo := &MockBusStopRepository{}
		service := NewBusStopService(mockRepo)

		err := service.DeleteById("1")
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})

	t.Run("Delete with error from repo", func(t *testing.T) {
		mockRepo := &MockBusStopRepository{deleteByIdErr: errors.New("Database error")}
		service := NewBusStopService(mockRepo)

		err := service.DeleteById("1")
		if err == nil || err.Error() != "Database error" {
			t.Errorf("Expected 'Database error', got %v", err)
		}
	})
}

func TestBusStopService_UpdateById(t *testing.T) {
	busStop := &models.BusStop{ID: "1", Lat: 55.7522, Long: 37.6156, Name: "Stop B"}

	t.Run("Update existing bus stop", func(t *testing.T) {
		mockRepo := &MockBusStopRepository{}
		service := NewBusStopService(mockRepo)

		err := service.UpdateById(busStop)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})

	t.Run("Update with error from repo", func(t *testing.T) {
		mockRepo := &MockBusStopRepository{updateByIdErr: errors.New("Database error")}
		service := NewBusStopService(mockRepo)

		err := service.UpdateById(busStop)
		if err == nil || err.Error() != "Database error" {
			t.Errorf("Expected 'Database error', got %v", err)
		}
	})
}

package service

import (
	"backend/pkg/models"
	"errors"
	"github.com/google/uuid"

	"testing"
	"time"
)

type MockDriverRepository struct {
	//repository.IDriverRepository
	getByIdResp       *models.Driver
	getByIdErr        error
	getByPassportResp *models.Driver
	getByPassportErr  error
	addErr            error
	getAllResp        []models.Driver
	deleteByIdErr     error
	updateByIdErr     error
}

func (m *MockDriverRepository) GetById(id string) (*models.Driver, error) {
	return m.getByIdResp, m.getByIdErr
}

func (m *MockDriverRepository) GetByPassportSeries(series string) (*models.Driver, error) {
	return m.getByPassportResp, m.getByPassportErr
}

func (m *MockDriverRepository) Add(driver *models.Driver) error {
	return m.addErr
}

func (m *MockDriverRepository) GetAll() ([]models.Driver, error) {
	return m.getAllResp, nil
}

func (m *MockDriverRepository) DeleteById(id string) error {
	return m.deleteByIdErr
}

func (m *MockDriverRepository) UpdateById(driver *models.Driver) error {
	return m.updateByIdErr
}

func TestDriverService_GetById(t *testing.T) {
	fixedTime, _ := time.Parse(time.RFC3339, "2022-11-11T11:11:11Z")
	driver := &models.Driver{
		ID:             uuid.New().String(),
		Name:           "John",
		Surname:        "Doe",
		Patronymic:     "Ivanovich",
		BirthDate:      fixedTime,
		PassportSeries: "AB123456",
		Snils:          "123-456-789 00",
		LicenseSeries:  "CD789012",
	}

	t.Run("Get existing driver by ID", func(t *testing.T) {
		mockRepo := &MockDriverRepository{getByIdResp: driver}
		service := NewDriverService(mockRepo)

		result, err := service.GetById(driver.ID)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if result == nil || result.ID != driver.ID {
			t.Errorf("Expected driver with ID %s, got %v", driver.ID, result)
		}
	})

	t.Run("Get non-existent driver by ID", func(t *testing.T) {
		mockRepo := &MockDriverRepository{getByIdErr: errors.New("Driver not found")}
		service := NewDriverService(mockRepo)

		_, err := service.GetById(uuid.New().String())
		if err == nil || err.Error() != "Driver not found" {
			t.Errorf("Expected 'Driver not found' error, got %v", err)
		}
	})
}

func TestDriverService_GetByPassportSeries(t *testing.T) {
	fixedTime, _ := time.Parse(time.RFC3339, "2022-11-11T11:11:11Z")
	driver := &models.Driver{
		ID:             uuid.New().String(),
		Name:           "John",
		Surname:        "Doe",
		Patronymic:     "Ivanovich",
		BirthDate:      fixedTime,
		PassportSeries: "AB123456",
		Snils:          "123-456-789 00",
		LicenseSeries:  "CD789012",
	}

	t.Run("Get existing driver by passport series", func(t *testing.T) {
		mockRepo := &MockDriverRepository{getByPassportResp: driver}
		service := NewDriverService(mockRepo)

		result, err := service.GetByPassportSeries(driver.PassportSeries)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if result == nil || result.PassportSeries != driver.PassportSeries {
			t.Errorf("Expected driver with passport series %s, got %v", driver.PassportSeries, result)
		}
	})

	t.Run("Get non-existent driver by passport series", func(t *testing.T) {
		mockRepo := &MockDriverRepository{getByPassportErr: errors.New("Driver not found")}
		service := NewDriverService(mockRepo)

		_, err := service.GetByPassportSeries("XY789012")
		if err == nil || err.Error() != "Driver not found" {
			t.Errorf("Expected 'Driver not found' error, got %v", err)
		}
	})
}

func TestDriverService_Add(t *testing.T) {
	fixedTime, _ := time.Parse(time.RFC3339, "2022-11-11T11:11:11Z")
	driver := &models.Driver{
		Name:           "John",
		Surname:        "Doe",
		Patronymic:     "Ivanovich",
		BirthDate:      fixedTime,
		PassportSeries: "AB123456",
		Snils:          "123-456-789 00",
		LicenseSeries:  "CD789012",
	}

	t.Run("Add new driver", func(t *testing.T) {
		mockRepo := &MockDriverRepository{}
		service := NewDriverService(mockRepo)

		err := service.Add(driver)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})

	t.Run("Add with error from repo", func(t *testing.T) {
		mockRepo := &MockDriverRepository{addErr: errors.New("Database error")}
		service := NewDriverService(mockRepo)

		err := service.Add(driver)
		if err == nil || err.Error() != "Database error" {
			t.Errorf("Expected 'Database error', got %v", err)
		}
	})
}

func TestDriverService_GetAll(t *testing.T) {
	fixedTime, _ := time.Parse(time.RFC3339, "2022-11-11T11:11:11Z")
	driver1 := models.Driver{
		ID:             uuid.New().String(),
		Name:           "John",
		Surname:        "Doe",
		Patronymic:     "Ivanovich",
		BirthDate:      fixedTime,
		PassportSeries: "AB123456",
		Snils:          "123-456-789 00",
		LicenseSeries:  "CD789012",
	}
	driver2 := models.Driver{
		ID:             uuid.New().String(),
		Name:           "Jane",
		Surname:        "Doe",
		Patronymic:     "Ivanovna",
		BirthDate:      fixedTime,
		PassportSeries: "XY789012",
		Snils:          "987-654-321 00",
		LicenseSeries:  "EF345678",
	}

	t.Run("Get all drivers", func(t *testing.T) {
		mockRepo := &MockDriverRepository{getAllResp: []models.Driver{driver1, driver2}}
		service := NewDriverService(mockRepo)

		drivers := service.GetAll()
		if len(drivers) != 2 {
			t.Errorf("Expected 2 drivers, got %d", len(drivers))
		}
	})

	t.Run("Get all from empty repo", func(t *testing.T) {
		mockRepo := &MockDriverRepository{getAllResp: []models.Driver{}}
		service := NewDriverService(mockRepo)

		drivers := service.GetAll()
		if len(drivers) != 0 {
			t.Errorf("Expected 0 drivers, got %d", len(drivers))
		}
	})
}

func TestDriverService_DeleteById(t *testing.T) {
	t.Run("Delete existing driver", func(t *testing.T) {
		mockRepo := &MockDriverRepository{}
		service := NewDriverService(mockRepo)

		err := service.DeleteById(uuid.New().String())
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})

	t.Run("Delete with error from repo", func(t *testing.T) {
		mockRepo := &MockDriverRepository{deleteByIdErr: errors.New("Database error")}
		service := NewDriverService(mockRepo)

		err := service.DeleteById(uuid.New().String())
		if err == nil || err.Error() != "Database error" {
			t.Errorf("Expected 'Database error', got %v", err)
		}
	})
}

func TestDriverService_UpdateById(t *testing.T) {
	fixedTime, _ := time.Parse(time.RFC3339, "2022-11-11T11:11:11Z")
	driver := &models.Driver{
		ID:             uuid.New().String(),
		Name:           "John",
		Surname:        "Doe",
		Patronymic:     "Ivanovich",
		BirthDate:      fixedTime,
		PassportSeries: "AB123456",
		Snils:          "123-456-789 00",
		LicenseSeries:  "CD789012",
	}

	t.Run("Update existing driver", func(t *testing.T) {
		mockRepo := &MockDriverRepository{}
		service := NewDriverService(mockRepo)

		err := service.UpdateById(driver)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})

	t.Run("Update with error from repo", func(t *testing.T) {
		mockRepo := &MockDriverRepository{updateByIdErr: errors.New("Database error")}
		service := NewDriverService(mockRepo)

		err := service.UpdateById(driver)
		if err == nil || err.Error() != "Database error" {
			t.Errorf("Expected 'Database error', got %v", err)
		}
	})
}

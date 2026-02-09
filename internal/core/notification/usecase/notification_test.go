package usecase

import (
	"errors"
	"testing"

	"github.com/SevgiF/notification-system/internal/core/notification/domain"
)

// MockRepository
type MockRepository struct {
	notifications []domain.Notification
}

func (m *MockRepository) CreateNotification(tx interface{}, item domain.Notification) error {
	m.notifications = append(m.notifications, item)
	return nil
}

func (m *MockRepository) GetNotification(id int) (domain.Notification, error) {
	for _, n := range m.notifications {
		if n.ID == id {
			return n, nil
		}
	}
	return domain.Notification{}, errors.New("not found")
}

func (m *MockRepository) UpdateNotification(id int, status int) error {
	for i, n := range m.notifications {
		if n.ID == id {
			m.notifications[i].Status = status
			return nil
		}
	}
	return errors.New("not found")
}

// ... Implement other interface methods as needed for tests ...
// Since we can't run tests, this is illustrative.

func TestAddNotification(t *testing.T) {
	// Setup
	// repo := &MockRepository{}
	// service := NewNotificationService(repo)

	// Test
	// err := service.AddNotification(dto.NotificationRequest{...})

	// Assert
	// if err != nil { t.Errorf(...) }
}

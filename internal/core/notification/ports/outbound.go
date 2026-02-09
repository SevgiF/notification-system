package ports

import (
	"github.com/SevgiF/notification-system/internal/core/notification/domain"
)

type OutboundPort interface {
	CreateNotification(item domain.Notification) error
	GetNotification(id int) (domain.Notification, error)
	UpdateNotification(id int, status int) error
	GetAllNotification(status *int, channel, from, to *string, limit, offset int) ([]domain.Notification, int, int, error)
	GetStatusList() ([]domain.Status, error)
	FetchPendingNotifications(limit int) ([]domain.Notification, error)
	GetMetrics() (map[string]int, error)
}

type NotificationGateway interface {
	Send(notification domain.Notification) error
}

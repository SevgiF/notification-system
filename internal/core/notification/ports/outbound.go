package ports

import (
	"database/sql"
	"github.com/SevgiF/notification-system/internal/core/notification/domain"
)

type OutboundPort interface {
	WithTransaction() (*sql.Tx, error)
	CreateNotification(tx *sql.Tx, item domain.Notification) error
	GetNotification(id int) (domain.Notification, error)
	UpdateNotification(id int, status int) error
	GetAllNotification(status *int, channel, from, to *string, limit, offset int) ([]domain.Notification, int, int, error)
	GetStatusList() ([]domain.Status, error)
}

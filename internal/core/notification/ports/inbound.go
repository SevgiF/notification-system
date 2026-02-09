package ports

import "github.com/SevgiF/notification-system/internal/adapter/http/fiber/dto"

type InboundPort interface {
	AddNotification(item dto.NotificationRequest) error
	BulkAddNotification(items []dto.NotificationRequest) error
	NotificationDetail(id int) (dto.NotificationResponse, error)
	CancelNotification(id int) error
	NotificationList(status *int, channel *string, from *string, to *string, limitQuery *int, pageQuery *int) ([]dto.NotificationResponse, dto.Pagination, dto.Filter, error)
	StatusList() ([]dto.Status, error)
	GetMetrics() (map[string]interface{}, error)
}

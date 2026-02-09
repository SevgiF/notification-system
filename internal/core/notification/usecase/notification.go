package usecase

import (
	"log"

	"github.com/SevgiF/notification-system/internal/adapter/http/fiber/dto"
	"github.com/SevgiF/notification-system/internal/core/notification/domain"
	"github.com/SevgiF/notification-system/internal/core/notification/ports"
)

func NewNotificationService(repo ports.OutboundPort) *NotificationService {
	return &NotificationService{
		repo: repo,
	}
}

type NotificationService struct {
	repo ports.OutboundPort
}

func (s *NotificationService) GetMetrics() (map[string]interface{}, error) {
	metrics, err := s.repo.GetMetrics()
	if err != nil {
		return nil, err
	}

	// Convert checks:
	result := make(map[string]interface{})
	for k, v := range metrics {
		result[k] = v
	}
	return result, nil
}

func (s *NotificationService) StatusList() ([]dto.Status, error) {
	items, err := s.repo.GetStatusList()
	if err != nil {
		return nil, err
	}

	statuses := ToStatusResponseList(items)
	return statuses, nil
}

func (s *NotificationService) NotificationList(status *int, channel *string, from *string, to *string, limitQuery *int, pageQuery *int) ([]dto.NotificationResponse, dto.Pagination, dto.Filter, error) {
	limit := domain.DEFAULT_LIMIT
	page := 1

	if pageQuery != nil && *pageQuery > 0 {
		page = *pageQuery
	}

	if limitQuery != nil && domain.MIN_LIMIT <= *limitQuery && *limitQuery <= domain.MAX_LIMIT {
		limit = *limitQuery
	}

	offset := (page - 1) * limit

	items, totalCount, totalPage, err := s.repo.GetAllNotification(status, channel, from, to, limit, offset)
	if err != nil {
		log.Println(err.Error())
		return nil, dto.Pagination{}, dto.Filter{}, err
	}

	pagination := dto.Pagination{
		Limit:       limit,
		CurrentPage: page,
		TotalCount:  totalCount,
		TotalPage:   totalPage,
	}

	filter := dto.Filter{
		From:    from,
		To:      to,
		Status:  status,
		Channel: channel,
	}

	notifications := ToNotificationResponseList(items)
	return notifications, pagination, filter, nil
}

func (s *NotificationService) CancelNotification(id int) error {
	return s.repo.UpdateNotification(id, domain.CANCELED)
}

func (s *NotificationService) NotificationDetail(id int) (dto.NotificationResponse, error) {
	item, err := s.repo.GetNotification(id)
	if err != nil {
		return dto.NotificationResponse{}, err
	}

	notification := ToNotificationResponse(item)
	return notification, nil
}

func (s *NotificationService) BulkAddNotification(items []dto.NotificationRequest) (err error) {
	notifications := NotificationListFromRequest(items)

	for _, notification := range notifications {
		notification.Status = domain.CREATED
		err = s.repo.CreateNotification(notification)
		if err != nil {
			return
		}
	}
	return
}

func (s *NotificationService) AddNotification(item dto.NotificationRequest) (err error) {
	notification := NotificationFromRequest(item)
	notification.Status = domain.CREATED
	return s.repo.CreateNotification(notification)
}

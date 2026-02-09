package usecase

import (
	"github.com/SevgiF/notification-system/internal/adapter/http/fiber/dto"
	"github.com/SevgiF/notification-system/internal/core/notification/domain"
)

func ToNotificationResponse(n domain.Notification) dto.NotificationResponse {
	return dto.NotificationResponse{
		ID:        n.ID,
		Recipient: n.Recipient,
		Channel:   n.Channel,
		Content:   n.Content,
		Priority:  n.Priority,
		Status:    n.Status,
	}
}

func ToNotificationResponseList(notifications []domain.Notification) []dto.NotificationResponse {
	responses := make([]dto.NotificationResponse, 0, len(notifications))

	for _, n := range notifications {
		responses = append(responses, ToNotificationResponse(n))
	}

	return responses
}

func NotificationFromRequest(r dto.NotificationRequest) domain.Notification {
	return domain.Notification{
		Recipient: r.Recipient,
		Channel:   r.Channel,
		Content:   r.Content,
		Priority:  r.Priority,
	}
}

func NotificationListFromRequest(requests []dto.NotificationRequest) []domain.Notification {
	notifications := make([]domain.Notification, 0, len(requests))

	for _, n := range requests {
		notifications = append(notifications, NotificationFromRequest(n))
	}

	return notifications
}

func ToStatusResponse(s domain.Status) dto.Status {
	return dto.Status{
		Code:        s.Code,
		Description: s.Description,
	}
}

func ToStatusResponseList(statuses []domain.Status) []dto.Status {
	responses := make([]dto.Status, 0, len(statuses))

	for _, s := range statuses {
		responses = append(responses, ToStatusResponse(s))
	}

	return responses
}

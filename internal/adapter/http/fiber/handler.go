package fiber

import (
	"github.com/SevgiF/notification-system/internal/adapter/http/fiber/dto"
	"github.com/SevgiF/notification-system/internal/core/notification/ports"
	"github.com/gofiber/fiber/v3"
	"strconv"
)

func NewNotificationHandler(s ports.InboundPort) *NotificationHandler {
	return &NotificationHandler{s: s}
}

type NotificationHandler struct {
	s ports.InboundPort
}

func (h *NotificationHandler) AddNotification(c fiber.Ctx) error {
	var n dto.NotificationRequest

	if err := c.Bind().Body(&n); err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	err := h.s.AddNotification(n)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.SendStatus(fiber.StatusCreated)
}

func (h *NotificationHandler) BulkAddNotification(c fiber.Ctx) error {
	var n []dto.NotificationRequest

	if err := c.Bind().Body(&n); err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	err := h.s.BulkAddNotification(n)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.SendStatus(fiber.StatusCreated)
}

func (h *NotificationHandler) NotificationDetail(c fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	notification, err := h.s.NotificationDetail(id)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	response := dto.Response{
		Message: "",
		Data:    notification,
	}
	return c.Status(fiber.StatusOK).JSON(response)
}

func (h *NotificationHandler) CancelNotification(c fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	err = h.s.CancelNotification(id)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.SendStatus(fiber.StatusOK)
}

func (h *NotificationHandler) NotificationList(c fiber.Ctx) error {
	status, _ := strconv.Atoi(c.Query("status"))
	channel := c.Query("channel")
	from := c.Query("from")
	to := c.Query("to")
	limitQuery, _ := strconv.Atoi(c.Query("limit"))
	pageQuery, _ := strconv.Atoi(c.Query("page"))

	notifications, pagination, filters, err := h.s.NotificationList(&status, &channel, &from, &to, &limitQuery, &pageQuery)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	response := dto.Response{
		Message:    "",
		Pagination: pagination,
		Filter:     filters,
		Data:       notifications,
	}

	return c.Status(fiber.StatusOK).JSON(response)

}

func (h *NotificationHandler) StatusList(c fiber.Ctx) error {
	status, err := h.s.StatusList()
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	response := dto.Response{
		Message: "",
		Data:    status,
	}
	return c.Status(fiber.StatusOK).JSON(response)
}

package fiber

import (
	"log"
	"strconv"
	"time"

	"github.com/SevgiF/notification-system/internal/adapter/http/fiber/dto"
	"github.com/SevgiF/notification-system/internal/core/notification/ports"
	"github.com/gofiber/fiber/v3"
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
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	err := h.s.AddNotification(n)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
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

	log.Println("Getting notification list")
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
		Pagination: &pagination,
		Filter:     &filters,
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

func (h *NotificationHandler) GetMetrics(c fiber.Ctx) error {
	metrics, err := h.s.GetMetrics()
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.Status(fiber.StatusOK).JSON(metrics)
}

func (h *NotificationHandler) HealthCheck(c fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "ok",
		"timestamp": time.Now(),
	})
}

package notification

import (
	"database/sql"
	handler "github.com/SevgiF/notification-system/internal/adapter/http/fiber"
	"github.com/SevgiF/notification-system/internal/adapter/repository/mysql"
	"github.com/SevgiF/notification-system/internal/core/notification/usecase"
	"github.com/gofiber/fiber/v3"
)

func SetupNotification(app *fiber.App, db *sql.DB) {
	repository := mysql.NewNotificationRepository(db)
	service := usecase.NewNotificationService(repository)
	h := handler.NewNotificationHandler(service)

	app.Post("notification", h.AddNotification)
	app.Get("notification", h.NotificationList)
	app.Post("notification/batch", h.BulkAddNotification)
	app.Get("notification/:id", h.NotificationDetail)
	app.Put("notification/cancel/:id", h.CancelNotification)
	app.Get("status-list", h.StatusList)
}

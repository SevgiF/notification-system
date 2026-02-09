package notification

import (
	"database/sql"
	"github.com/SevgiF/notification-system/internal/adapter/gateway"
	handler "github.com/SevgiF/notification-system/internal/adapter/http/fiber"
	"github.com/SevgiF/notification-system/internal/adapter/repository/mysql"
	"github.com/SevgiF/notification-system/internal/core/notification/usecase"
	env "github.com/SevgiF/notification-system/pkg/environment"
	"github.com/gofiber/fiber/v3"
	"log"
)

func SetupNotification(app *fiber.App, db *sql.DB) *usecase.WorkerPool {
	repository := mysql.NewNotificationRepository(db)

	// Gateway
	webhookURL := env.GetEnvOrFail("WEBHOOK_URL")
	gateway := gateway.NewWebhookGateway(webhookURL)

	service := usecase.NewNotificationService(repository)
	h := handler.NewNotificationHandler(service)

	// Worker Pool (3 workers as default)
	workerPool := usecase.NewWorkerPool(repository, gateway, 3)

	log.Println("Setting up notification handler")
	app.Post("notification", h.AddNotification)
	app.Get("notification", h.NotificationList)
	app.Post("notification/batch", h.BulkAddNotification)
	app.Get("notification/:id", h.NotificationDetail)
	app.Put("notification/cancel/:id", h.CancelNotification)
	app.Get("status-list", h.StatusList)

	// Observability
	app.Get("metrics", h.GetMetrics)
	app.Get("health", h.HealthCheck)

	return workerPool
}

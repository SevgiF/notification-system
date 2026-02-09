package usecase

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/SevgiF/notification-system/internal/core/notification/domain"
	"github.com/SevgiF/notification-system/internal/core/notification/ports"
	"golang.org/x/time/rate"
)

type WorkerPool struct {
	repo         ports.OutboundPort
	gateway      ports.NotificationGateway
	workers      int
	rateLimiters map[string]*rate.Limiter
	quit         chan struct{}
	wg           sync.WaitGroup
}

func NewWorkerPool(repo ports.OutboundPort, gateway ports.NotificationGateway, workers int) *WorkerPool {
	return &WorkerPool{
		repo:         repo,
		gateway:      gateway,
		workers:      workers,
		rateLimiters: make(map[string]*rate.Limiter),
		quit:         make(chan struct{}),
	}
}

func (wp *WorkerPool) Start() {
	// Initialize rate limiters for each channel
	// 100 messages per second per channel
	wp.rateLimiters["sms"] = rate.NewLimiter(100, 1)
	wp.rateLimiters["email"] = rate.NewLimiter(100, 1)
	wp.rateLimiters["push"] = rate.NewLimiter(100, 1)

	for i := 0; i < wp.workers; i++ {
		wp.wg.Add(1)
		go wp.worker()
	}
}

func (wp *WorkerPool) Stop() {
	close(wp.quit)
	wp.wg.Wait()
}

func (wp *WorkerPool) worker() {
	defer wp.wg.Done()

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-wp.quit:
			return
		case <-ticker.C:
			wp.processNotifications()
		}
	}
}

func (wp *WorkerPool) processNotifications() {
	// Fetch pending notifications
	notifications, err := wp.repo.FetchPendingNotifications(wp.workers) // Fetch up to 'workers' count or batch size
	if err != nil {
		log.Printf("Error fetching notifications: %v", err)
		return
	}

	for _, n := range notifications {
		wp.wg.Add(1)
		go func(notification domain.Notification) {
			defer wp.wg.Done()
			wp.processNotification(notification)
		}(n)
	}
}

func (wp *WorkerPool) processNotification(n domain.Notification) {
	limiter, exists := wp.rateLimiters[n.Channel]
	if !exists {
		// Default limiter if channel unknown, or just log error
		// For now, assume known channels or use a default
		limiter = rate.NewLimiter(10, 1)
	}

	// 1. Rate Limit
	if err := limiter.Wait(context.Background()); err != nil {
		log.Printf("Rate limiter error for notification %d: %v", n.ID, err)
		// Should retry? Or just fail? Rate limiter Wait blocks until allowed, so error usually means ctx cancelled.
		return
	}

	// 2. Send
	err := wp.gateway.Send(n)

	// 3. Update Status
	status := domain.SENT
	if err != nil {
		log.Printf("Failed to send notification %d: %v", n.ID, err)
		status = domain.FAILED
	}

	if err := wp.repo.UpdateNotification(n.ID, status); err != nil {
		log.Printf("Failed to update status for notification %d: %v", n.ID, err)
	}
}

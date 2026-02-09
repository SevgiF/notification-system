package gateway

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/SevgiF/notification-system/internal/core/notification/domain"
)

type WebhookGateway struct {
	client *http.Client
	url    string
}

func NewWebhookGateway(url string) *WebhookGateway {
	return &WebhookGateway{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
		url: url,
	}
}

func (g *WebhookGateway) Send(n domain.Notification) error {
	payload := map[string]interface{}{
		"to":      n.Recipient,
		"channel": n.Channel,
		"content": n.Content,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequest("POST", g.url, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := g.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	}

	return fmt.Errorf("webhook replied with status: %d", resp.StatusCode)
}

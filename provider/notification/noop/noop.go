// Package noop provides a no-op implementation of the Notification provider.
package noop

import (
	"context"

	"github.com/gondolia/gondolia/provider"
	"github.com/gondolia/gondolia/provider/notification"
)

func init() {
	provider.Register[notification.NotificationProvider]("notification", "noop",
		provider.Metadata{
			Name:        "noop",
			DisplayName: "No-Op Notification Provider",
			Category:    "notification",
			Version:     "1.0.0",
			Description: "A no-operation notification provider for development and testing",
			ConfigSpec:  []provider.ConfigField{},
		},
		NewProvider,
	)
}

// Provider is a no-op Notification provider.
type Provider struct{}

// NewProvider creates a new no-op Notification provider.
func NewProvider(config map[string]any) (notification.NotificationProvider, error) {
	return &Provider{}, nil
}

func (p *Provider) Send(ctx context.Context, msg notification.Message) (*notification.SendResult, error) {
	return &notification.SendResult{
		MessageID: "noop-msg-001",
		Status:    "sent",
		Error:     "",
	}, nil
}

func (p *Provider) SendBatch(ctx context.Context, msgs []notification.Message) ([]notification.SendResult, error) {
	results := make([]notification.SendResult, len(msgs))
	for i := range msgs {
		results[i] = notification.SendResult{
			MessageID: "noop-msg-001",
			Status:    "sent",
			Error:     "",
		}
	}
	return results, nil
}

func (p *Provider) Channels() []string {
	return []string{"email", "sms", "push"}
}

func (p *Provider) Metadata() notification.Metadata {
	return notification.Metadata{
		Name:     "noop",
		Channels: []string{"email", "sms", "push"},
	}
}

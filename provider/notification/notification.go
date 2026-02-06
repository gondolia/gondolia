// Package notification defines the interface for notification channel integrations.
package notification

import "context"

// NotificationProvider abstracts notification channels (Email, SMS, Push, etc.).
type NotificationProvider interface {
	// Send sends a notification.
	Send(ctx context.Context, msg Message) (*SendResult, error)

	// SendBatch sends multiple notifications.
	SendBatch(ctx context.Context, msgs []Message) ([]SendResult, error)

	// Channels returns the supported channels.
	Channels() []string

	// Metadata returns provider information.
	Metadata() Metadata
}

// Message represents a notification message.
type Message struct {
	Channel      string // "email", "sms", "push"
	To           []Recipient
	From         string
	Subject      string
	Body         string // HTML for Email, Plain-Text for SMS/Push
	BodyPlain    string // Plain-Text fallback for Email
	Template     string // Template ID (optional)
	TemplateData map[string]any
	Metadata     map[string]string
	Attachments  []Attachment
}

// Recipient represents a message recipient.
type Recipient struct {
	Address string // Email, phone number, or device token
	Name    string
}

// Attachment represents an email attachment.
type Attachment struct {
	Filename    string
	ContentType string
	Data        []byte
}

// SendResult is the result of sending a notification.
type SendResult struct {
	MessageID string
	Status    string // "sent", "queued", "failed"
	Error     string
}

// Metadata provides information about the notification provider.
type Metadata struct {
	Name     string
	Channels []string
}

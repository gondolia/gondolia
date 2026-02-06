// Package payment defines the interface for payment gateway integrations.
package payment

import "context"

// PaymentProvider abstracts payment gateways (Saferpay, Stripe, Adyen, etc.).
type PaymentProvider interface {
	// Initialize starts a payment session.
	Initialize(ctx context.Context, req InitializeRequest) (*PaymentSession, error)

	// Authorize checks and authorizes a payment.
	Authorize(ctx context.Context, sessionID string) (*AuthorizationResult, error)

	// Capture captures an authorized payment.
	Capture(ctx context.Context, transactionID string, amount *Amount) (*CaptureResult, error)

	// Cancel cancels an authorized payment.
	Cancel(ctx context.Context, transactionID string) error

	// Refund refunds a captured payment (partial or full).
	Refund(ctx context.Context, transactionID string, amount Amount) (*RefundResult, error)

	// HandleWebhook processes provider-specific webhooks.
	HandleWebhook(ctx context.Context, payload []byte, headers map[string]string) (*WebhookEvent, error)

	// Metadata returns provider information.
	Metadata() Metadata
}

// InitializeRequest contains data for initializing a payment session.
type InitializeRequest struct {
	OrderID        string
	Amount         Amount
	Currency       string
	Description    string
	ReturnURL      string
	WebhookURL     string
	PaymentMethods []string // e.g. ["VISA", "MASTERCARD", "TWINT"]
	CustomerEmail  string
	Metadata       map[string]string
}

// PaymentSession represents an initialized payment session.
type PaymentSession struct {
	SessionID   string
	RedirectURL string
	Token       string
	ExpiresAt   string
}

// AuthorizationResult is the result of authorizing a payment.
type AuthorizationResult struct {
	TransactionID  string
	Status         string // "authorized", "failed", "pending"
	Amount         Amount
	PaymentMethod  string
	CardDisplay    string // "xxxx xxxx xxxx 1234"
	LiabilityShift bool
	Raw            map[string]any // Provider-specific raw data
}

// CaptureResult is the result of capturing a payment.
type CaptureResult struct {
	TransactionID string
	Status        string
	CapturedAt    string
}

// RefundResult is the result of refunding a payment.
type RefundResult struct {
	RefundID string
	Status   string
	Amount   Amount
}

// Amount represents a monetary amount.
type Amount struct {
	Value    int64  // Cents
	Currency string // ISO 4217
}

// WebhookEvent represents a webhook event from the payment provider.
type WebhookEvent struct {
	Type          string // "payment.authorized", "payment.captured", "payment.failed"
	TransactionID string
	OrderID       string
	Data          map[string]any
}

// Metadata provides information about the payment provider.
type Metadata struct {
	Name             string
	SupportedMethods []string
	TestMode         bool
}

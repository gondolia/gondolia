// Package noop provides a no-op implementation of the Payment provider.
package noop

import (
	"context"

	"github.com/gondolia/gondolia/provider"
	"github.com/gondolia/gondolia/provider/payment"
)

func init() {
	provider.Register[payment.PaymentProvider]("payment", "noop",
		provider.Metadata{
			Name:        "noop",
			DisplayName: "No-Op Payment Provider",
			Category:    "payment",
			Version:     "1.0.0",
			Description: "A no-operation payment provider for development and testing",
			ConfigSpec:  []provider.ConfigField{},
		},
		NewProvider,
	)
}

// Provider is a no-op Payment provider.
type Provider struct{}

// NewProvider creates a new no-op Payment provider.
func NewProvider(config map[string]any) (payment.PaymentProvider, error) {
	return &Provider{}, nil
}

func (p *Provider) Initialize(ctx context.Context, req payment.InitializeRequest) (*payment.PaymentSession, error) {
	return &payment.PaymentSession{
		SessionID:   "noop-session-001",
		RedirectURL: req.ReturnURL,
		Token:       "noop-token",
		ExpiresAt:   "",
	}, nil
}

func (p *Provider) Authorize(ctx context.Context, sessionID string) (*payment.AuthorizationResult, error) {
	return &payment.AuthorizationResult{
		TransactionID:  "noop-transaction-001",
		Status:         "authorized",
		Amount:         payment.Amount{Value: 0, Currency: "USD"},
		PaymentMethod:  "noop",
		CardDisplay:    "",
		LiabilityShift: false,
		Raw:            make(map[string]any),
	}, nil
}

func (p *Provider) Capture(ctx context.Context, transactionID string, amount *payment.Amount) (*payment.CaptureResult, error) {
	return &payment.CaptureResult{
		TransactionID: transactionID,
		Status:        "captured",
		CapturedAt:    "",
	}, nil
}

func (p *Provider) Cancel(ctx context.Context, transactionID string) error {
	return nil
}

func (p *Provider) Refund(ctx context.Context, transactionID string, amount payment.Amount) (*payment.RefundResult, error) {
	return &payment.RefundResult{
		RefundID: "noop-refund-001",
		Status:   "refunded",
		Amount:   amount,
	}, nil
}

func (p *Provider) HandleWebhook(ctx context.Context, payload []byte, headers map[string]string) (*payment.WebhookEvent, error) {
	return &payment.WebhookEvent{
		Type:          "payment.authorized",
		TransactionID: "noop-transaction-001",
		OrderID:       "",
		Data:          make(map[string]any),
	}, nil
}

func (p *Provider) Metadata() payment.Metadata {
	return payment.Metadata{
		Name:             "noop",
		SupportedMethods: []string{},
		TestMode:         true,
	}
}

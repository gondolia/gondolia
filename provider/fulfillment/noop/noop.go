// Package noop provides a no-op implementation of the Fulfillment provider.
package noop

import (
	"context"
	"time"

	"github.com/gondolia/gondolia/provider"
	"github.com/gondolia/gondolia/provider/fulfillment"
)

func init() {
	provider.Register[fulfillment.FulfillmentProvider]("fulfillment", "noop",
		provider.Metadata{
			Name:        "noop",
			DisplayName: "No-Op Fulfillment Provider",
			Category:    "fulfillment",
			Version:     "1.0.0",
			Description: "A no-operation fulfillment provider for development and testing",
			ConfigSpec:  []provider.ConfigField{},
		},
		NewProvider,
	)
}

// Provider is a no-op Fulfillment provider.
type Provider struct{}

// NewProvider creates a new no-op Fulfillment provider.
func NewProvider(config map[string]any) (fulfillment.FulfillmentProvider, error) {
	return &Provider{}, nil
}

func (p *Provider) CreateShipment(ctx context.Context, req fulfillment.ShipmentRequest) (*fulfillment.ShipmentResult, error) {
	estimatedDate := time.Now().Add(7 * 24 * time.Hour)
	return &fulfillment.ShipmentResult{
		ShipmentID:     "noop-shipment-001",
		TrackingNumber: "NOOP123456789",
		Label:          []byte{},
		EstimatedDate:  &estimatedDate,
	}, nil
}

func (p *Provider) GetShipmentStatus(ctx context.Context, shipmentID string) (*fulfillment.ShipmentStatus, error) {
	return &fulfillment.ShipmentStatus{
		ShipmentID:     shipmentID,
		Status:         "created",
		TrackingNumber: "NOOP123456789",
		Events:         []fulfillment.TrackingEvent{},
	}, nil
}

func (p *Provider) CancelShipment(ctx context.Context, shipmentID string) error {
	return nil
}

func (p *Provider) GetTrackingURL(ctx context.Context, trackingNumber string) (string, error) {
	return "https://example.com/track/" + trackingNumber, nil
}

func (p *Provider) CalculateShipping(ctx context.Context, req fulfillment.ShippingCalcRequest) ([]fulfillment.ShippingOption, error) {
	return []fulfillment.ShippingOption{
		{
			Service:       "standard",
			Name:          "Standard Shipping",
			Price:         0,
			Currency:      "USD",
			EstimatedDays: 7,
		},
	}, nil
}

func (p *Provider) Metadata() fulfillment.Metadata {
	return fulfillment.Metadata{
		Name:     "noop",
		Carriers: []string{},
	}
}

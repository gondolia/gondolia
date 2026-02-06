// Package fulfillment defines the interface for logistics/shipping system integrations.
package fulfillment

import (
	"context"
	"time"
)

// FulfillmentProvider abstracts logistics/shipping systems (DHL, UPS, Swiss Post, etc.).
type FulfillmentProvider interface {
	// CreateShipment creates a shipment order.
	CreateShipment(ctx context.Context, req ShipmentRequest) (*ShipmentResult, error)

	// GetShipmentStatus retrieves shipment status.
	GetShipmentStatus(ctx context.Context, shipmentID string) (*ShipmentStatus, error)

	// CancelShipment cancels a shipment order.
	CancelShipment(ctx context.Context, shipmentID string) error

	// GetTrackingURL returns the tracking URL.
	GetTrackingURL(ctx context.Context, trackingNumber string) (string, error)

	// CalculateShipping calculates shipping costs.
	CalculateShipping(ctx context.Context, req ShippingCalcRequest) ([]ShippingOption, error)

	// Metadata returns provider information.
	Metadata() Metadata
}

// ShipmentRequest contains data for creating a shipment.
type ShipmentRequest struct {
	OrderID  string
	From     Address
	To       Address
	Packages []Package
	Service  string // e.g. "standard", "express"
	Notes    string
}

// ShipmentResult is the result of creating a shipment.
type ShipmentResult struct {
	ShipmentID     string
	TrackingNumber string
	Label          []byte // PDF label
	EstimatedDate  *time.Time
}

// ShipmentStatus represents the status of a shipment.
type ShipmentStatus struct {
	ShipmentID     string
	Status         string // "created", "picked_up", "in_transit", "delivered", "failed"
	TrackingNumber string
	Events         []TrackingEvent
}

// TrackingEvent represents a tracking event.
type TrackingEvent struct {
	Timestamp   time.Time
	Status      string
	Location    string
	Description string
}

// ShippingCalcRequest contains data for calculating shipping costs.
type ShippingCalcRequest struct {
	From     Address
	To       Address
	Packages []Package
}

// ShippingOption represents a shipping option with cost.
type ShippingOption struct {
	Service       string
	Name          string
	Price         int64
	Currency      string
	EstimatedDays int
}

// Address represents a shipping address.
type Address struct {
	Name       string
	Street     string
	PostalCode string
	City       string
	Country    string // ISO 3166-1 alpha-2
}

// Package represents a package to be shipped.
type Package struct {
	WeightGrams int
	LengthCm    int
	WidthCm     int
	HeightCm    int
}

// Metadata provides information about the fulfillment provider.
type Metadata struct {
	Name     string
	Carriers []string
}

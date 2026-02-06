// Package erp defines the interface for ERP system integrations (SAP, Microsoft Dynamics, etc.).
package erp

import (
	"context"
	"time"
)

// ERPProvider abstracts the communication with an ERP system.
// Implementations may include SAP R/3, SAP S/4HANA, Microsoft Dynamics, etc.
type ERPProvider interface {
	// --- Order Management ---

	// CreateOrder transmits an order to the ERP system.
	CreateOrder(ctx context.Context, req CreateOrderRequest) (*CreateOrderResult, error)

	// SimulateOrder calculates pricing and availability without committing.
	SimulateOrder(ctx context.Context, req SimulateOrderRequest) (*SimulateOrderResult, error)

	// GetOrderStatus retrieves the current status of an order.
	GetOrderStatus(ctx context.Context, orderID string) (*OrderStatus, error)

	// --- Inventory ---

	// GetProductAvailability retrieves stock levels for products.
	GetProductAvailability(ctx context.Context, skus []string) ([]ProductStock, error)

	// --- Pricing ---

	// GetTierPrices retrieves tiered pricing for customer/product combinations.
	GetTierPrices(ctx context.Context, req TierPriceRequest) ([]TierPrice, error)

	// --- Company/Customer Data ---

	// SyncCompany synchronizes company data from the ERP system.
	SyncCompany(ctx context.Context, erpCustomerID string) (*CompanyData, error)

	// GetCompanyAddresses retrieves all addresses for a company.
	GetCompanyAddresses(ctx context.Context, erpCustomerID string) ([]Address, error)

	// --- Reports ---

	// GetOrderHistory retrieves historical orders.
	GetOrderHistory(ctx context.Context, req ReportFilter) ([]OrderReport, error)

	// GetShipmentHistory retrieves historical shipments.
	GetShipmentHistory(ctx context.Context, req ReportFilter) ([]ShipmentReport, error)

	// GetInvoiceHistory retrieves historical invoices.
	GetInvoiceHistory(ctx context.Context, req ReportFilter) ([]InvoiceReport, error)

	// --- Metadata ---

	// Metadata returns information about this provider.
	Metadata() Metadata
}

// --- Request/Response Types ---

// CreateOrderRequest contains all data needed to create an order in the ERP.
type CreateOrderRequest struct {
	TenantConfig TenantConfig
	Order        Order
	Customer     Customer
	ShipTo       Address
	BillTo       Address
}

// CreateOrderResult is the response from creating an order.
type CreateOrderResult struct {
	ERPOrderNumber string
	Items          []OrderItemResult
	Messages       []Message
}

// SimulateOrderRequest contains data for order simulation (pricing/availability check).
type SimulateOrderRequest struct {
	TenantConfig TenantConfig
	Items        []SimulateItem
	Customer     Customer
	ShipTo       Address
	DesiredDate  *time.Time
}

// SimulateOrderResult is the response from order simulation.
type SimulateOrderResult struct {
	Items    []SimulatedItem
	Totals   Totals
	Schedule []DeliverySchedule
	Messages []Message
}

// OrderStatus represents the current status of an order in the ERP.
type OrderStatus struct {
	ERPOrderNumber string
	Status         string
	Items          []OrderItemStatus
}

// ProductStock represents inventory information for a product.
type ProductStock struct {
	SKU          string
	PlantCode    string
	PlantName    string
	Quantity     float64
	Unit         string
	LeadTimeDays int
}

// TierPriceRequest contains parameters for retrieving tiered pricing.
type TierPriceRequest struct {
	CustomerID string
	SKUs       []string
	Currency   string
}

// TierPrice represents a quantity-based price tier.
type TierPrice struct {
	SKU       string
	MinQty    float64
	Price     float64
	Currency  string
	ValidFrom time.Time
	ValidTo   time.Time
}

// CompanyData represents company master data from the ERP.
type CompanyData struct {
	ERPCustomerID string
	Name          string
	TaxID         string
	PaymentTerms  string
	Currency      string
	CreditLimit   float64
	Attributes    map[string]string
}

// Address represents a business address.
type Address struct {
	ID         string
	Name       string
	Street     string
	PostalCode string
	City       string
	Country    string // ISO 3166-1 alpha-2
	Region     string
}

// ReportFilter contains parameters for historical data queries.
type ReportFilter struct {
	CustomerID string
	DateFrom   time.Time
	DateTo     time.Time
	Limit      int
	Offset     int
}

// OrderReport represents a historical order.
type OrderReport struct {
	ERPOrderNumber string
	OrderDate      time.Time
	CustomerPO     string
	TotalAmount    float64
	Currency       string
	Status         string
}

// ShipmentReport represents a historical shipment.
type ShipmentReport struct {
	ShipmentID     string
	OrderNumber    string
	ShipDate       time.Time
	TrackingNumber string
	Carrier        string
}

// InvoiceReport represents a historical invoice.
type InvoiceReport struct {
	InvoiceNumber string
	OrderNumber   string
	InvoiceDate   time.Time
	DueDate       time.Time
	Amount        float64
	Currency      string
	Status        string
}

// TenantConfig contains tenant-specific ERP configuration.
type TenantConfig struct {
	SalesOrg    string // SAP: VKORG
	DistChannel string // SAP: VTWEG
	Division    string // SAP: SPART
	Currency    string
	Language    string
}

// Order represents an order to be created.
type Order struct {
	ExternalID string
	CustomerPO string
	Items      []OrderItem
	Notes      string
}

// OrderItem represents a line item in an order.
type OrderItem struct {
	SKU            string
	Quantity       float64
	Unit           string
	RequestedPrice *float64 // Optional: for contract prices
	Plant          string   // Target plant/warehouse
}

// Customer represents customer data for order creation.
type Customer struct {
	ERPCustomerID string
	SoldToParty   string // SAP: KUNNR (Sold-To)
	ShipToParty   string // SAP: Ship-To
	BillToParty   string // SAP: Bill-To
	PayerParty    string // SAP: Payer
}

// SimulateItem represents an item for order simulation.
type SimulateItem struct {
	SKU      string
	Quantity float64
	Unit     string
}

// SimulatedItem is the result of simulating an order item.
type SimulatedItem struct {
	SKU         string
	Quantity    float64
	Unit        string
	UnitPrice   float64
	TotalPrice  float64
	Available   bool
	LeadTimeDays int
}

// OrderItemResult represents the result of creating an order item.
type OrderItemResult struct {
	SKU            string
	ItemNumber     string
	ConfirmedQty   float64
	ConfirmedPrice float64
}

// OrderItemStatus represents the status of an order item.
type OrderItemStatus struct {
	ItemNumber string
	SKU        string
	Status     string
	ShippedQty float64
}

// Totals represents order totals.
type Totals struct {
	Subtotal    float64
	Tax         float64
	Shipping    float64
	Total       float64
	Currency    string
}

// DeliverySchedule represents a scheduled delivery date.
type DeliverySchedule struct {
	Date     time.Time
	Quantity float64
	Items    []string // SKUs
}

// Message represents a message from the ERP system.
type Message struct {
	Type    string // "success", "warning", "error", "info"
	Code    string
	Message string
}

// Metadata provides information about the ERP provider implementation.
type Metadata struct {
	Name         string
	Version      string
	Protocol     string // "soap", "odata", "rest"
	Capabilities []string
}

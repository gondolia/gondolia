// Package noop provides a no-op implementation of the ERP provider.
// This is useful for development and testing when no ERP system is configured.
package noop

import (
	"context"

	"github.com/gondolia/gondolia/provider"
	"github.com/gondolia/gondolia/provider/erp"
)

func init() {
	provider.Register[erp.ERPProvider]("erp", "noop",
		provider.Metadata{
			Name:        "noop",
			DisplayName: "No-Op ERP Provider",
			Category:    "erp",
			Version:     "1.0.0",
			Description: "A no-operation ERP provider for development and testing",
			ConfigSpec:  []provider.ConfigField{},
		},
		NewProvider,
	)
}

// Provider is a no-op ERP provider that returns empty/default values.
type Provider struct{}

// NewProvider creates a new no-op ERP provider.
func NewProvider(config map[string]any) (erp.ERPProvider, error) {
	return &Provider{}, nil
}

func (p *Provider) CreateOrder(ctx context.Context, req erp.CreateOrderRequest) (*erp.CreateOrderResult, error) {
	return &erp.CreateOrderResult{
		ERPOrderNumber: "NOOP-00000001",
		Items:          []erp.OrderItemResult{},
		Messages:       []erp.Message{},
	}, nil
}

func (p *Provider) SimulateOrder(ctx context.Context, req erp.SimulateOrderRequest) (*erp.SimulateOrderResult, error) {
	items := make([]erp.SimulatedItem, len(req.Items))
	for i, item := range req.Items {
		items[i] = erp.SimulatedItem{
			SKU:          item.SKU,
			Quantity:     item.Quantity,
			Unit:         item.Unit,
			UnitPrice:    0,
			TotalPrice:   0,
			Available:    true,
			LeadTimeDays: 0,
		}
	}

	return &erp.SimulateOrderResult{
		Items: items,
		Totals: erp.Totals{
			Subtotal: 0,
			Tax:      0,
			Shipping: 0,
			Total:    0,
			Currency: req.TenantConfig.Currency,
		},
		Schedule: []erp.DeliverySchedule{},
		Messages: []erp.Message{},
	}, nil
}

func (p *Provider) GetOrderStatus(ctx context.Context, orderID string) (*erp.OrderStatus, error) {
	return &erp.OrderStatus{
		ERPOrderNumber: orderID,
		Status:         "unknown",
		Items:          []erp.OrderItemStatus{},
	}, nil
}

func (p *Provider) GetProductAvailability(ctx context.Context, skus []string) ([]erp.ProductStock, error) {
	stocks := make([]erp.ProductStock, len(skus))
	for i, sku := range skus {
		stocks[i] = erp.ProductStock{
			SKU:          sku,
			PlantCode:    "0000",
			PlantName:    "Default Plant",
			Quantity:     999,
			Unit:         "PC",
			LeadTimeDays: 0,
		}
	}
	return stocks, nil
}

func (p *Provider) GetTierPrices(ctx context.Context, req erp.TierPriceRequest) ([]erp.TierPrice, error) {
	return []erp.TierPrice{}, nil
}

func (p *Provider) SyncCompany(ctx context.Context, erpCustomerID string) (*erp.CompanyData, error) {
	return &erp.CompanyData{
		ERPCustomerID: erpCustomerID,
		Name:          "Demo Company",
		TaxID:         "",
		PaymentTerms:  "NET30",
		Currency:      "USD",
		CreditLimit:   0,
		Attributes:    make(map[string]string),
	}, nil
}

func (p *Provider) GetCompanyAddresses(ctx context.Context, erpCustomerID string) ([]erp.Address, error) {
	return []erp.Address{}, nil
}

func (p *Provider) GetOrderHistory(ctx context.Context, req erp.ReportFilter) ([]erp.OrderReport, error) {
	return []erp.OrderReport{}, nil
}

func (p *Provider) GetShipmentHistory(ctx context.Context, req erp.ReportFilter) ([]erp.ShipmentReport, error) {
	return []erp.ShipmentReport{}, nil
}

func (p *Provider) GetInvoiceHistory(ctx context.Context, req erp.ReportFilter) ([]erp.InvoiceReport, error) {
	return []erp.InvoiceReport{}, nil
}

func (p *Provider) Metadata() erp.Metadata {
	return erp.Metadata{
		Name:         "noop",
		Version:      "1.0.0",
		Protocol:     "none",
		Capabilities: []string{},
	}
}

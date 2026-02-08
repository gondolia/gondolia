package domain

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

// Company represents a company/business entity
type Company struct {
	ID       uuid.UUID `json:"id"`
	TenantID uuid.UUID `json:"tenant_id"`

	// SAP Mapping
	SAPCompanyNumber string  `json:"sap_company_number"`
	SAPCustomerGroup *string `json:"sap_customer_group,omitempty"`
	SAPShippingPlant *string `json:"sap_shipping_plant,omitempty"`
	SAPOffice        *string `json:"sap_office,omitempty"`
	SAPPaymentType   *string `json:"sap_payment_type,omitempty"`
	SAPPriceGroup    *string `json:"sap_price_group,omitempty"`

	// Profile
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
	Email       *string `json:"email,omitempty"`
	Currency    string  `json:"currency"`

	// Address
	Street      *string `json:"street,omitempty"`
	HouseNumber *string `json:"house_number,omitempty"`
	ZIP         *string `json:"zip,omitempty"`
	City        *string `json:"city,omitempty"`
	Country     string  `json:"country"`

	// Contact
	Phone *string `json:"phone,omitempty"`
	Fax   *string `json:"fax,omitempty"`
	URL   *string `json:"url,omitempty"`

	// Config
	Config              map[string]any `json:"config,omitempty"`
	DesiredDeliveryDays []string       `json:"desired_delivery_days,omitempty"`
	DefaultShippingNote *string        `json:"default_shipping_note,omitempty"`
	DisableOrderFeature bool           `json:"disable_order_feature"`

	// Branding
	CustomPrimaryColor   *string `json:"custom_primary_color,omitempty"`
	CustomSecondaryColor *string `json:"custom_secondary_color,omitempty"`

	// Status
	IsActive  bool       `json:"is_active"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

// NewCompany creates a new company with defaults
func NewCompany(tenantID uuid.UUID, sapCompanyNumber, name string) *Company {
	now := time.Now()
	return &Company{
		ID:               uuid.New(),
		TenantID:         tenantID,
		SAPCompanyNumber: strings.TrimSpace(sapCompanyNumber),
		Name:             strings.TrimSpace(name),
		Currency:         "CHF",
		Country:          "CH",
		Config:           make(map[string]any),
		IsActive:         false,
		CreatedAt:        now,
		UpdatedAt:        now,
	}
}

// FullAddress returns formatted address string
func (c *Company) FullAddress() string {
	var parts []string

	if c.Street != nil {
		addr := *c.Street
		if c.HouseNumber != nil {
			addr += " " + *c.HouseNumber
		}
		parts = append(parts, addr)
	}

	if c.ZIP != nil || c.City != nil {
		var line string
		if c.ZIP != nil {
			line = *c.ZIP
		}
		if c.City != nil {
			if line != "" {
				line += " "
			}
			line += *c.City
		}
		parts = append(parts, line)
	}

	parts = append(parts, c.Country)
	return strings.Join(parts, ", ")
}

// CreateCompanyRequest represents a request to create a company
type CreateCompanyRequest struct {
	SAPCompanyNumber string  `json:"sap_company_number" binding:"required"`
	Name             string  `json:"name" binding:"required,min=1,max=255"`
	Email            *string `json:"email,omitempty" binding:"omitempty,email"`
	Currency         *string `json:"currency,omitempty"`
	Country          *string `json:"country,omitempty"`

	// Address
	Street      *string `json:"street,omitempty"`
	HouseNumber *string `json:"house_number,omitempty"`
	ZIP         *string `json:"zip,omitempty"`
	City        *string `json:"city,omitempty"`

	// SAP
	SAPCustomerGroup *string `json:"sap_customer_group,omitempty"`
	SAPShippingPlant *string `json:"sap_shipping_plant,omitempty"`
	SAPOffice        *string `json:"sap_office,omitempty"`
}

// UpdateCompanyRequest represents a request to update a company
type UpdateCompanyRequest struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	Email       *string `json:"email,omitempty"`

	// Address
	Street      *string `json:"street,omitempty"`
	HouseNumber *string `json:"house_number,omitempty"`
	ZIP         *string `json:"zip,omitempty"`
	City        *string `json:"city,omitempty"`
	Country     *string `json:"country,omitempty"`

	// Contact
	Phone *string `json:"phone,omitempty"`
	Fax   *string `json:"fax,omitempty"`
	URL   *string `json:"url,omitempty"`

	// Config
	DesiredDeliveryDays []string `json:"desired_delivery_days,omitempty"`
	DefaultShippingNote *string  `json:"default_shipping_note,omitempty"`
	DisableOrderFeature *bool    `json:"disable_order_feature,omitempty"`

	// Branding
	CustomPrimaryColor   *string `json:"custom_primary_color,omitempty"`
	CustomSecondaryColor *string `json:"custom_secondary_color,omitempty"`

	IsActive *bool `json:"is_active,omitempty"`
}

// CompanyFilter represents filter options for listing companies
type CompanyFilter struct {
	TenantID         uuid.UUID
	SAPCompanyNumber *string
	IsActive         *bool
	Search           *string
	Limit            int
	Offset           int
}

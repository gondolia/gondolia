// Package crm defines the interface for CRM system integrations.
package crm

import "context"

// CRMProvider abstracts CRM systems (MS Dynamics, Salesforce, etc.).
type CRMProvider interface {
	// SyncContact synchronizes a contact with the CRM.
	SyncContact(ctx context.Context, contact Contact) (*SyncResult, error)

	// GetAccount retrieves company data from the CRM.
	GetAccount(ctx context.Context, accountID string) (*Account, error)

	// ListAccounts retrieves a list of companies.
	ListAccounts(ctx context.Context, filter AccountFilter) ([]Account, error)

	// Metadata returns provider information.
	Metadata() Metadata
}

// Contact represents a CRM contact.
type Contact struct {
	ExternalID string
	Email      string
	FirstName  string
	LastName   string
	Phone      string
	Company    string
	Attributes map[string]any
}

// Account represents a CRM account (company).
type Account struct {
	ExternalID    string
	Name          string
	ERPCustomerID string
	Attributes    map[string]any
}

// AccountFilter contains filter criteria for listing accounts.
type AccountFilter struct {
	Query  string
	Limit  int
	Offset int
}

// SyncResult is the result of syncing a contact.
type SyncResult struct {
	ExternalID string
	Action     string // "created", "updated", "unchanged"
}

// Metadata provides information about the CRM provider.
type Metadata struct {
	Name string
}

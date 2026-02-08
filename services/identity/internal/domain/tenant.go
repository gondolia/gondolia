package domain

import (
	"time"

	"github.com/google/uuid"
)

// Tenant represents a tenant in the multi-tenant system
type Tenant struct {
	ID        uuid.UUID      `json:"id"`
	Code      string         `json:"code"`
	Name      string         `json:"name"`
	Config    map[string]any `json:"config,omitempty"`
	IsActive  bool           `json:"is_active"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

// NewTenant creates a new tenant with defaults
func NewTenant(code, name string) *Tenant {
	now := time.Now()
	return &Tenant{
		ID:        uuid.New(),
		Code:      code,
		Name:      name,
		Config:    make(map[string]any),
		IsActive:  true,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

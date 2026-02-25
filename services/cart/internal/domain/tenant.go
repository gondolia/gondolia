package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// Tenant represents a tenant in the system
type Tenant struct {
	ID        uuid.UUID `json:"id"`
	Code      string    `json:"code"`
	Name      string    `json:"name"`
	Config    []byte    `json:"config,omitempty"` // JSONB
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

var (
	ErrTenantNotFound = errors.New("tenant not found")
)

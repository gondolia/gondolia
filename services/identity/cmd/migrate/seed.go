package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

// seedDatabase creates test data for development
func seedDatabase(ctx context.Context, databaseURL string) error {
	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		return fmt.Errorf("connecting to database: %w", err)
	}
	defer pool.Close()

	// Check if already seeded
	var count int
	err = pool.QueryRow(ctx, "SELECT COUNT(*) FROM tenants").Scan(&count)
	if err != nil {
		return fmt.Errorf("checking existing data: %w", err)
	}
	if count > 0 {
		fmt.Println("Database already seeded, skipping...")
		return nil
	}

	fmt.Println("Seeding database...")

	// Create test tenant
	tenantID := uuid.MustParse("00000000-0000-0000-0000-000000000001")
	_, err = pool.Exec(ctx, `
		INSERT INTO tenants (id, code, name, is_active, config, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, tenantID, "demo", "Demo Tenant", true, "{}", time.Now(), time.Now())
	if err != nil {
		return fmt.Errorf("creating tenant: %w", err)
	}
	fmt.Println("  Created tenant: demo")

	// Create test company
	companyID := uuid.MustParse("00000000-0000-0000-0000-000000000002")
	_, err = pool.Exec(ctx, `
		INSERT INTO companies (
			id, tenant_id, sap_company_number, name, is_active, currency, country,
			created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`, companyID, tenantID, "1000", "Demo Company GmbH", true, "EUR", "DE", time.Now(), time.Now())
	if err != nil {
		return fmt.Errorf("creating company: %w", err)
	}
	fmt.Println("  Created company: Demo Company GmbH (SAP: 1000)")

	// Create admin role
	adminRoleID := uuid.MustParse("00000000-0000-0000-0000-000000000003")
	adminPermissions := `{
		"company.manage": true,
		"company.manage-users-and-roles-and-permissions": true,
		"company.manage-settings": true,
		"company.manage-addresses": true,
		"company.manage-custom-skus": true,
		"company.manage-watchlists": true,
		"company.order-data.see-orders": true,
		"company.order-data.see-invoices": true,
		"company.order-data.see-shipments": true,
		"company.order-data.see-reshipments": true,
		"company.order-data.see-credits": true,
		"sales.create-order": true
	}`
	_, err = pool.Exec(ctx, `
		INSERT INTO roles (id, tenant_id, company_id, name, permissions, is_system, created_at, updated_at)
		VALUES ($1, $2, NULL, $3, $4, $5, $6, $7)
	`, adminRoleID, tenantID, "Administrator", adminPermissions, true, time.Now(), time.Now())
	if err != nil {
		return fmt.Errorf("creating admin role: %w", err)
	}
	fmt.Println("  Created role: Administrator (system)")

	// Create user role
	userRoleID := uuid.MustParse("00000000-0000-0000-0000-000000000004")
	userPermissions := `{
		"company.manage-watchlists": true,
		"company.order-data.see-orders": true,
		"company.order-data.see-invoices": true,
		"company.order-data.see-shipments": true,
		"sales.create-order": true
	}`
	_, err = pool.Exec(ctx, `
		INSERT INTO roles (id, tenant_id, company_id, name, permissions, is_system, created_at, updated_at)
		VALUES ($1, $2, NULL, $3, $4, $5, $6, $7)
	`, userRoleID, tenantID, "Benutzer", userPermissions, true, time.Now(), time.Now())
	if err != nil {
		return fmt.Errorf("creating user role: %w", err)
	}
	fmt.Println("  Created role: Benutzer (system)")

	// Create admin user
	adminUserID := uuid.MustParse("00000000-0000-0000-0000-000000000005")
	passwordHash, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	_, err = pool.Exec(ctx, `
		INSERT INTO users (
			id, tenant_id, is_active, is_imported, is_salesmaster, sso_only,
			email, password_hash, firstname, lastname, default_language,
			default_company_id, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
	`, adminUserID, tenantID, true, false, true, false,
		"admin@demo.local", string(passwordHash), "Admin", "User", "de",
		companyID, time.Now(), time.Now())
	if err != nil {
		return fmt.Errorf("creating admin user: %w", err)
	}
	fmt.Println("  Created user: admin@demo.local (password: admin123)")

	// Assign admin to company with admin role
	_, err = pool.Exec(ctx, `
		INSERT INTO user_companies (user_id, company_id, role_id, user_type, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, adminUserID, companyID, adminRoleID, 0, time.Now(), time.Now())
	if err != nil {
		return fmt.Errorf("assigning admin to company: %w", err)
	}

	// Create test user
	testUserID := uuid.MustParse("00000000-0000-0000-0000-000000000006")
	testPasswordHash, _ := bcrypt.GenerateFromPassword([]byte("test123"), bcrypt.DefaultCost)
	_, err = pool.Exec(ctx, `
		INSERT INTO users (
			id, tenant_id, is_active, is_imported, is_salesmaster, sso_only,
			email, password_hash, firstname, lastname, default_language,
			default_company_id, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
	`, testUserID, tenantID, true, false, false, false,
		"user@demo.local", string(testPasswordHash), "Test", "User", "de",
		companyID, time.Now(), time.Now())
	if err != nil {
		return fmt.Errorf("creating test user: %w", err)
	}
	fmt.Println("  Created user: user@demo.local (password: test123)")

	// Assign test user to company with user role
	_, err = pool.Exec(ctx, `
		INSERT INTO user_companies (user_id, company_id, role_id, user_type, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, testUserID, companyID, userRoleID, 1, time.Now(), time.Now())
	if err != nil {
		return fmt.Errorf("assigning test user to company: %w", err)
	}

	fmt.Println("\nSeed completed! Test credentials:")
	fmt.Println("  Admin: admin@demo.local / admin123")
	fmt.Println("  User:  user@demo.local / test123")
	fmt.Println("  Tenant: demo")

	return nil
}

package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/lib/pq"

	"github.com/gondolia/gondolia/services/identity/internal/domain"
)

type CompanyRepository struct {
	db *DB
}

func NewCompanyRepository(db *DB) *CompanyRepository {
	return &CompanyRepository{db: db}
}

func (r *CompanyRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Company, error) {
	query := `
		SELECT id, tenant_id, sap_company_number, sap_customer_group, sap_shipping_plant,
		       sap_office, sap_payment_type, sap_price_group, name, description, email,
		       currency, street, house_number, zip, city, country, phone, fax, url,
		       config, desired_delivery_days, default_shipping_note, disable_order_feature,
		       custom_primary_color, custom_secondary_color, is_active,
		       created_at, updated_at, deleted_at
		FROM companies
		WHERE id = $1 AND deleted_at IS NULL
	`

	return r.scanCompany(r.db.Pool.QueryRow(ctx, query, id))
}

func (r *CompanyRepository) GetBySAPNumber(ctx context.Context, tenantID uuid.UUID, sapNumber string) (*domain.Company, error) {
	query := `
		SELECT id, tenant_id, sap_company_number, sap_customer_group, sap_shipping_plant,
		       sap_office, sap_payment_type, sap_price_group, name, description, email,
		       currency, street, house_number, zip, city, country, phone, fax, url,
		       config, desired_delivery_days, default_shipping_note, disable_order_feature,
		       custom_primary_color, custom_secondary_color, is_active,
		       created_at, updated_at, deleted_at
		FROM companies
		WHERE tenant_id = $1 AND sap_company_number = $2 AND deleted_at IS NULL
	`

	return r.scanCompany(r.db.Pool.QueryRow(ctx, query, tenantID, sapNumber))
}

func (r *CompanyRepository) List(ctx context.Context, filter domain.CompanyFilter) ([]domain.Company, int, error) {
	var conditions []string
	var args []any
	argNum := 1

	conditions = append(conditions, fmt.Sprintf("tenant_id = $%d", argNum))
	args = append(args, filter.TenantID)
	argNum++

	conditions = append(conditions, "deleted_at IS NULL")

	if filter.SAPCompanyNumber != nil {
		conditions = append(conditions, fmt.Sprintf("sap_company_number = $%d", argNum))
		args = append(args, *filter.SAPCompanyNumber)
		argNum++
	}

	if filter.IsActive != nil {
		conditions = append(conditions, fmt.Sprintf("is_active = $%d", argNum))
		args = append(args, *filter.IsActive)
		argNum++
	}

	if filter.Search != nil {
		conditions = append(conditions, fmt.Sprintf("LOWER(name) LIKE LOWER($%d)", argNum))
		args = append(args, "%"+*filter.Search+"%")
		argNum++
	}

	whereClause := strings.Join(conditions, " AND ")

	// Count query
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM companies WHERE %s", whereClause)
	var total int
	if err := r.db.Pool.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	// Data query
	query := fmt.Sprintf(`
		SELECT id, tenant_id, sap_company_number, sap_customer_group, sap_shipping_plant,
		       sap_office, sap_payment_type, sap_price_group, name, description, email,
		       currency, street, house_number, zip, city, country, phone, fax, url,
		       config, desired_delivery_days, default_shipping_note, disable_order_feature,
		       custom_primary_color, custom_secondary_color, is_active,
		       created_at, updated_at, deleted_at
		FROM companies
		WHERE %s
		ORDER BY name
		LIMIT $%d OFFSET $%d
	`, whereClause, argNum, argNum+1)

	args = append(args, filter.Limit, filter.Offset)

	rows, err := r.db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var companies []domain.Company
	for rows.Next() {
		company, err := r.scanCompanyFromRows(rows)
		if err != nil {
			return nil, 0, err
		}
		companies = append(companies, *company)
	}

	return companies, total, rows.Err()
}

func (r *CompanyRepository) Create(ctx context.Context, company *domain.Company) error {
	query := `
		INSERT INTO companies (
			id, tenant_id, sap_company_number, sap_customer_group, sap_shipping_plant,
			sap_office, sap_payment_type, sap_price_group, name, description, email,
			currency, street, house_number, zip, city, country, phone, fax, url,
			config, desired_delivery_days, default_shipping_note, disable_order_feature,
			custom_primary_color, custom_secondary_color, is_active, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17,
			$18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29
		)
	`

	configJSON, err := json.Marshal(company.Config)
	if err != nil {
		return err
	}

	_, err = r.db.Pool.Exec(ctx, query,
		company.ID, company.TenantID, company.SAPCompanyNumber, company.SAPCustomerGroup,
		company.SAPShippingPlant, company.SAPOffice, company.SAPPaymentType, company.SAPPriceGroup,
		company.Name, company.Description, company.Email, company.Currency,
		company.Street, company.HouseNumber, company.ZIP, company.City, company.Country,
		company.Phone, company.Fax, company.URL, configJSON,
		pq.Array(company.DesiredDeliveryDays), company.DefaultShippingNote, company.DisableOrderFeature,
		company.CustomPrimaryColor, company.CustomSecondaryColor, company.IsActive,
		company.CreatedAt, company.UpdatedAt,
	)
	return err
}

func (r *CompanyRepository) Update(ctx context.Context, company *domain.Company) error {
	query := `
		UPDATE companies SET
			sap_customer_group = $2, sap_shipping_plant = $3, sap_office = $4,
			sap_payment_type = $5, sap_price_group = $6, name = $7, description = $8,
			email = $9, currency = $10, street = $11, house_number = $12, zip = $13,
			city = $14, country = $15, phone = $16, fax = $17, url = $18, config = $19,
			desired_delivery_days = $20, default_shipping_note = $21, disable_order_feature = $22,
			custom_primary_color = $23, custom_secondary_color = $24, is_active = $25,
			updated_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
	`

	configJSON, err := json.Marshal(company.Config)
	if err != nil {
		return err
	}

	result, err := r.db.Pool.Exec(ctx, query,
		company.ID, company.SAPCustomerGroup, company.SAPShippingPlant, company.SAPOffice,
		company.SAPPaymentType, company.SAPPriceGroup, company.Name, company.Description,
		company.Email, company.Currency, company.Street, company.HouseNumber, company.ZIP,
		company.City, company.Country, company.Phone, company.Fax, company.URL, configJSON,
		pq.Array(company.DesiredDeliveryDays), company.DefaultShippingNote, company.DisableOrderFeature,
		company.CustomPrimaryColor, company.CustomSecondaryColor, company.IsActive,
	)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return domain.ErrCompanyNotFound
	}

	return nil
}

func (r *CompanyRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE companies SET deleted_at = NOW(), updated_at = NOW() WHERE id = $1 AND deleted_at IS NULL`

	result, err := r.db.Pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return domain.ErrCompanyNotFound
	}

	return nil
}

func (r *CompanyRepository) scanCompany(row pgx.Row) (*domain.Company, error) {
	var c domain.Company
	var configJSON []byte

	err := row.Scan(
		&c.ID, &c.TenantID, &c.SAPCompanyNumber, &c.SAPCustomerGroup, &c.SAPShippingPlant,
		&c.SAPOffice, &c.SAPPaymentType, &c.SAPPriceGroup, &c.Name, &c.Description, &c.Email,
		&c.Currency, &c.Street, &c.HouseNumber, &c.ZIP, &c.City, &c.Country,
		&c.Phone, &c.Fax, &c.URL, &configJSON, pq.Array(&c.DesiredDeliveryDays),
		&c.DefaultShippingNote, &c.DisableOrderFeature,
		&c.CustomPrimaryColor, &c.CustomSecondaryColor, &c.IsActive,
		&c.CreatedAt, &c.UpdatedAt, &c.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrCompanyNotFound
		}
		return nil, err
	}

	if configJSON != nil {
		if err := json.Unmarshal(configJSON, &c.Config); err != nil {
			return nil, err
		}
	}

	return &c, nil
}

func (r *CompanyRepository) scanCompanyFromRows(rows pgx.Rows) (*domain.Company, error) {
	var c domain.Company
	var configJSON []byte

	err := rows.Scan(
		&c.ID, &c.TenantID, &c.SAPCompanyNumber, &c.SAPCustomerGroup, &c.SAPShippingPlant,
		&c.SAPOffice, &c.SAPPaymentType, &c.SAPPriceGroup, &c.Name, &c.Description, &c.Email,
		&c.Currency, &c.Street, &c.HouseNumber, &c.ZIP, &c.City, &c.Country,
		&c.Phone, &c.Fax, &c.URL, &configJSON, pq.Array(&c.DesiredDeliveryDays),
		&c.DefaultShippingNote, &c.DisableOrderFeature,
		&c.CustomPrimaryColor, &c.CustomSecondaryColor, &c.IsActive,
		&c.CreatedAt, &c.UpdatedAt, &c.DeletedAt,
	)
	if err != nil {
		return nil, err
	}

	if configJSON != nil {
		if err := json.Unmarshal(configJSON, &c.Config); err != nil {
			return nil, err
		}
	}

	return &c, nil
}

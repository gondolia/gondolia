package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"

	"github.com/gondolia/gondolia/services/catalog/internal/domain"
)

// formatOptionLabel generates a human-readable i18n label from an option code
// e.g. "1_5kw" -> {"de": "1,5 kW", "en": "1.5 kW"}, "230v" -> {"de": "230V", "en": "230V"}
func formatOptionLabel(code string) map[string]string {
	label := code

	// Common unit suffixes
	units := map[string]string{
		"kw": " kW", "kva": " kVA", "v": "V", "a": "A",
		"mm": " mm", "cm": " cm", "m": " m", "kg": " kg", "g": " g",
		"l": " L", "ml": " mL", "bar": " bar", "rpm": " rpm",
	}

	// Try to extract numeric part and unit
	for suffix, formatted := range units {
		if strings.HasSuffix(strings.ToLower(label), suffix) {
			numPart := label[:len(label)-len(suffix)]
			// Replace underscores with dots for numeric values
			numPart = strings.ReplaceAll(numPart, "_", ".")
			deLbl := strings.ReplaceAll(numPart, ".", ",") + formatted
			enLbl := numPart + formatted
			return map[string]string{"de": deLbl, "en": enLbl}
		}
	}

	// No unit found: humanize the code (replace underscores, title case)
	label = strings.ReplaceAll(label, "_", " ")
	label = strings.ToUpper(label[:1]) + label[1:]
	return map[string]string{"de": label, "en": label}
}

// formatAxisLabel generates a human-readable i18n label from an axis attribute_code
// e.g. "power_rating" -> {"de": "Leistung", "en": "Power Rating"}
func formatAxisLabel(code string) map[string]string {
	knownAxes := map[string]map[string]string{
		"power_rating":  {"de": "Leistung", "en": "Power Rating"},
		"voltage":       {"de": "Spannung", "en": "Voltage"},
		"mounting_type": {"de": "Bauform", "en": "Mounting Type"},
		"size":          {"de": "Grösse", "en": "Size"},
		"color":         {"de": "Farbe", "en": "Color"},
		"material":      {"de": "Material", "en": "Material"},
		"weight":        {"de": "Gewicht", "en": "Weight"},
		"length":        {"de": "Länge", "en": "Length"},
		"width":         {"de": "Breite", "en": "Width"},
		"height":        {"de": "Höhe", "en": "Height"},
	}

	if labels, ok := knownAxes[code]; ok {
		return labels
	}

	// Fallback: humanize the code
	label := strings.ReplaceAll(code, "_", " ")
	words := strings.Fields(label)
	for i, w := range words {
		words[i] = strings.ToUpper(w[:1]) + w[1:]
	}
	humanized := strings.Join(words, " ")
	return map[string]string{"de": humanized, "en": humanized}
}

// GetProductWithVariants retrieves a product with all its variants (if parent)
func (r *ProductRepository) GetProductWithVariants(ctx context.Context, id uuid.UUID) (*domain.Product, error) {
	product, err := r.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// If not a variant parent, return as-is
	if product.ProductType != domain.ProductTypeVariantParent {
		return product, nil
	}

	// Load variant axes
	axes, err := r.GetVariantAxes(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to load variant axes: %w", err)
	}
	// Load all variants (need these first to build axis options)
	variants, err := r.ListVariants(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to load variants: %w", err)
	}

	// Convert to compact variant representation
	product.Variants = make([]domain.ProductVariant, len(variants))
	for i, variant := range variants {
		axisValues, _ := r.GetAxisValues(ctx, variant.ID)
		axisValuesMap := make(map[string]string)
		for _, av := range axisValues {
			axisValuesMap[av.AxisAttributeCode] = av.OptionCode
		}

		product.Variants[i] = domain.ProductVariant{
			ID:         variant.ID,
			SKU:        variant.SKU,
			AxisValues: axisValuesMap,
			Status:     variant.Status,
			Images:     variant.Images,
		}
	}

	// Build axis options from actual variant data and generate labels
	for i := range axes {
		optionSet := make(map[string]bool)
		var options []domain.AxisOption
		pos := 0
		for _, v := range product.Variants {
			if code, ok := v.AxisValues[axes[i].AttributeCode]; ok && !optionSet[code] {
				optionSet[code] = true
				options = append(options, domain.AxisOption{
					Code:     code,
					Label:    formatOptionLabel(code),
					Position: pos,
				})
				pos++
			}
		}
		axes[i].Options = options
		axes[i].Label = formatAxisLabel(axes[i].AttributeCode)
	}
	product.VariantAxes = axes

	return product, nil
}

// ListVariants returns all variants for a parent product
func (r *ProductRepository) ListVariants(ctx context.Context, parentID uuid.UUID, status ...domain.ProductStatus) ([]domain.Product, error) {
	query := `
		SELECT id, tenant_id, product_type, parent_id, sku, name, description, category_ids,
		       attributes, status, images, pim_identifier, last_synced_at,
		       created_at, updated_at, deleted_at
		FROM products
		WHERE parent_id = $1 AND deleted_at IS NULL
	`

	args := []any{parentID}

	if len(status) > 0 {
		query += ` AND status = ANY($2)`
		statusStrs := make([]string, len(status))
		for i, s := range status {
			statusStrs[i] = string(s)
		}
		args = append(args, statusStrs)
	}

	query += ` ORDER BY sku`

	rows, err := r.db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var variants []domain.Product
	for rows.Next() {
		variant, err := r.scanProductFromRows(rows)
		if err != nil {
			return nil, err
		}
		variants = append(variants, *variant)
	}

	return variants, rows.Err()
}

// FindVariantByAxisValues finds a variant matching the given axis values
func (r *ProductRepository) FindVariantByAxisValues(ctx context.Context, parentID uuid.UUID, axisValues map[string]string) (*domain.Product, error) {
	if len(axisValues) == 0 {
		return nil, fmt.Errorf("no axis values provided")
	}

	// Build a query that finds a variant with matching axis values
	query := `
		SELECT DISTINCT p.id, p.tenant_id, p.product_type, p.parent_id, p.sku, p.name, 
		       p.description, p.category_ids, p.attributes, p.status, p.images, 
		       p.pim_identifier, p.last_synced_at, p.created_at, p.updated_at, p.deleted_at
		FROM products p
		WHERE p.parent_id = $1 
		  AND p.product_type = 'variant'
		  AND p.deleted_at IS NULL
		  AND (
		      SELECT COUNT(*)
		      FROM variant_axis_values vav
		      JOIN variant_axes va ON vav.axis_id = va.id
		      WHERE vav.variant_id = p.id
	`

	args := []any{parentID}
	argNum := 2

	// Add conditions for each axis value
	conditions := make([]string, 0, len(axisValues))
	for attrCode, optCode := range axisValues {
		conditions = append(conditions, fmt.Sprintf("(va.attribute_code = $%d AND vav.option_code = $%d)", argNum, argNum+1))
		args = append(args, attrCode, optCode)
		argNum += 2
	}

	query += ` AND (` + joinWithOR(conditions) + `)`
	query += fmt.Sprintf(`) = $%d`, argNum)
	args = append(args, len(axisValues))

	return r.scanProduct(r.db.Pool.QueryRow(ctx, query, args...))
}

// GetAvailableAxisValues returns available axis option codes based on current selection
func (r *ProductRepository) GetAvailableAxisValues(ctx context.Context, parentID uuid.UUID, selected map[string]string) (map[string][]domain.AxisOption, error) {
	// Get all axes for this parent
	axes, err := r.GetVariantAxes(ctx, parentID)
	if err != nil {
		return nil, err
	}

	result := make(map[string][]domain.AxisOption)

	for _, axis := range axes {
		// If this axis is already selected, skip
		if _, ok := selected[axis.AttributeCode]; ok {
			continue
		}

		// Query for available option codes for this axis
		query := `
			SELECT DISTINCT vav.option_code
			FROM variant_axis_values vav
			JOIN variant_axes va ON vav.axis_id = va.id
			JOIN products p ON vav.variant_id = p.id
			WHERE va.product_id = $1
			  AND va.attribute_code = $2
			  AND p.status = 'active'
			  AND p.deleted_at IS NULL
		`

		args := []any{parentID, axis.AttributeCode}
		argNum := 3

		// Add filters for selected axes
		if len(selected) > 0 {
			query += ` AND p.id IN (
				SELECT vav_inner.variant_id
				FROM variant_axis_values vav_inner
				JOIN variant_axes va_inner ON vav_inner.axis_id = va_inner.id
				WHERE va_inner.product_id = $1
			`

			for selAttrCode, selOptCode := range selected {
				query += fmt.Sprintf(" AND (va_inner.attribute_code = $%d AND vav_inner.option_code = $%d)", argNum, argNum+1)
				args = append(args, selAttrCode, selOptCode)
				argNum += 2
			}

			query += ` GROUP BY vav_inner.variant_id
				HAVING COUNT(DISTINCT va_inner.attribute_code) = $` + fmt.Sprintf("%d", argNum)
			args = append(args, len(selected))
			query += `)`
		}

		rows, err := r.db.Pool.Query(ctx, query, args...)
		if err != nil {
			return nil, err
		}

		var options []domain.AxisOption
		for rows.Next() {
			var optCode string
			if err := rows.Scan(&optCode); err != nil {
				rows.Close()
				return nil, err
			}
			available := true
			options = append(options, domain.AxisOption{
				Code:      optCode,
				Available: &available,
			})
		}
		rows.Close()

		result[axis.AttributeCode] = options
	}

	return result, nil
}

// SetVariantAxes sets the variant axes for a parent product
func (r *ProductRepository) SetVariantAxes(ctx context.Context, parentID uuid.UUID, axes []domain.VariantAxis) error {
	// Start transaction
	tx, err := r.db.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Delete existing axes
	_, err = tx.Exec(ctx, "DELETE FROM variant_axes WHERE product_id = $1", parentID)
	if err != nil {
		return err
	}

	// Insert new axes
	for _, axis := range axes {
		_, err = tx.Exec(ctx, `
			INSERT INTO variant_axes (id, product_id, attribute_code, position)
			VALUES ($1, $2, $3, $4)
		`, axis.ID, parentID, axis.AttributeCode, axis.Position)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

// GetVariantAxes retrieves all variant axes for a parent product
func (r *ProductRepository) GetVariantAxes(ctx context.Context, parentID uuid.UUID) ([]domain.VariantAxis, error) {
	query := `
		SELECT id, product_id, attribute_code, position
		FROM variant_axes
		WHERE product_id = $1
		ORDER BY position
	`

	rows, err := r.db.Pool.Query(ctx, query, parentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var axes []domain.VariantAxis
	for rows.Next() {
		var axis domain.VariantAxis
		if err := rows.Scan(&axis.ID, &axis.ProductID, &axis.AttributeCode, &axis.Position); err != nil {
			return nil, err
		}
		axes = append(axes, axis)
	}

	return axes, rows.Err()
}

// SetAxisValues sets the axis values for a variant product
func (r *ProductRepository) SetAxisValues(ctx context.Context, variantID uuid.UUID, values []domain.AxisValueEntry) error {
	// Start transaction
	tx, err := r.db.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Delete existing values
	_, err = tx.Exec(ctx, "DELETE FROM variant_axis_values WHERE variant_id = $1", variantID)
	if err != nil {
		return err
	}

	// Insert new values
	for _, value := range values {
		_, err = tx.Exec(ctx, `
			INSERT INTO variant_axis_values (variant_id, axis_id, option_code)
			VALUES ($1, $2, $3)
		`, variantID, value.AxisID, value.OptionCode)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

// GetAxisValues retrieves all axis values for a variant product
func (r *ProductRepository) GetAxisValues(ctx context.Context, variantID uuid.UUID) ([]domain.AxisValueEntry, error) {
	query := `
		SELECT vav.variant_id, vav.axis_id, va.attribute_code, vav.option_code
		FROM variant_axis_values vav
		JOIN variant_axes va ON vav.axis_id = va.id
		WHERE vav.variant_id = $1
		ORDER BY va.position
	`

	rows, err := r.db.Pool.Query(ctx, query, variantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var values []domain.AxisValueEntry
	for rows.Next() {
		var value domain.AxisValueEntry
		if err := rows.Scan(&value.VariantID, &value.AxisID, &value.AxisAttributeCode, &value.OptionCode); err != nil {
			return nil, err
		}
		values = append(values, value)
	}

	return values, rows.Err()
}

// Helper function to join conditions with OR
func joinWithOR(conditions []string) string {
	if len(conditions) == 0 {
		return ""
	}
	result := conditions[0]
	for i := 1; i < len(conditions); i++ {
		result += " OR " + conditions[i]
	}
	return result
}

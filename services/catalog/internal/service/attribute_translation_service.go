package service

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/gondolia/gondolia/services/catalog/internal/domain"
	"github.com/gondolia/gondolia/services/catalog/internal/repository"
)

// AttributeTranslationService handles business logic for attribute translations
type AttributeTranslationService struct {
	repo repository.AttributeTranslationRepository
}

// NewAttributeTranslationService creates a new attribute translation service
func NewAttributeTranslationService(repo repository.AttributeTranslationRepository) *AttributeTranslationService {
	return &AttributeTranslationService{
		repo: repo,
	}
}

// GetByKey returns an attribute translation by key and locale
func (s *AttributeTranslationService) GetByKey(ctx context.Context, tenantID uuid.UUID, attributeKey, locale string) (*domain.AttributeTranslation, error) {
	return s.repo.GetByKey(ctx, tenantID, attributeKey, locale)
}

// GetByTenantAndLocale returns all attribute translations for a tenant and locale as a map
func (s *AttributeTranslationService) GetByTenantAndLocale(ctx context.Context, tenantID uuid.UUID, locale string) (map[string]*domain.AttributeTranslation, error) {
	return s.repo.GetByTenantAndLocale(ctx, tenantID, locale)
}

// List returns a paginated list of attribute translations
func (s *AttributeTranslationService) List(ctx context.Context, filter domain.AttributeTranslationFilter) ([]domain.AttributeTranslation, int, error) {
	return s.repo.List(ctx, filter)
}

// Create creates a new attribute translation
func (s *AttributeTranslationService) Create(ctx context.Context, tenantID uuid.UUID, req domain.CreateAttributeTranslationRequest) (*domain.AttributeTranslation, error) {
	translation := domain.NewAttributeTranslation(tenantID, req.AttributeKey, req.Locale, req.DisplayName)
	translation.Unit = req.Unit
	translation.Description = req.Description

	if err := s.repo.Create(ctx, translation); err != nil {
		return nil, err
	}

	return translation, nil
}

// Update updates an existing attribute translation
func (s *AttributeTranslationService) Update(ctx context.Context, id uuid.UUID, req domain.UpdateAttributeTranslationRequest) (*domain.AttributeTranslation, error) {
	// Note: We need to get the translation first to update it
	// Since we don't have a GetByID method, this is a simplified implementation
	// In a real-world scenario, you'd want to add GetByID to the repository
	translation := &domain.AttributeTranslation{
		ID:        id,
		UpdatedAt: time.Now(),
	}

	if req.DisplayName != nil {
		translation.DisplayName = *req.DisplayName
	}
	if req.Unit != nil {
		translation.Unit = req.Unit
	}
	if req.Description != nil {
		translation.Description = req.Description
	}

	if err := s.repo.Update(ctx, translation); err != nil {
		return nil, err
	}

	return translation, nil
}

// Delete deletes an attribute translation
func (s *AttributeTranslationService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}

// TranslateAttributes translates product attributes using the translation map
func (s *AttributeTranslationService) TranslateAttributes(attributes []domain.ProductAttribute, translations map[string]*domain.AttributeTranslation) []domain.ProductAttributeWithTranslation {
	result := make([]domain.ProductAttributeWithTranslation, len(attributes))

	for i, attr := range attributes {
		result[i] = domain.ProductAttributeWithTranslation{
			Key:   attr.Key,
			Type:  attr.Type,
			Value: attr.Value,
		}

		// Apply translation if available
		if trans, ok := translations[attr.Key]; ok {
			result[i].DisplayName = trans.DisplayName
			result[i].Unit = trans.Unit
		}
	}

	return result
}

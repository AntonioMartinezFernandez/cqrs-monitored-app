package domain_validation

import (
	"github.com/AntonioMartinezFernandez/cqrs-monitored-app/pkg/domain"
)

const validationErrorDefaultMessage = "validation error has been raised"

type ValidationError struct {
	items map[string]any

	domain.RootDomainError
}

func NewValidationErrorWithMetadata(keyedMetadata ...ValidationMetadata) *ValidationError {
	metadata := make(map[string]any, len(keyedMetadata))
	for _, raw := range keyedMetadata {
		metadata[raw.key] = raw.value
	}

	return &ValidationError{
		items: metadata,
	}
}

func (v *ValidationError) Error() string {
	return validationErrorDefaultMessage
}

func (v *ValidationError) ExtraItems() map[string]any {
	return v.items
}

type ValidationMetadata struct {
	key   string
	value any
}

func NewValidationMetadata(key string, value any) ValidationMetadata {
	return ValidationMetadata{
		key:   key,
		value: value,
	}
}

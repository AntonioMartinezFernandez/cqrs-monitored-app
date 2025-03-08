package domain_validation

import (
	"github.com/AntonioMartinezFernandez/cqrs-monitored-app/pkg/domain"
)

const rootDomainErrorsDetailsKey = "details"
const rootDomainValidationErrorMessage = "domain validation error occurred"

type DomainValidationError struct {
	domain.RootDomainError

	items map[string]any
}

func NewDomainValidationError(errors *ValidationErrors, domainErr error) *DomainValidationError {
	return &DomainValidationError{
		RootDomainError: domain.NewDomainErrorWithPrevious(domainErr),
		items:           buildDomainValidationErrorExtraItemsWithErrorDetails(errors),
	}
}

func (rve *DomainValidationError) Error() string {
	return rootDomainValidationErrorMessage
}

func (rve *DomainValidationError) ExtraItems() map[string]any {
	return rve.items
}

func (rve *DomainValidationError) ErrorDetails() []map[string]any {
	errors := make([]map[string]any, 0)

	if rve.items == nil {
		return errors
	}

	validationErrors, detailsKeyExists := rve.items[rootDomainErrorsDetailsKey]
	if !detailsKeyExists {
		return errors
	}

	errors = append(errors, validationErrors.([]map[string]any)...)

	return errors
}

func buildDomainValidationErrorExtraItemsWithErrorDetails(errors *ValidationErrors) map[string]any {
	details := map[string]any{
		rootDomainErrorsDetailsKey: make([]map[string]any, 0),
	}

	for _, e := range errors.Errors() {
		details[rootDomainErrorsDetailsKey] = append(
			details[rootDomainErrorsDetailsKey].([]map[string]any),
			e.ExtraItems(),
		)
	}

	return details
}

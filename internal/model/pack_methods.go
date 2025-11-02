package model

// Validate validates the PackRequest
func (r *PackRequest) Validate() error {
	if r.Quantity <= 0 {
		return NewValidationError("quantity must be greater than 0")
	}
	return nil
}

// GetValidPackSizes returns only valid (positive) pack sizes from the request
func (r *PackRequest) GetValidPackSizes() []int {
	if len(r.PackSizes) == 0 {
		return nil
	}

	validSizes := make([]int, 0, len(r.PackSizes))
	for _, size := range r.PackSizes {
		if size > 0 {
			validSizes = append(validSizes, size)
		}
	}
	return validSizes
}

// HasPackSizes returns true if the request contains pack sizes
func (r *PackRequest) HasPackSizes() bool {
	return len(r.PackSizes) > 0
}

// CalculateTotals calculates and sets TotalItems and TotalPacks from PackBreakdown
func (r *PackResponse) CalculateTotals() {
	totalItems := 0
	totalPacks := 0

	for packSize, count := range r.PackBreakdown {
		totalItems += packSize * count
		totalPacks += count
	}

	r.TotalItems = totalItems
	r.TotalPacks = totalPacks
}

// NewPackResponse creates a new PackResponse with calculated totals
func NewPackResponse(quantity int, breakdown map[int]int, packSizesUsed []int) *PackResponse {
	response := &PackResponse{
		Quantity:      quantity,
		PackBreakdown: breakdown,
		PackSizesUsed: packSizesUsed,
	}
	response.CalculateTotals()
	return response
}

// NewErrorResponse creates a new ErrorResponse
func NewErrorResponse(error, message string) ErrorResponse {
	return ErrorResponse{
		Error:   error,
		Message: message,
	}
}

// NewValidationError creates an error response for validation failures
func NewValidationError(message string) error {
	return &ValidationError{Message: message}
}

// ValidationError represents a validation error
type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}

// IsValidationError checks if an error is a ValidationError
func IsValidationError(err error) bool {
	_, ok := err.(*ValidationError)
	return ok
}

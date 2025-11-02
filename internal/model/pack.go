package model

// PackSize represents a pack size configuration
type PackSize struct {
	ID   int `json:"id"`
	Size int `json:"size"`
}

// PackRequest represents the request to calculate pack distribution
type PackRequest struct {
	Quantity  int   `json:"quantity" binding:"required,min=1"`
	PackSizes []int `json:"pack_sizes,omitempty"`
}

// PackResponse represents the response with pack distribution
type PackResponse struct {
	Quantity      int         `json:"quantity"`
	TotalItems    int         `json:"total_items"`
	TotalPacks    int         `json:"total_packs"`
	PackBreakdown map[int]int `json:"pack_breakdown"`
	PackSizesUsed []int       `json:"pack_sizes_used"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

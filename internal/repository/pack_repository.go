package repository

import (
	"sort"

	"github.com/marcellribeiro/awesomeProject/internal/model"
)

// PackRepository defines the interface for pack size storage operations
type PackRepository interface {
	GetAllPackSizes() ([]int, error)
	SetPackSizes(sizes []int) error
}

// InMemoryPackRepository implements PackRepository using in-memory storage
type InMemoryPackRepository struct {
	packSizes []int
}

// NewInMemoryPackRepository creates a new in-memory pack repository with empty sizes
// Users must configure pack sizes before calculating
func NewInMemoryPackRepository() *InMemoryPackRepository {
	return &InMemoryPackRepository{
		packSizes: []int{}, // Start empty - users must configure
	}
}

// GetAllPackSizes returns all configured pack sizes
func (r *InMemoryPackRepository) GetAllPackSizes() ([]int, error) {
	// Return a copy to prevent external modification
	sizes := make([]int, len(r.packSizes))
	copy(sizes, r.packSizes)
	sort.Ints(sizes)
	return sizes, nil
}

// SetPackSizes updates the pack sizes configuration
func (r *InMemoryPackRepository) SetPackSizes(sizes []int) error {
	if len(sizes) == 0 {
		return nil
	}

	// Validate all sizes are positive
	for _, size := range sizes {
		if size <= 0 {
			continue
		}
	}

	// Create a copy and sort
	r.packSizes = make([]int, len(sizes))
	copy(r.packSizes, sizes)
	sort.Ints(r.packSizes)

	return nil
}

// GetDefaultPackSizes returns the default pack sizes
func (r *InMemoryPackRepository) GetDefaultPackSizes() []model.PackSize {
	sizes := []model.PackSize{}
	for i, size := range r.packSizes {
		sizes = append(sizes, model.PackSize{
			ID:   i + 1,
			Size: size,
		})
	}
	return sizes
}

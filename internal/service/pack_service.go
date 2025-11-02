package service

import (
	"fmt"
	"sort"

	"github.com/marcellribeiro/awesomeProject/internal/model"
	"github.com/marcellribeiro/awesomeProject/internal/repository"
	"github.com/marcellribeiro/awesomeProject/pkg/calculator"
)

// PackService defines the interface for pack calculation business logic
type PackService interface {
	CalculatePackDistribution(request *model.PackRequest) (*model.PackResponse, error)
	GetAvailablePackSizes() ([]int, error)
	UpdatePackSizes(sizes []int) error
}

// packService implements PackService
type packService struct {
	calculator calculator.PackCalculator
	repository repository.PackRepository
}

// NewPackService creates a new pack service instance
func NewPackService(calc calculator.PackCalculator, repo repository.PackRepository) PackService {
	return &packService{
		calculator: calc,
		repository: repo,
	}
}

// CalculatePackDistribution calculates the optimal pack distribution for a given quantity
func (s *packService) CalculatePackDistribution(request *model.PackRequest) (*model.PackResponse, error) {
	// Validate request
	if err := request.Validate(); err != nil {
		return nil, err
	}

	// Get pack sizes (use provided ones or fetch from repository)
	var packSizes []int
	var err error

	if request.HasPackSizes() {
		packSizes = request.GetValidPackSizes()
	} else {
		packSizes, err = s.repository.GetAllPackSizes()
		if err != nil {
			return nil, fmt.Errorf("failed to get pack sizes: %w", err)
		}
	}

	if len(packSizes) == 0 {
		return nil, model.NewValidationError("no valid pack sizes available")
	}

	// Calculate optimal distribution
	breakdown, err := s.calculator.Calculate(request.Quantity, packSizes)
	if err != nil {
		return nil, fmt.Errorf("calculation failed: %w", err)
	}

	// Build response with calculated totals
	return model.NewPackResponse(request.Quantity, breakdown, packSizes), nil
}

// GetAvailablePackSizes returns all configured pack sizes
func (s *packService) GetAvailablePackSizes() ([]int, error) {
	sizes, err := s.repository.GetAllPackSizes()
	if err != nil {
		return nil, fmt.Errorf("failed to get pack sizes: %w", err)
	}
	sort.Ints(sizes)
	return sizes, nil
}

// UpdatePackSizes updates the configured pack sizes
func (s *packService) UpdatePackSizes(sizes []int) error {
	if len(sizes) == 0 {
		return fmt.Errorf("pack sizes cannot be empty")
	}

	// Validate all sizes are positive
	for _, size := range sizes {
		if size <= 0 {
			return fmt.Errorf("all pack sizes must be positive, got: %d", size)
		}
	}

	return s.repository.SetPackSizes(sizes)
}

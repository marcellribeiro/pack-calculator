package service

import (
	"testing"

	"github.com/marcellribeiro/awesomeProject/internal/model"
	"github.com/marcellribeiro/awesomeProject/internal/repository"
	"github.com/marcellribeiro/awesomeProject/pkg/calculator"
)

func TestPackService_CalculatePackDistribution(t *testing.T) {
	repo := repository.NewInMemoryPackRepository()
	calc := calculator.NewDynamicPackCalculator()
	service := NewPackService(calc, repo)

	// Setup default pack sizes for tests that don't provide custom sizes
	repo.SetPackSizes([]int{250, 500, 1000, 2000, 5000})

	tests := []struct {
		name      string
		request   *model.PackRequest
		wantError bool
	}{
		{
			name: "Valid request with default pack sizes",
			request: &model.PackRequest{
				Quantity: 251,
			},
			wantError: false,
		},
		{
			name: "Valid request with custom pack sizes",
			request: &model.PackRequest{
				Quantity:  100,
				PackSizes: []int{25, 50, 100},
			},
			wantError: false,
		},
		{
			name: "Invalid quantity - zero",
			request: &model.PackRequest{
				Quantity: 0,
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := service.CalculatePackDistribution(tt.request)

			if tt.wantError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if result == nil {
				t.Errorf("Expected result but got nil")
				return
			}

			if result.TotalItems < tt.request.Quantity {
				t.Errorf("Total items %v is less than requested %v", result.TotalItems, tt.request.Quantity)
			}
		})
	}
}

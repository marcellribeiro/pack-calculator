package calculator

import "testing"

func TestDynamicPackCalculator_Calculate(t *testing.T) {
	calc := NewDynamicPackCalculator()

	tests := []struct {
		name      string
		quantity  int
		packSizes []int
		wantTotal int
		wantPacks int
	}{
		{
			name:      "Order 1 item",
			quantity:  1,
			packSizes: []int{250, 500, 1000, 2000, 5000},
			wantTotal: 250,
			wantPacks: 1,
		},
		{
			name:      "Order 250 items",
			quantity:  250,
			packSizes: []int{250, 500, 1000, 2000, 5000},
			wantTotal: 250,
			wantPacks: 1,
		},
		{
			name:      "Order 251 items",
			quantity:  251,
			packSizes: []int{250, 500, 1000, 2000, 5000},
			wantTotal: 500,
			wantPacks: 1,
		},
		{
			name:      "Order 501 items",
			quantity:  501,
			packSizes: []int{250, 500, 1000, 2000, 5000},
			wantTotal: 750,
			wantPacks: 2,
		},
		{
			name:      "Order 12001 items",
			quantity:  12001,
			packSizes: []int{250, 500, 1000, 2000, 5000},
			wantTotal: 12250,
			wantPacks: 4,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := calc.Calculate(tt.quantity, tt.packSizes)
			if err != nil {
				t.Errorf("Calculate() error = %v", err)
				return
			}

			totalItems := 0
			totalPacks := 0
			for packSize, count := range result {
				totalItems += packSize * count
				totalPacks += count
			}

			if totalItems != tt.wantTotal {
				t.Errorf("Calculate() total items = %v, want %v. Breakdown: %v", totalItems, tt.wantTotal, result)
			}

			if totalPacks != tt.wantPacks {
				t.Errorf("Calculate() total packs = %v, want %v. Breakdown: %v", totalPacks, tt.wantPacks, result)
			}
		})
	}
}

func TestDynamicPackCalculator_EdgeCase(t *testing.T) {
	calc := NewDynamicPackCalculator()

	packSizes := []int{23, 31, 53}
	quantity := 500000

	result, err := calc.Calculate(quantity, packSizes)
	if err != nil {
		t.Fatalf("Calculate() error = %v", err)
	}

	actualTotal := 0
	for packSize, count := range result {
		actualTotal += packSize * count
	}

	t.Logf("Edge case result: %v", result)
	t.Logf("Quantity requested: %d", quantity)
	t.Logf("Total items shipped: %d", actualTotal)

	if actualTotal < quantity {
		t.Errorf("Shipped %d items but needed at least %d", actualTotal, quantity)
	}
}

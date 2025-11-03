package model

import (
	"errors"
	"reflect"
	"testing"
)

func TestPackRequest_Validate(t *testing.T) {
	tests := []struct {
		name     string
		quantity int
		wantErr  bool
	}{
		{"Valid quantity", 100, false},
		{"Valid quantity 1", 1, false},
		{"Zero quantity", 0, true},
		{"Negative quantity", -100, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &PackRequest{Quantity: tt.quantity}
			err := req.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPackRequest_GetValidPackSizes(t *testing.T) {
	tests := []struct {
		name      string
		packSizes []int
		want      []int
	}{
		{"Empty slice", []int{}, nil},
		{"Nil slice", nil, nil},
		{"Valid sizes", []int{250, 500, 1000}, []int{250, 500, 1000}},
		{"Mixed valid/invalid", []int{250, 0, 500, -100}, []int{250, 500}},
		{"Only invalid", []int{0, -100}, []int{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &PackRequest{Quantity: 100, PackSizes: tt.packSizes}
			got := req.GetValidPackSizes()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetValidPackSizes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPackRequest_HasPackSizes(t *testing.T) {
	tests := []struct {
		name      string
		packSizes []int
		want      bool
	}{
		{"Has sizes", []int{250, 500}, true},
		{"Empty slice", []int{}, false},
		{"Nil slice", nil, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &PackRequest{Quantity: 100, PackSizes: tt.packSizes}
			if got := req.HasPackSizes(); got != tt.want {
				t.Errorf("HasPackSizes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPackResponse_CalculateTotals(t *testing.T) {
	tests := []struct {
		name      string
		breakdown map[int]int
		wantItems int
		wantPacks int
	}{
		{"Empty", map[int]int{}, 0, 0},
		{"Single", map[int]int{250: 1}, 250, 1},
		{"Multiple", map[int]int{1000: 1, 250: 1}, 1250, 2},
		{"Same size", map[int]int{250: 4}, 1000, 4},
		{"Nil breakdown", nil, 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := &PackResponse{PackBreakdown: tt.breakdown}
			resp.CalculateTotals()
			if resp.TotalItems != tt.wantItems || resp.TotalPacks != tt.wantPacks {
				t.Errorf("CalculateTotals() = items:%d packs:%d, want items:%d packs:%d",
					resp.TotalItems, resp.TotalPacks, tt.wantItems, tt.wantPacks)
			}
		})
	}
}

func TestNewPackResponse(t *testing.T) {
	breakdown := map[int]int{250: 1}
	packSizes := []int{250, 500}

	resp := NewPackResponse(250, breakdown, packSizes)

	if resp == nil {
		t.Fatal("NewPackResponse() returned nil")
	}
	if resp.Quantity != 250 {
		t.Errorf("Quantity = %d, want 250", resp.Quantity)
	}
	if resp.TotalItems != 250 || resp.TotalPacks != 1 {
		t.Errorf("Totals not calculated correctly")
	}
}

func TestNewErrorResponse(t *testing.T) {
	resp := NewErrorResponse("error", "message")
	if resp.Error != "error" || resp.Message != "message" {
		t.Errorf("NewErrorResponse() = %v, want error:error message:message", resp)
	}
}

func TestNewValidationError(t *testing.T) {
	msg := "validation failed"
	err := NewValidationError(msg)

	if err == nil {
		t.Fatal("NewValidationError() returned nil")
	}
	if err.Error() != msg {
		t.Errorf("Error() = %v, want %v", err.Error(), msg)
	}
	if !IsValidationError(err) {
		t.Error("Not recognized as ValidationError")
	}
}

func TestIsValidationError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{"ValidationError", &ValidationError{Message: "test"}, true},
		{"Generic error", errors.New("error"), false},
		{"Nil", nil, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidationError(tt.err); got != tt.want {
				t.Errorf("IsValidationError() = %v, want %v", got, tt.want)
			}
		})
	}
}

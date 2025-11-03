package repository

import (
	"reflect"
	"testing"
)

func TestNewInMemoryPackRepository(t *testing.T) {
	repo := NewInMemoryPackRepository()

	if repo == nil {
		t.Fatal("Expected repository to be created, got nil")
	}
	if len(repo.packSizes) != 0 {
		t.Errorf("Expected empty packSizes, got length %d", len(repo.packSizes))
	}
}

func TestInMemoryPackRepository_GetAllPackSizes(t *testing.T) {
	tests := []struct {
		name     string
		initial  []int
		expected []int
	}{
		{"Empty", []int{}, []int{}},
		{"Single", []int{250}, []int{250}},
		{"Sorted", []int{250, 500, 1000}, []int{250, 500, 1000}},
		{"Unsorted", []int{1000, 250, 500}, []int{250, 500, 1000}},
		{"Duplicates", []int{250, 500, 250}, []int{250, 250, 500}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &InMemoryPackRepository{packSizes: tt.initial}
			sizes, err := repo.GetAllPackSizes()

			if err != nil {
				t.Errorf("GetAllPackSizes() error = %v", err)
			}
			if !reflect.DeepEqual(sizes, tt.expected) {
				t.Errorf("GetAllPackSizes() = %v, want %v", sizes, tt.expected)
			}
		})
	}
}

func TestInMemoryPackRepository_SetPackSizes(t *testing.T) {
	tests := []struct {
		name     string
		input    []int
		expected []int
	}{
		{"Empty", []int{}, []int{}},
		{"Single", []int{250}, []int{250}},
		{"Multiple", []int{250, 500, 1000}, []int{250, 500, 1000}},
		{"Unsorted", []int{5000, 250, 1000}, []int{250, 1000, 5000}},
		{"Duplicates", []int{250, 500, 250}, []int{250, 250, 500}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewInMemoryPackRepository()
			err := repo.SetPackSizes(tt.input)

			if err != nil {
				t.Errorf("SetPackSizes() error = %v", err)
			}

			sizes, _ := repo.GetAllPackSizes()
			if !reflect.DeepEqual(sizes, tt.expected) {
				t.Errorf("After SetPackSizes(), got %v, want %v", sizes, tt.expected)
			}
		})
	}
}

func TestInMemoryPackRepository_GetDefaultPackSizes(t *testing.T) {
	tests := []struct {
		name  string
		sizes []int
		count int
	}{
		{"Empty", []int{}, 0},
		{"Single", []int{250}, 1},
		{"Multiple", []int{250, 500, 1000}, 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &InMemoryPackRepository{packSizes: tt.sizes}
			result := repo.GetDefaultPackSizes()

			if len(result) != tt.count {
				t.Errorf("GetDefaultPackSizes() returned %d items, want %d", len(result), tt.count)
			}

			for i, ps := range result {
				if ps.ID != i+1 {
					t.Errorf("PackSize[%d].ID = %d, want %d", i, ps.ID, i+1)
				}
				if ps.Size != tt.sizes[i] {
					t.Errorf("PackSize[%d].Size = %d, want %d", i, ps.Size, tt.sizes[i])
				}
			}
		})
	}
}

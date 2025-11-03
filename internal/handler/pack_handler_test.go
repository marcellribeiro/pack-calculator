package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/marcellribeiro/awesomeProject/internal/model"
)

// Mock service for testing
type mockPackService struct {
	calculateFunc       func(request *model.PackRequest) (*model.PackResponse, error)
	getPackSizesFunc    func() ([]int, error)
	updatePackSizesFunc func(sizes []int) error
}

func (m *mockPackService) CalculatePackDistribution(request *model.PackRequest) (*model.PackResponse, error) {
	if m.calculateFunc != nil {
		return m.calculateFunc(request)
	}
	return nil, errors.New("not implemented")
}

func (m *mockPackService) GetAvailablePackSizes() ([]int, error) {
	if m.getPackSizesFunc != nil {
		return m.getPackSizesFunc()
	}
	return nil, errors.New("not implemented")
}

func (m *mockPackService) UpdatePackSizes(sizes []int) error {
	if m.updatePackSizesFunc != nil {
		return m.updatePackSizesFunc(sizes)
	}
	return errors.New("not implemented")
}

func TestNewPackHandler(t *testing.T) {
	mockService := &mockPackService{}
	handler := NewPackHandler(mockService)

	if handler == nil {
		t.Fatal("Expected handler to be created, got nil")
	}

	if handler.service == nil {
		t.Error("Expected service to be set, got nil")
	}
}

func TestPackHandler_CalculatePacks(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		requestBody    map[string]interface{}
		mockResponse   *model.PackResponse
		mockError      error
		expectedStatus int
		checkError     bool
	}{
		{
			name: "Valid request - successful calculation",
			requestBody: map[string]interface{}{
				"quantity": 1001,
			},
			mockResponse: &model.PackResponse{
				Quantity:   1001,
				TotalItems: 1250,
				TotalPacks: 2,
			},
			mockError:      nil,
			expectedStatus: http.StatusOK,
			checkError:     false,
		},
		{
			name: "Valid request with pack sizes",
			requestBody: map[string]interface{}{
				"quantity":   250,
				"pack_sizes": []int{250, 500, 1000},
			},
			mockResponse: &model.PackResponse{
				Quantity:   250,
				TotalItems: 250,
				TotalPacks: 1,
			},
			mockError:      nil,
			expectedStatus: http.StatusOK,
			checkError:     false,
		},
		{
			name: "Service returns error",
			requestBody: map[string]interface{}{
				"quantity": 1001,
			},
			mockResponse:   nil,
			mockError:      errors.New("calculation failed"),
			expectedStatus: http.StatusBadRequest,
			checkError:     true,
		},
		{
			name: "Invalid quantity - zero",
			requestBody: map[string]interface{}{
				"quantity": 0,
			},
			mockResponse:   nil,
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
			checkError:     true,
		},
		{
			name: "Invalid quantity - negative",
			requestBody: map[string]interface{}{
				"quantity": -100,
			},
			mockResponse:   nil,
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
			checkError:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &mockPackService{
				calculateFunc: func(request *model.PackRequest) (*model.PackResponse, error) {
					if tt.mockError != nil {
						return nil, tt.mockError
					}
					return tt.mockResponse, nil
				},
			}

			handler := NewPackHandler(mockService)

			// Setup Gin context
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			// Create request body
			body, _ := json.Marshal(tt.requestBody)
			c.Request, _ = http.NewRequest("POST", "/api/calculate", bytes.NewBuffer(body))
			c.Request.Header.Set("Content-Type", "application/json")

			// Call handler directly
			handler.CalculatePacks(c)

			// Check status code
			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			// Check response
			var response map[string]interface{}
			json.Unmarshal(w.Body.Bytes(), &response)

			if tt.checkError {
				if _, hasError := response["error"]; !hasError {
					t.Error("Expected error in response")
				}
			} else {
				if quantity, ok := response["quantity"]; ok {
					if int(quantity.(float64)) != tt.mockResponse.Quantity {
						t.Errorf("Expected quantity %d, got %v", tt.mockResponse.Quantity, quantity)
					}
				}
			}
		})
	}
}

func TestPackHandler_GetPackSizes(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		mockSizes      []int
		mockError      error
		expectedStatus int
	}{
		{
			name:           "Successfully get pack sizes",
			mockSizes:      []int{250, 500, 1000},
			mockError:      nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Empty pack sizes",
			mockSizes:      []int{},
			mockError:      nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Service returns error",
			mockSizes:      nil,
			mockError:      errors.New("database error"),
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:           "Single pack size",
			mockSizes:      []int{250},
			mockError:      nil,
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &mockPackService{
				getPackSizesFunc: func() ([]int, error) {
					return tt.mockSizes, tt.mockError
				},
			}

			handler := NewPackHandler(mockService)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request, _ = http.NewRequest("GET", "/api/pack-sizes", nil)

			handler.GetPackSizes(c)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			var response map[string]interface{}
			json.Unmarshal(w.Body.Bytes(), &response)

			if tt.mockError == nil {
				if packSizes, ok := response["pack_sizes"]; ok {
					sizes := packSizes.([]interface{})
					if len(sizes) != len(tt.mockSizes) {
						t.Errorf("Expected %d pack sizes, got %d", len(tt.mockSizes), len(sizes))
					}
				}
			}
		})
	}
}

func TestPackHandler_UpdatePackSizes(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		requestBody    map[string]interface{}
		mockError      error
		expectedStatus int
	}{
		{
			name: "Successfully update pack sizes",
			requestBody: map[string]interface{}{
				"pack_sizes": []int{250, 500, 1000},
			},
			mockError:      nil,
			expectedStatus: http.StatusOK,
		},
		{
			name: "Update with single pack size",
			requestBody: map[string]interface{}{
				"pack_sizes": []int{250},
			},
			mockError:      nil,
			expectedStatus: http.StatusOK,
		},
		{
			name: "Service returns error",
			requestBody: map[string]interface{}{
				"pack_sizes": []int{250, 500, 1000},
			},
			mockError:      errors.New("validation error"),
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Missing pack_sizes field",
			requestBody: map[string]interface{}{
				"wrong_field": []int{250, 500},
			},
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Empty pack sizes array",
			requestBody: map[string]interface{}{
				"pack_sizes": []int{},
			},
			mockError:      errors.New("pack sizes cannot be empty"),
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &mockPackService{
				updatePackSizesFunc: func(sizes []int) error {
					return tt.mockError
				},
			}

			handler := NewPackHandler(mockService)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			body, _ := json.Marshal(tt.requestBody)
			c.Request, _ = http.NewRequest("PUT", "/api/pack-sizes", bytes.NewBuffer(body))
			c.Request.Header.Set("Content-Type", "application/json")

			handler.UpdatePackSizes(c)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			var response map[string]interface{}
			json.Unmarshal(w.Body.Bytes(), &response)

			if tt.expectedStatus == http.StatusOK {
				if msg, ok := response["message"]; !ok || msg != "Pack sizes updated successfully" {
					t.Error("Expected success message")
				}
			} else {
				if _, hasError := response["error"]; !hasError {
					t.Error("Expected error in response")
				}
			}
		})
	}
}

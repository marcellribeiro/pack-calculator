package handler

// API Documentation using JSON Schema

type APIDocumentation struct {
	OpenAPI    string                        `json:"openapi"`
	Info       APIInfo                       `json:"info"`
	Servers    []APIServer                   `json:"servers"`
	Paths      map[string]map[string]APIPath `json:"paths"`
	Components APIComponents                 `json:"components"`
}

type APIInfo struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Version     string `json:"version"`
}

type APIServer struct {
	URL         string `json:"url"`
	Description string `json:"description"`
}

type APIPath struct {
	Summary     string                 `json:"summary"`
	Description string                 `json:"description,omitempty"`
	RequestBody *APIRequestBody        `json:"requestBody,omitempty"`
	Parameters  []APIParameter         `json:"parameters,omitempty"`
	Responses   map[string]APIResponse `json:"responses"`
}

type APIRequestBody struct {
	Required bool                  `json:"required"`
	Content  map[string]APIContent `json:"content"`
}

type APIContent struct {
	Schema APISchema `json:"schema"`
}

type APISchema struct {
	Ref        string                 `json:"$ref,omitempty"`
	Type       string                 `json:"type,omitempty"`
	Properties map[string]APIProperty `json:"properties,omitempty"`
	Required   []string               `json:"required,omitempty"`
	Items      *APISchema             `json:"items,omitempty"`
}

type APIProperty struct {
	Type        string      `json:"type"`
	Description string      `json:"description,omitempty"`
	Example     interface{} `json:"example,omitempty"`
	Minimum     *int        `json:"minimum,omitempty"`
}

type APIParameter struct {
	Name        string    `json:"name"`
	In          string    `json:"in"`
	Required    bool      `json:"required"`
	Schema      APISchema `json:"schema"`
	Description string    `json:"description,omitempty"`
}

type APIResponse struct {
	Description string                `json:"description"`
	Content     map[string]APIContent `json:"content,omitempty"`
}

type APIComponents struct {
	Schemas map[string]APISchema `json:"schemas"`
}

func GetAPIDocumentation() APIDocumentation {
	minOne := 1

	return APIDocumentation{
		OpenAPI: "3.0.0",
		Info: APIInfo{
			Title:       "Pack Calculator API",
			Description: "API for calculating optimal pack distributions based on configurable pack sizes",
			Version:     "1.0.0",
		},
		Servers: []APIServer{
			{
				URL:         "http://localhost:8080",
				Description: "Local development server",
			},
		},
		Paths: map[string]map[string]APIPath{
			"/health": {
				"get": {
					Summary:     "Health Check",
					Description: "Check if the service is running",
					Responses: map[string]APIResponse{
						"200": {
							Description: "Service is healthy",
							Content: map[string]APIContent{
								"application/json": {
									Schema: APISchema{
										Type: "object",
										Properties: map[string]APIProperty{
											"status": {
												Type:    "string",
												Example: "healthy",
											},
											"service": {
												Type:    "string",
												Example: "pack-calculator",
											},
										},
									},
								},
							},
						},
					},
				},
			},
			"/api/pack-sizes": {
				"get": {
					Summary:     "Get Pack Sizes",
					Description: "Retrieve all configured pack sizes",
					Responses: map[string]APIResponse{
						"200": {
							Description: "List of pack sizes",
							Content: map[string]APIContent{
								"application/json": {
									Schema: APISchema{
										Type: "object",
										Properties: map[string]APIProperty{
											"pack_sizes": {
												Type:    "array",
												Example: []int{250, 500, 1000},
											},
										},
									},
								},
							},
						},
					},
				},
				"put": {
					Summary:     "Update Pack Sizes",
					Description: "Update the configured pack sizes",
					RequestBody: &APIRequestBody{
						Required: true,
						Content: map[string]APIContent{
							"application/json": {
								Schema: APISchema{
									Ref: "#/components/schemas/UpdatePackSizesRequest",
								},
							},
						},
					},
					Responses: map[string]APIResponse{
						"200": {
							Description: "Pack sizes updated successfully",
							Content: map[string]APIContent{
								"application/json": {
									Schema: APISchema{
										Type: "object",
										Properties: map[string]APIProperty{
											"message": {
												Type:    "string",
												Example: "Pack sizes updated successfully",
											},
											"pack_sizes": {
												Type:    "array",
												Example: []int{250, 500, 1000},
											},
										},
									},
								},
							},
						},
						"400": {
							Description: "Invalid request",
							Content: map[string]APIContent{
								"application/json": {
									Schema: APISchema{
										Ref: "#/components/schemas/ErrorResponse",
									},
								},
							},
						},
					},
				},
			},
			"/api/calculate": {
				"post": {
					Summary:     "Calculate Pack Distribution",
					Description: "Calculate the optimal pack distribution for a given quantity. Rules: 1) Only whole packs 2) Minimize total items 3) Minimize number of packs",
					RequestBody: &APIRequestBody{
						Required: true,
						Content: map[string]APIContent{
							"application/json": {
								Schema: APISchema{
									Ref: "#/components/schemas/PackRequest",
								},
							},
						},
					},
					Responses: map[string]APIResponse{
						"200": {
							Description: "Calculation successful",
							Content: map[string]APIContent{
								"application/json": {
									Schema: APISchema{
										Ref: "#/components/schemas/PackResponse",
									},
								},
							},
						},
						"400": {
							Description: "Invalid request or calculation failed",
							Content: map[string]APIContent{
								"application/json": {
									Schema: APISchema{
										Ref: "#/components/schemas/ErrorResponse",
									},
								},
							},
						},
					},
				},
			},
		},
		Components: APIComponents{
			Schemas: map[string]APISchema{
				"PackRequest": {
					Type: "object",
					Properties: map[string]APIProperty{
						"quantity": {
							Type:        "integer",
							Description: "Number of items to order",
							Example:     251,
							Minimum:     &minOne,
						},
						"pack_sizes": {
							Type:        "array",
							Description: "Optional custom pack sizes (if not provided, uses configured pack sizes)",
							Example:     []int{250, 500, 1000},
						},
					},
					Required: []string{"quantity"},
				},
				"PackResponse": {
					Type: "object",
					Properties: map[string]APIProperty{
						"quantity": {
							Type:        "integer",
							Description: "Original requested quantity",
							Example:     251,
						},
						"total_items": {
							Type:        "integer",
							Description: "Total items that will be shipped",
							Example:     500,
						},
						"total_packs": {
							Type:        "integer",
							Description: "Total number of packs",
							Example:     1,
						},
						"pack_breakdown": {
							Type:        "object",
							Description: "Breakdown of packs by size",
							Example:     map[string]int{"500": 1},
						},
						"pack_sizes_used": {
							Type:        "array",
							Description: "Pack sizes that were used for calculation",
							Example:     []int{250, 500, 1000},
						},
					},
				},
				"UpdatePackSizesRequest": {
					Type: "object",
					Properties: map[string]APIProperty{
						"pack_sizes": {
							Type:        "array",
							Description: "Array of pack sizes (positive integers)",
							Example:     []int{250, 500, 1000, 2000, 5000},
						},
					},
					Required: []string{"pack_sizes"},
				},
				"ErrorResponse": {
					Type: "object",
					Properties: map[string]APIProperty{
						"error": {
							Type:        "string",
							Description: "Error message",
							Example:     "Invalid request",
						},
						"message": {
							Type:        "string",
							Description: "Detailed error message",
							Example:     "quantity must be greater than 0",
						},
					},
				},
			},
		},
	}
}

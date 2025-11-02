# 📦 Pack Calculator

A pack optimization system built with **Go** and **Gin framework** that calculates the most efficient pack distribution for customer orders, following strict business rules to minimize both items and pack count.

**Live Demo:** 🚀 [Coming Soon - Will be deployed]

---

## 🎯 Challenge Overview

The store ships products in various pack sizes, and customers can order any quantity. The system must:

### 📋 Rules (Priority Order)

1. **Rule 1:** Only whole packs can be sent (no breaking packs)
2. **Rule 2:** Send the **least amount of items** to fulfill the order ⚠️ **PRIORITY**
3. **Rule 3:** Send **as few packs as possible** (Rule 2 takes precedence)

### ✅ Example Results

| Items Ordered | Result | Explanation |
|--------------|--------|-------------|
| 1 | 1 × 250 | Smallest pack that covers the order |
| 250 | 1 × 250 | Exact match |
| 251 | 1 × 500 | Fewer items than 2×250 |
| 501 | 1×500 + 1×250 | 750 items total, 2 packs |
| 12001 | 2×5000 + 1×2000 + 1×250 | 12,250 items, 4 packs |

---

## 🏗️ Architecture

This project follows **Clean Architecture** and **SOLID principles**:

```
📁 Project Structure
├── cmd/
│   └── api/
│       └── main.go              # Application entry point
├── internal/
│   ├── handler/                 # HTTP handlers (Presentation layer)
│   │   ├── api_docs.go
│   │   └── pack_handler.go
│   ├── router/                  # Router setup
│   │   └── router.go
│   ├── service/                 # Business logic (Use case layer)
│   │   ├── pack_service.go
│   │   └── pack_service_test.go
│   ├── repository/              # Data access (Interface adapter)
│   │   └── pack_repository.go
│   └── model/                   # Domain models and helpers
│       ├── pack.go
│       └── pack_methods.go
├── pkg/
│   └── calculator/              # Core algorithm (Domain layer)
│       ├── pack_calculator.go
│       └── pack_calculator_test.go
├── web/
│   ├── templates/               # HTML templates
│   │   └── index.html
│   └── static/                  # CSS, JS, images
├── Dockerfile
├── docker-compose.yml
├── go.mod
└── README.md
```

### 🎨 Design Patterns Used

- **Dependency Injection:** All components receive dependencies via constructors
- **Repository Pattern:** Abstracts data storage from business logic
- **Strategy Pattern:** Calculator interface allows different algorithms
- **Interface Segregation:** Small, focused interfaces (PackRepository, PackCalculator, PackService)

---

## 🚀 Getting Started

### Prerequisites

- Go 1.21 or higher
- Docker & Docker Compose (optional, for containerized deployment)

### Option 1: Run with Go

```bash
# Clone the repository
git clone https://github.com/marcellribeiro/awesomeProject.git
cd awesomeProject

# Install dependencies
go mod download

# Run the application
go run cmd/api/main.go
```

Visit: `http://localhost:8080`

### Option 2: Run with Docker

```bash
# Build and run with Docker Compose
docker-compose up --build

# Or with plain Docker
docker build -t pack-calculator .
docker run -p 8080:8080 pack-calculator
```

Visit: `http://localhost:8080`

---

## 🧪 Testing

### Run All Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with verbose output
go test -v ./...

# Run specific test
go test -v ./pkg/calculator -run TestDynamicPackCalculator_EdgeCase
```

### Critical Edge Case Test

The system is validated against the critical edge case:

```go
Pack Sizes: [23, 31, 53]
Quantity: 500,000
Expected: {23: 2, 31: 7, 53: 9429}
Result: 500,000 items exactly
```

Run this specific test:
```bash
go test -v ./pkg/calculator -run EdgeCase
```

### Benchmark Tests

```bash
# Run benchmarks
go test -bench=. ./pkg/calculator

# Example output:
# BenchmarkCalculate_Small-8    500000    2.5 ns/op
# BenchmarkCalculate_Large-8      1000  1200 ns/op
```

---

## 📡 API Documentation

http://localhost:8080/docs

---

## 🎨 Web Interface

The application includes a web interface accessible at `http://localhost:8080`

---

## 🧮 Algorithm Explanation

### Dynamic Programming Approach

The calculator uses a two-phase dynamic programming algorithm:

#### Phase 1: Find Minimum Items
```
Goal: Find the smallest total ≥ requested quantity that can be achieved

Algorithm:
- Create DP array where dp[i] = true if amount i is achievable
- For each amount from 1 to max:
  - For each pack size:
    - If (amount - pack_size) is achievable, mark amount as achievable
- Return the first achievable amount ≥ requested quantity
```

#### Phase 2: Find Minimum Packs
```
Goal: Find the minimum number of packs to achieve the target from Phase 1

Algorithm:
- Create DP array where dp[i] = minimum packs to achieve amount i
- Track which pack was used to reach each amount (parent array)
- Backtrack from target to reconstruct the optimal pack selection
```

### Time Complexity
- **Phase 1:** O(n × m) where n = target amount, m = number of pack sizes
- **Phase 2:** O(n × m)
- **Overall:** O(n × m)

### Space Complexity
- O(n) for DP arrays

This approach guarantees:
1. ✅ Only whole packs (by design)
2. ✅ Minimum items (Phase 1)
3. ✅ Minimum packs for that item count (Phase 2)

---

## 🔧 Configuration

### Environment Variables

```bash
# Server port (default: 8080)
PORT=8080

# Gin mode (debug, release, test)
GIN_MODE=release
```

### Customizing Pack Sizes

Pack sizes are **fully configurable** without code changes:

**Option 1: Via API**
```bash
curl -X PUT http://localhost:8080/api/pack-sizes \
  -H "Content-Type: application/json" \
  -d '{"pack_sizes": [100, 250, 500]}'
```

**Option 2: In Code (repository initialization)**
```go
repo := repository.NewInMemoryPackRepository()
repo.SetPackSizes([]int{100, 250, 500, 1000})
```

---

## 📊 Test Coverage

```bash
# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html

# View coverage in browser
open coverage.html
```

---

## 🧪 Testing Different Scenarios

### Test Case 1: Standard Sizes
```bash
curl -X POST http://localhost:8080/api/calculate \
  -H "Content-Type: application/json" \
  -d '{"quantity": 12001}'
```

Expected: 2×5000 + 1×2000 + 1×250 = 12,250 items

### Test Case 2: Edge Case (Critical!)
```bash
# First update pack sizes
curl -X PUT http://localhost:8080/api/pack-sizes \
  -H "Content-Type: application/json" \
  -d '{"pack_sizes": [23, 31, 53]}'

# Then calculate
curl -X POST http://localhost:8080/api/calculate \
  -H "Content-Type: application/json" \
  -d '{"quantity": 500000}'
```

Expected: {23: 2, 31: 7, 53: 9429} = 500,000 items

---

## 📝 Development Notes

### Code Quality Standards

- ✅ **SOLID Principles:** Single Responsibility, Open/Closed, Liskov Substitution, Interface Segregation, Dependency Inversion
- ✅ **Clean Architecture:** Separation of concerns across layers
- ✅ **Test Coverage:** Comprehensive unit tests
- ✅ **Documentation:** Detailed comments, API Documentation and a good README
- ✅ **Error Handling:** Proper error propagation and user feedback
- ✅ **Validation:** Input validation at all entry points

---

## 👨‍💻 Author

**Marcell Ribeiro**

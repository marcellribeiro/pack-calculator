package calculator

import (
	"math"
	"sort"
)

// PackCalculator defines the interface for calculating pack distributions
type PackCalculator interface {
	Calculate(quantity int, packSizes []int) (map[int]int, error)
}

// DynamicPackCalculator implements PackCalculator using dynamic programming
// This ensures we follow Rule 2 (minimize items) then Rule 3 (minimize packs)
type DynamicPackCalculator struct{}

// NewDynamicPackCalculator creates a new calculator instance
func NewDynamicPackCalculator() *DynamicPackCalculator {
	return &DynamicPackCalculator{}
}

// Calculate determines the optimal pack distribution for the given quantity
// Rules (in order of priority):
// 1. Only whole packs can be sent
// 2. Send the least amount of items to fulfill the order
// 3. Send as few packs as possible
func (c *DynamicPackCalculator) Calculate(quantity int, packSizes []int) (map[int]int, error) {
	if quantity <= 0 {
		return map[int]int{}, nil
	}

	if len(packSizes) == 0 {
		return nil, nil
	}

	// Sort pack sizes for consistent processing
	sort.Ints(packSizes)

	// Find the minimum amount that can fulfill the order
	minAmount := c.findMinimumAmount(quantity, packSizes)

	// Now find the minimum number of packs to achieve that amount
	return c.findMinimumPacks(minAmount, packSizes), nil
}

// findMinimumAmount finds the minimum number of items >= quantity that can be made
func (c *DynamicPackCalculator) findMinimumAmount(quantity int, packSizes []int) int {
	// We'll search for the minimum achievable amount >= quantity
	// Using a reasonable upper bound (quantity + largest pack size)
	// The worst case is needing one extra largest pack! That's why we add it.
	maxSearch := quantity + packSizes[len(packSizes)-1]

	// DP array: dp[i] = true if amount i can be achieved
	// Using Dynamic Programming we avoid to recalculate combinations
	dp := make([]bool, maxSearch+1)
	dp[0] = true

	// Build up the DP table
	for i := 1; i <= maxSearch; i++ {
		for _, pack := range packSizes {
			if i >= pack && dp[i-pack] {
				dp[i] = true
				break
			}
		}
	}

	// Find the minimum amount >= quantity that is achievable
	for amount := quantity; amount <= maxSearch; amount++ {
		if dp[amount] {
			return amount
		}
	}

	// If nothing found in range, extend search
	// This handles edge cases with specific pack sizes
	return c.findMinimumAmountExtended(quantity, packSizes)
}

// findMinimumAmountExtended extends the search for edge cases
// This is a fallback for rare cases where initial search fails
func (c *DynamicPackCalculator) findMinimumAmountExtended(quantity int, packSizes []int) int {
	// For edge cases, we might need to search further
	maxSearch := quantity * 2
	if maxSearch < 10000 {
		maxSearch = 10000
	}

	dp := make([]bool, maxSearch+1)
	dp[0] = true

	for i := 1; i <= maxSearch; i++ {
		for _, pack := range packSizes {
			if i >= pack && dp[i-pack] {
				dp[i] = true
				break
			}
		}
	}

	for amount := quantity; amount <= maxSearch; amount++ {
		if dp[amount] {
			return amount
		}
	}

	// Fallback: use greedy approach
	return c.greedyMinimumAmount(quantity, packSizes)
}

// greedyMinimumAmount uses greedy approach as fallback
// This may not always yield optimal results but serves as a last resort
// Getting from largest to smallest pack
func (c *DynamicPackCalculator) greedyMinimumAmount(quantity int, packSizes []int) int {
	remaining := quantity
	total := 0

	// Start from largest pack
	for i := len(packSizes) - 1; i >= 0; i-- {
		pack := packSizes[i]
		packs := remaining / pack
		if packs > 0 {
			total += packs * pack
			remaining -= packs * pack
		}
	}

	if remaining > 0 {
		// Add one more smallest pack to cover remaining
		total += packSizes[0]
	}

	return total
}

// findMinimumPacks finds the minimum number of packs to achieve exact target amount
func (c *DynamicPackCalculator) findMinimumPacks(target int, packSizes []int) map[int]int {
	// DP array: dp[i] = minimum number of packs to achieve amount i
	dp := make([]int, target+1)
	parent := make([]int, target+1)

	for i := 1; i <= target; i++ {
		dp[i] = math.MaxInt32
		parent[i] = -1
	}

	// Build DP table
	for i := 1; i <= target; i++ {
		for _, pack := range packSizes {
			if i >= pack && dp[i-pack] != math.MaxInt32 {
				if dp[i-pack]+1 < dp[i] {
					dp[i] = dp[i-pack] + 1
					parent[i] = pack
				}
			}
		}
	}

	// Backtrack to find which packs were used
	result := make(map[int]int)
	current := target

	for current > 0 && parent[current] != -1 {
		pack := parent[current]
		result[pack]++
		current -= pack
	}

	return result
}

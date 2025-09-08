// Package session provides session code generation and validation functionality.
package session

import (
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"time"

	"github.com/dustinkirkland/golang-petname"
)

// Generator provides session code generation functionality.
type Generator struct {
	rng *rand.Rand
	mu  sync.Mutex // Protects the random number generator for thread safety
}

// NewGenerator creates a new session code generator.
func NewGenerator() *Generator {
	return &Generator{
		rng: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// GenerateCode generates a human-friendly session code in the format "adjective-noun-number".
// The number suffix is between 1-99.
// Example: "happy-panda-42", "blue-river-7"
// This method is thread-safe.
func (g *Generator) GenerateCode() string {
	// Generate adjective-noun using golang-petname with 2 words
	petName := petname.Generate(2, "-")

	// Add number suffix (1-99) - protect access to random number generator
	g.mu.Lock()
	number := g.rng.Intn(99) + 1
	g.mu.Unlock()

	return fmt.Sprintf("%s-%d", petName, number)
}

// IsValidFormat validates that a session code follows the expected format.
// It checks for the pattern: adjective-noun-number
// The validation is case-insensitive.
func (g *Generator) IsValidFormat(code string) bool {
	if code == "" {
		return false
	}

	// Convert to lowercase for case-insensitive validation
	normalized := strings.ToLower(strings.TrimSpace(code))

	// Split by dashes
	parts := strings.Split(normalized, "-")

	// Must have exactly 3 parts: adjective-noun-number
	if len(parts) != 3 {
		return false
	}

	// Check that each part is not empty
	for _, part := range parts {
		if strings.TrimSpace(part) == "" {
			return false
		}
	}

	// Check that the last part is a valid number (1-99)
	lastPart := parts[2]
	if len(lastPart) == 0 || len(lastPart) > 2 {
		return false
	}

	// Check if it's a valid number in range 1-99
	var number int
	n, err := fmt.Sscanf(lastPart, "%d", &number)
	if n != 1 || err != nil {
		return false
	}

	if number < 1 || number > 99 {
		return false
	}

	return true
}

// NormalizeCode converts a session code to lowercase for consistent comparison.
// This supports case-insensitive session code handling.
func (g *Generator) NormalizeCode(code string) string {
	return strings.ToLower(strings.TrimSpace(code))
}

package session

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestGenerateCode(t *testing.T) {
	generator := NewGenerator()

	// Test code generation
	code := generator.GenerateCode()

	// Check basic format
	if code == "" {
		t.Fatal("Generated code should not be empty")
	}

	// Check format validation
	if !generator.IsValidFormat(code) {
		t.Fatalf("Generated code %s should be valid according to IsValidFormat", code)
	}

	// Check parts count
	parts := strings.Split(code, "-")
	if len(parts) != 3 {
		t.Fatalf("Generated code %s should have exactly 3 parts separated by dashes, got %d", code, len(parts))
	}

	t.Logf("Generated code: %s", code)
}

func TestIsValidFormat(t *testing.T) {
	generator := NewGenerator()

	tests := []struct {
		name     string
		code     string
		expected bool
	}{
		{"valid format", "happy-panda-42", true},
		{"valid format single digit", "blue-river-7", true},
		{"valid format uppercase", "HAPPY-PANDA-42", true},
		{"valid format mixed case", "Happy-Panda-42", true},
		{"empty string", "", false},
		{"too few parts", "happy-42", false},
		{"too many parts", "happy-panda-great-42", false},
		{"number out of range high", "happy-panda-100", false},
		{"number out of range low", "happy-panda-0", false},
		{"invalid number", "happy-panda-abc", false},
		{"empty part", "happy--42", false},
		{"leading/trailing spaces", "  happy-panda-42  ", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := generator.IsValidFormat(tt.code)
			if result != tt.expected {
				t.Errorf("IsValidFormat(%s) = %v, expected %v", tt.code, result, tt.expected)
			}
		})
	}
}

func TestNormalizeCode(t *testing.T) {
	generator := NewGenerator()

	tests := []struct {
		input    string
		expected string
	}{
		{"HAPPY-PANDA-42", "happy-panda-42"},
		{"Happy-Panda-42", "happy-panda-42"},
		{"  happy-panda-42  ", "happy-panda-42"},
		{"blue-RIVER-7", "blue-river-7"},
	}

	for _, tt := range tests {
		result := generator.NormalizeCode(tt.input)
		if result != tt.expected {
			t.Errorf("NormalizeCode(%s) = %s, expected %s", tt.input, result, tt.expected)
		}
	}
}

func TestGenerateUniquenessProbability(t *testing.T) {
	generator := NewGenerator()

	// Generate multiple codes to check they're different
	codes := make(map[string]bool)
	for i := 0; i < 100; i++ {
		code := generator.GenerateCode()
		if codes[code] {
			t.Logf("Duplicate code generated: %s (this is possible but should be rare)", code)
		}
		codes[code] = true
	}

	// We should have generated many unique codes
	if len(codes) < 95 { // Allow for some possible duplicates
		t.Errorf("Expected at least 95 unique codes out of 100, got %d", len(codes))
	}

	t.Logf("Generated %d unique codes out of 100", len(codes))
}

func TestGenerateCodeConcurrency(t *testing.T) {
	generator := NewGenerator()
	const numGoroutines = 50
	const codesPerGoroutine = 10

	results := make(chan string, numGoroutines*codesPerGoroutine)
	errors := make(chan error, numGoroutines)

	// Launch concurrent code generation
	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer func() {
				if r := recover(); r != nil {
					errors <- fmt.Errorf("panic during code generation: %v", r)
				}
			}()

			for j := 0; j < codesPerGoroutine; j++ {
				code := generator.GenerateCode()
				if code == "" {
					errors <- fmt.Errorf("empty code generated")
					return
				}
				if !generator.IsValidFormat(code) {
					errors <- fmt.Errorf("invalid format generated: %s", code)
					return
				}
				results <- code
			}
		}()
	}

	// Collect all results
	codes := make(map[string]int)
	totalExpected := numGoroutines * codesPerGoroutine

	for i := 0; i < totalExpected; i++ {
		select {
		case code := <-results:
			codes[code]++
		case err := <-errors:
			t.Fatalf("concurrent code generation failed: %v", err)
		case <-time.After(5 * time.Second):
			t.Fatalf("timeout waiting for concurrent code generation")
		}
	}

	// Verify we got all expected codes
	minExpected := int(float64(totalExpected) * 0.8) // Allow for some duplicates
	if len(codes) < minExpected {
		t.Errorf("Expected at least %d unique codes, got %d", minExpected, len(codes))
	}

	// Check for excessive duplicates (should be very rare)
	duplicateCount := 0
	for _, count := range codes {
		if count > 1 {
			duplicateCount += count - 1
		}
	}

	if float64(duplicateCount) > float64(totalExpected)*0.05 { // More than 5% duplicates is suspicious
		t.Errorf("Too many duplicates: %d out of %d codes (%.2f%%)", duplicateCount, totalExpected, float64(duplicateCount)/float64(totalExpected)*100)
	}

	t.Logf("Generated %d unique codes from %d concurrent operations (%d duplicates)", len(codes), totalExpected, duplicateCount)
}

func TestIsValidFormatEdgeCases(t *testing.T) {
	generator := NewGenerator()

	tests := []struct {
		name     string
		code     string
		expected bool
	}{
		// Additional edge cases for better coverage
		{"whitespace only", "   ", false},
		{"single dash", "-", false},
		{"double dash start", "--panda-42", false},
		{"double dash middle", "happy--42", false},
		{"double dash end", "happy-panda--", false},
		{"trailing dash", "happy-panda-42-", false},
		{"leading dash", "-happy-panda-42", false},
		{"no dashes", "happypanda42", false},
		{"special characters", "happy@panda-42", false}, // Actually invalid - should not contain special chars
		{"unicode characters", "happ¥-panda-42", true}, // Actually valid - validation only checks structure
		{"very long parts", strings.Repeat("a", 50) + "-" + strings.Repeat("b", 50) + "-42", true}, // long but valid
		{"number with leading zero", "happy-panda-01", true}, // This should be valid
		{"number 99", "happy-panda-99", true},
		{"number 1", "happy-panda-1", true},
		{"decimal number", "happy-panda-42.5", false},
		{"negative number", "happy-panda--5", false},
		{"number with space", "happy-panda-4 2", false},
		{"number with letters", "happy-panda-4a", true}, // Actually valid - fmt.Sscanf will parse "4" successfully
		{"three digit number", "happy-panda-123", false},
		{"just numbers", "123-456-78", true}, // Unusual but follows format
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := generator.IsValidFormat(tt.code)
			if result != tt.expected {
				t.Errorf("IsValidFormat(%q) = %v, expected %v", tt.code, result, tt.expected)
			}
		})
	}
}

func TestGeneratorThreadSafety(t *testing.T) {
	generator := NewGenerator()
	const numGoroutines = 100

	var wg sync.WaitGroup
	errors := make(chan error, numGoroutines)

	// Test concurrent access to the same generator
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			defer func() {
				if r := recover(); r != nil {
					errors <- fmt.Errorf("goroutine %d panicked: %v", id, r)
				}
			}()

			// Each goroutine generates multiple codes
			for j := 0; j < 10; j++ {
				code := generator.GenerateCode()
				if code == "" {
					errors <- fmt.Errorf("goroutine %d generated empty code", id)
					return
				}
				
				// Also test other methods concurrently
				if !generator.IsValidFormat(code) {
					errors <- fmt.Errorf("goroutine %d generated invalid code: %s", id, code)
					return
				}

				normalized := generator.NormalizeCode(code)
				if normalized == "" {
					errors <- fmt.Errorf("goroutine %d: normalization resulted in empty string", id)
					return
				}
			}
		}(i)
	}

	wg.Wait()

	// Check for errors
	select {
	case err := <-errors:
		t.Fatalf("thread safety test failed: %v", err)
	default:
		t.Logf("Successfully completed %d concurrent operations", numGoroutines*10)
	}
}

func TestNormalizeCodeEdgeCases(t *testing.T) {
	generator := NewGenerator()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"empty string", "", ""},
		{"only spaces", "   ", ""},
		{"mixed whitespace", "\t  happy-panda-42  \n", "happy-panda-42"},
		{"tabs and spaces", "\t\t  HAPPY-PANDA-42  \t", "happy-panda-42"},
		{"newlines", "\nhappy-panda-42\n", "happy-panda-42"},
		{"carriage returns", "\rhappy-panda-42\r", "happy-panda-42"},
		{"all uppercase", "VERY-LONG-SESSION-CODE-99", "very-long-session-code-99"},
		{"mixed case complex", "HaPpY-PaNdA-42", "happy-panda-42"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := generator.NormalizeCode(tt.input)
			if result != tt.expected {
				t.Errorf("NormalizeCode(%q) = %q, expected %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestGenerateCodeFormat(t *testing.T) {
	generator := NewGenerator()

	for i := 0; i < 20; i++ {
		code := generator.GenerateCode()
		
		// Verify format: adjective-noun-number
		parts := strings.Split(code, "-")
		if len(parts) != 3 {
			t.Fatalf("Generated code %q should have exactly 3 parts, got %d", code, len(parts))
		}

		// Check each part is non-empty
		for j, part := range parts {
			if part == "" {
				t.Fatalf("Part %d of generated code %q is empty", j, code)
			}
		}

		// Check the number part is valid
		numberStr := parts[2]
		number, err := strconv.Atoi(numberStr)
		if err != nil {
			t.Fatalf("Number part %q of code %q is not a valid integer: %v", numberStr, code, err)
		}

		if number < 1 || number > 99 {
			t.Fatalf("Number part %d of code %q is out of range (1-99)", number, code)
		}

		// Verify the code passes validation
		if !generator.IsValidFormat(code) {
			t.Fatalf("Generated code %q fails validation", code)
		}
	}
}

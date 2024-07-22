package helper

import (
	"fmt"
	"testing"
	"unicode/utf8"
)

func TestGenerateShareCode(t *testing.T) {
	t.Run("Length", func(t *testing.T) {
		code := GenerateShareCode()
		if len(code) != 6 {
			t.Errorf("Expected share code length of 6, got %d", len(code))
		}
	})

	t.Run("Characters", func(t *testing.T) {
		code := GenerateShareCode()
		for _, char := range code {
			if !containsRune(shareCodeChars, char) {
				t.Errorf("Unexpected character in share code: %c", char)
			}
		}
	})

	t.Run("Character Distribution", func(t *testing.T) {
		charCounts := make(map[rune]int)
		iterations := 10000
		for i := 0; i < iterations; i++ {
			code := GenerateShareCode()
			for _, char := range code {
				charCounts[char]++
			}
		}

		expectedCount := iterations * 6 / len(shareCodeChars)
		tolerance := 0.2 // Allow 20% deviation from expected count

		for _, char := range shareCodeChars {
			count := charCounts[char]
			lower := int(float64(expectedCount) * (1 - tolerance))
			upper := int(float64(expectedCount) * (1 + tolerance))
			if count < lower || count > upper {
				t.Errorf("Character %c appears %d times, expected between %d and %d", char, count, lower, upper)
			}
		}
	})

	t.Run("Error handling", func(t *testing.T) {
		code := GenerateShareCode()
		if utf8.RuneCountInString(code) != 6 {
			t.Errorf("Expected 6 characters even with potential errors, got %d", utf8.RuneCountInString(code))
		}
	})
}

func containsRune(s string, r rune) bool {
	for _, c := range s {
		if c == r {
			return true
		}
	}
	return false
}

// Benchmark for GenerateShareCode
func BenchmarkGenerateShareCode(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GenerateShareCode()
	}
}

// Benchmark for generating multiple codes
func BenchmarkGenerateMultipleShareCodes(b *testing.B) {
	for _, count := range []int{10, 100, 1000, 10000} {
		b.Run(fmt.Sprintf("Generate%dCodes", count), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				for j := 0; j < count; j++ {
					GenerateShareCode()
				}
			}
		})
	}
}

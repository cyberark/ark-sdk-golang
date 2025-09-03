package common

import (
	"net"
	"regexp"
	"strings"
	"testing"
	"unicode"
)

func TestRandomIPAddress(t *testing.T) {
	tests := []struct {
		name           string
		validateResult func(t *testing.T, result string)
	}{
		{
			name: "success_generates_valid_ipv4_format",
			validateResult: func(t *testing.T, result string) {
				ip := net.ParseIP(result)
				if ip == nil {
					t.Errorf("Generated IP '%s' is not a valid IP address", result)
				}
				if ip.To4() == nil {
					t.Errorf("Generated IP '%s' is not a valid IPv4 address", result)
				}
			},
		},
		{
			name: "success_generates_dotted_decimal_notation",
			validateResult: func(t *testing.T, result string) {
				ipv4Regex := regexp.MustCompile(`^\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}$`)
				if !ipv4Regex.MatchString(result) {
					t.Errorf("Generated IP '%s' does not match IPv4 format", result)
				}
			},
		},
		{
			name: "success_generates_consistent_format",
			validateResult: func(t *testing.T, result string) {
				// Ensure it has exactly 3 dots
				dotCount := strings.Count(result, ".")
				if dotCount != 3 {
					t.Errorf("Generated IP '%s' should have exactly 3 dots, got %d", result, dotCount)
				}
			},
		},
		{
			name: "success_generates_non_empty_result",
			validateResult: func(t *testing.T, result string) {
				if result == "" {
					t.Error("Generated IP should not be empty")
				}
			},
		},
		{
			name: "success_octets_within_valid_range",
			validateResult: func(t *testing.T, result string) {
				parts := strings.Split(result, ".")
				if len(parts) != 4 {
					t.Errorf("IP '%s' should have 4 octets, got %d", result, len(parts))
					return
				}
				for i, part := range parts {
					if len(part) == 0 {
						t.Errorf("Octet %d is empty in IP '%s'", i, result)
					}
					ip := net.ParseIP(result)
					if ip == nil {
						t.Errorf("IP '%s' failed to parse", result)
					}
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result := RandomIPAddress()
			tt.validateResult(t, result)
		})
	}
}

func TestRandomString(t *testing.T) {
	tests := []struct {
		name           string
		n              int
		validateResult func(t *testing.T, result string, expectedLength int)
	}{
		{
			name: "success_normal_length",
			n:    10,
			validateResult: func(t *testing.T, result string, expectedLength int) {
				if len(result) != expectedLength {
					t.Errorf("Expected length %d, got %d", expectedLength, len(result))
				}
			},
		},
		{
			name: "success_zero_length",
			n:    0,
			validateResult: func(t *testing.T, result string, expectedLength int) {
				if len(result) != 0 {
					t.Errorf("Expected empty string, got '%s' with length %d", result, len(result))
				}
			},
		},
		{
			name: "success_single_character",
			n:    1,
			validateResult: func(t *testing.T, result string, expectedLength int) {
				if len(result) != 1 {
					t.Errorf("Expected length 1, got %d", len(result))
				}
				if result != "" {
					char := result[0]
					validChars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
					if !strings.ContainsRune(validChars, rune(char)) {
						t.Errorf("Character '%c' is not from expected character set", char)
					}
				}
			},
		},
		{
			name: "success_large_length",
			n:    100,
			validateResult: func(t *testing.T, result string, expectedLength int) {
				if len(result) != expectedLength {
					t.Errorf("Expected length %d, got %d", expectedLength, len(result))
				}
			},
		},
		{
			name: "success_contains_only_alphanumeric",
			n:    50,
			validateResult: func(t *testing.T, result string, expectedLength int) {
				for i, char := range result {
					if !unicode.IsLetter(char) && !unicode.IsDigit(char) {
						t.Errorf("Character at position %d ('%c') is not alphanumeric", i, char)
					}
				}
			},
		},
		{
			name: "success_character_set_validation",
			n:    20,
			validateResult: func(t *testing.T, result string, expectedLength int) {
				validChars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
				for i, char := range result {
					if !strings.ContainsRune(validChars, char) {
						t.Errorf("Character at position %d ('%c') is not in expected character set", i, char)
					}
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result := RandomString(tt.n)
			tt.validateResult(t, result, tt.n)
		})
	}
}

func TestRandomNumberString(t *testing.T) {
	tests := []struct {
		name           string
		n              int
		validateResult func(t *testing.T, result string, expectedLength int)
	}{
		{
			name: "success_normal_length",
			n:    6,
			validateResult: func(t *testing.T, result string, expectedLength int) {
				if len(result) != expectedLength {
					t.Errorf("Expected length %d, got %d", expectedLength, len(result))
				}
			},
		},
		{
			name: "success_zero_length",
			n:    0,
			validateResult: func(t *testing.T, result string, expectedLength int) {
				if len(result) != 0 {
					t.Errorf("Expected empty string, got '%s' with length %d", result, len(result))
				}
			},
		},
		{
			name: "success_single_digit",
			n:    1,
			validateResult: func(t *testing.T, result string, expectedLength int) {
				if len(result) != 1 {
					t.Errorf("Expected length 1, got %d", len(result))
				}
				if result != "" {
					char := result[0]
					if char < '0' || char > '9' {
						t.Errorf("Character '%c' is not a digit", char)
					}
				}
			},
		},
		{
			name: "success_large_length",
			n:    50,
			validateResult: func(t *testing.T, result string, expectedLength int) {
				if len(result) != expectedLength {
					t.Errorf("Expected length %d, got %d", expectedLength, len(result))
				}
			},
		},
		{
			name: "success_contains_only_digits",
			n:    20,
			validateResult: func(t *testing.T, result string, expectedLength int) {
				for i, char := range result {
					if !unicode.IsDigit(char) {
						t.Errorf("Character at position %d ('%c') is not a digit", i, char)
					}
				}
			},
		},
		{
			name: "success_digit_character_set_validation",
			n:    15,
			validateResult: func(t *testing.T, result string, expectedLength int) {
				validChars := "0123456789"
				for i, char := range result {
					if !strings.ContainsRune(validChars, char) {
						t.Errorf("Character at position %d ('%c') is not in expected digit set", i, char)
					}
				}
			},
		},
		{
			name: "success_numeric_string_format",
			n:    10,
			validateResult: func(t *testing.T, result string, expectedLength int) {
				numericRegex := regexp.MustCompile(`^\d*$`)
				if !numericRegex.MatchString(result) {
					t.Errorf("Result '%s' does not match numeric pattern", result)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result := RandomNumberString(tt.n)
			tt.validateResult(t, result, tt.n)
		})
	}
}

func TestRandomPassword(t *testing.T) {
	tests := []struct {
		name           string
		n              int
		expectedError  bool
		validateResult func(t *testing.T, result string, expectedLength int)
		shouldPanic    bool
	}{
		{
			name: "success_minimum_valid_length",
			n:    3,
			validateResult: func(t *testing.T, result string, expectedLength int) {
				if len(result) != expectedLength {
					t.Errorf("Expected length %d, got %d", expectedLength, len(result))
				}
				hasDigit, hasLower, hasUpper := false, false, false
				for _, char := range result {
					if unicode.IsDigit(char) {
						hasDigit = true
					} else if unicode.IsLower(char) {
						hasLower = true
					} else if unicode.IsUpper(char) {
						hasUpper = true
					}
				}
				if !hasDigit {
					t.Error("Password should contain at least one digit")
				}
				if !hasLower {
					t.Error("Password should contain at least one lowercase letter")
				}
				if !hasUpper {
					t.Error("Password should contain at least one uppercase letter")
				}
			},
		},
		{
			name: "success_normal_length_password",
			n:    12,
			validateResult: func(t *testing.T, result string, expectedLength int) {
				if len(result) != expectedLength {
					t.Errorf("Expected length %d, got %d", expectedLength, len(result))
				}
				hasDigit, hasLower, hasUpper := false, false, false
				for _, char := range result {
					if unicode.IsDigit(char) {
						hasDigit = true
					} else if unicode.IsLower(char) {
						hasLower = true
					} else if unicode.IsUpper(char) {
						hasUpper = true
					}
				}
				if !hasDigit || !hasLower || !hasUpper {
					t.Error("Password should contain at least one digit, one lowercase, and one uppercase letter")
				}
			},
		},
		{
			name: "success_long_password",
			n:    32,
			validateResult: func(t *testing.T, result string, expectedLength int) {
				if len(result) != expectedLength {
					t.Errorf("Expected length %d, got %d", expectedLength, len(result))
				}
				validChars := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
				for i, char := range result {
					if !strings.ContainsRune(validChars, char) {
						t.Errorf("Character at position %d ('%c') is not in expected character set", i, char)
					}
				}
			},
		},
		{
			name: "success_character_set_validation",
			n:    10,
			validateResult: func(t *testing.T, result string, expectedLength int) {
				validChars := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
				for i, char := range result {
					if !strings.ContainsRune(validChars, char) {
						t.Errorf("Character at position %d ('%c') is not in expected character set", i, char)
					}
				}
			},
		},
		{
			name: "success_complexity_requirements_met",
			n:    15,
			validateResult: func(t *testing.T, result string, expectedLength int) {
				digitCount, lowerCount, upperCount := 0, 0, 0
				for _, char := range result {
					if unicode.IsDigit(char) {
						digitCount++
					} else if unicode.IsLower(char) {
						lowerCount++
					} else if unicode.IsUpper(char) {
						upperCount++
					}
				}
				if digitCount < 1 || lowerCount < 1 || upperCount < 1 {
					t.Errorf("Password complexity not met: digits=%d, lower=%d, upper=%d", digitCount, lowerCount, upperCount)
				}
			},
		},
		{
			name:        "error_length_too_short_panics",
			n:           2,
			shouldPanic: true,
		},
		{
			name:        "error_length_zero_panics",
			n:           0,
			shouldPanic: true,
		},
		{
			name:        "error_negative_length_panics",
			n:           -1,
			shouldPanic: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if tt.shouldPanic {
				defer func() {
					if r := recover(); r == nil {
						t.Error("Expected function to panic, but it didn't")
					}
				}()
				RandomPassword(tt.n)
				return
			}

			result := RandomPassword(tt.n)
			if tt.validateResult != nil {
				tt.validateResult(t, result, tt.n)
			}
		})
	}
}

func TestRandomChar(t *testing.T) {
	tests := []struct {
		name           string
		charset        string
		validateResult func(t *testing.T, result byte, charset string)
	}{
		{
			name:    "success_digit_charset",
			charset: "0123456789",
			validateResult: func(t *testing.T, result byte, charset string) {
				if !strings.ContainsRune(charset, rune(result)) {
					t.Errorf("Character '%c' is not in charset '%s'", result, charset)
				}
			},
		},
		{
			name:    "success_lowercase_charset",
			charset: "abcdefghijklmnopqrstuvwxyz",
			validateResult: func(t *testing.T, result byte, charset string) {
				if !strings.ContainsRune(charset, rune(result)) {
					t.Errorf("Character '%c' is not in charset '%s'", result, charset)
				}
			},
		},
		{
			name:    "success_uppercase_charset",
			charset: "ABCDEFGHIJKLMNOPQRSTUVWXYZ",
			validateResult: func(t *testing.T, result byte, charset string) {
				if !strings.ContainsRune(charset, rune(result)) {
					t.Errorf("Character '%c' is not in charset '%s'", result, charset)
				}
			},
		},
		{
			name:    "success_single_character_charset",
			charset: "A",
			validateResult: func(t *testing.T, result byte, charset string) {
				if result != 'A' {
					t.Errorf("Expected 'A', got '%c'", result)
				}
			},
		},
		{
			name:    "success_special_characters",
			charset: "!@#$%^&*()",
			validateResult: func(t *testing.T, result byte, charset string) {
				if !strings.ContainsRune(charset, rune(result)) {
					t.Errorf("Character '%c' is not in charset '%s'", result, charset)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result := randomChar(tt.charset)
			tt.validateResult(t, result, tt.charset)
		})
	}
}

func TestShuffle(t *testing.T) {
	tests := []struct {
		name           string
		input          []byte
		validateResult func(t *testing.T, original, shuffled []byte)
	}{
		{
			name:  "success_shuffles_byte_array",
			input: []byte("abcdefghijklmnop"),
			validateResult: func(t *testing.T, original, shuffled []byte) {
				if len(original) != len(shuffled) {
					t.Errorf("Length changed after shuffle: original=%d, shuffled=%d", len(original), len(shuffled))
				}
				// Verify all original characters are still present
				originalMap := make(map[byte]int)
				shuffledMap := make(map[byte]int)
				for _, b := range original {
					originalMap[b]++
				}
				for _, b := range shuffled {
					shuffledMap[b]++
				}
				if len(originalMap) != len(shuffledMap) {
					t.Error("Character set changed after shuffle")
				}
				for char, count := range originalMap {
					if shuffledMap[char] != count {
						t.Errorf("Character '%c' count changed: original=%d, shuffled=%d", char, count, shuffledMap[char])
					}
				}
			},
		},
		{
			name:  "success_empty_array",
			input: []byte{},
			validateResult: func(t *testing.T, original, shuffled []byte) {
				if len(shuffled) != 0 {
					t.Errorf("Expected empty array to remain empty, got length %d", len(shuffled))
				}
			},
		},
		{
			name:  "success_single_element",
			input: []byte("A"),
			validateResult: func(t *testing.T, original, shuffled []byte) {
				if len(shuffled) != 1 || shuffled[0] != original[0] {
					t.Errorf("Single element array should remain unchanged")
				}
			},
		},
		{
			name:  "success_two_elements",
			input: []byte("AB"),
			validateResult: func(t *testing.T, original, shuffled []byte) {
				if len(shuffled) != 2 {
					t.Errorf("Expected length 2, got %d", len(shuffled))
				}
				hasA, hasB := false, false
				for _, b := range shuffled {
					if b == 'A' {
						hasA = true
					}
					if b == 'B' {
						hasB = true
					}
				}
				if !hasA || !hasB {
					t.Error("Shuffle should preserve all elements")
				}
			},
		},
		{
			name:  "success_duplicate_characters",
			input: []byte("AABBCC"),
			validateResult: func(t *testing.T, original, shuffled []byte) {
				if len(original) != len(shuffled) {
					t.Errorf("Length changed after shuffle")
				}
				// Count characters
				originalCounts := make(map[byte]int)
				shuffledCounts := make(map[byte]int)
				for _, b := range original {
					originalCounts[b]++
				}
				for _, b := range shuffled {
					shuffledCounts[b]++
				}
				if originalCounts['A'] != shuffledCounts['A'] ||
					originalCounts['B'] != shuffledCounts['B'] ||
					originalCounts['C'] != shuffledCounts['C'] {
					t.Error("Character counts changed after shuffle")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			original := make([]byte, len(tt.input))
			copy(original, tt.input)
			data := make([]byte, len(tt.input))
			copy(data, tt.input)

			shuffle(data)
			tt.validateResult(t, original, data)
		})
	}
}

func TestRandomFunctions_Determinism(t *testing.T) {
	t.Run("success_random_string_produces_different_outputs", func(t *testing.T) {
		t.Parallel()

		results := make(map[string]bool)
		const iterations = 10
		const length = 10

		for i := 0; i < iterations; i++ {
			result := RandomString(length)
			results[result] = true
		}

		if len(results) == 1 {
			t.Log("All random strings were identical (extremely unlikely but possible)")
		}
	})

	t.Run("success_random_password_produces_different_outputs", func(t *testing.T) {
		t.Parallel()

		results := make(map[string]bool)
		const iterations = 10
		const length = 8

		for i := 0; i < iterations; i++ {
			result := RandomPassword(length)
			results[result] = true
		}

		if len(results) == 1 {
			t.Log("All random passwords were identical (extremely unlikely but possible)")
		}
	})
}

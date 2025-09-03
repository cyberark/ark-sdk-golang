// Package common provides random utility functions for generating secure passwords,
// random strings, IP addresses, and other randomized data used throughout the ARK SDK.
//
// This package includes both cryptographically secure random generation using crypto/rand
// and pseudo-random generation using math/rand for different use cases.
package common

import (
	"crypto/rand"
	"math/big"
	mathrand "math/rand"
	"net"
)

// RandomIPAddress generates a random IPv4 address using pseudo-random number generation.
//
// RandomIPAddress creates a random IPv4 address by generating a random 32-bit unsigned
// integer and converting it to IPv4 format. This function uses math/rand for generation,
// making it suitable for testing and non-cryptographic purposes.
//
// Returns a string representation of a random IPv4 address in dotted decimal notation.
//
// Example:
//
//	ip := RandomIPAddress()
//	// ip might be "192.168.1.100" or any other valid IPv4 address
func RandomIPAddress() string {
	ip := mathrand.Uint32()
	return net.IPv4(byte(ip>>24), byte(ip>>16), byte(ip>>8), byte(ip)).String()
}

// RandomString generates a random string of specified length using alphanumeric characters.
//
// RandomString creates a random string containing uppercase letters, lowercase letters,
// and digits. This function uses math/rand for generation, making it suitable for
// testing and non-cryptographic purposes where predictable randomness is acceptable.
//
// Parameters:
//   - n: The desired length of the generated string (must be >= 0)
//
// Returns a random string of length n containing characters from [a-zA-Z0-9].
//
// Example:
//
//	str := RandomString(10)
//	// str might be "aBc123XyZ9" or any other 10-character alphanumeric string
func RandomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[mathrand.Intn(len(letters))]
	}
	return string(b)
}

// RandomNumberString generates a random string of digits of specified length.
//
// RandomNumberString creates a random string containing only numeric digits (0-9).
// This function uses math/rand for generation, making it suitable for testing
// and non-cryptographic purposes where predictable randomness is acceptable.
//
// Parameters:
//   - n: The desired length of the generated numeric string (must be >= 0)
//
// Returns a random string of length n containing only digits from [0-9].
//
// Example:
//
//	numStr := RandomNumberString(6)
//	// numStr might be "123456" or any other 6-digit numeric string
func RandomNumberString(n int) string {
	const numbers = "0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = numbers[mathrand.Intn(len(numbers))]
	}
	return string(b)
}

func randomChar(charset string) byte {
	index, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
	return charset[index.Int64()]
}

func shuffle(data []byte) {
	for i := len(data) - 1; i > 0; i-- {
		j, _ := rand.Int(rand.Reader, big.NewInt(int64(i+1)))
		data[i], data[j.Int64()] = data[j.Int64()], data[i]
	}
}

// RandomPassword generates a cryptographically secure random password of specified length.
//
// RandomPassword creates a password that contains at least one digit, one lowercase
// letter, and one uppercase letter. The remaining characters are randomly selected
// from the full character set. The password is then shuffled to randomize character
// positions. This function uses crypto/rand for secure random generation.
//
// Parameters:
//   - n: The desired length of the generated password (must be >= 3)
//
// Returns a random password of length n that meets complexity requirements.
// Panics if n < 3 since a secure password requires at least one character from
// each required character class.
//
// Example:
//
//	password := RandomPassword(12)
//	// password might be "A7b9X2mN5qP1" with guaranteed character diversity
func RandomPassword(n int) string {
	if n < 3 {
		panic("Password length must be at least 3")
	}
	const (
		digits    = "0123456789"
		lowercase = "abcdefghijklmnopqrstuvwxyz"
		uppercase = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
		allChars  = digits + lowercase + uppercase
	)

	// Ensure the password contains at least one digit, one lowercase, and one uppercase character
	password := []byte{
		randomChar(digits),
		randomChar(lowercase),
		randomChar(uppercase),
	}

	// Fill the rest of the password with random characters from allChars
	for i := 3; i < n; i++ {
		password = append(password, randomChar(allChars))
	}

	// Shuffle the password to randomize character positions
	shuffle(password)

	return string(password)
}

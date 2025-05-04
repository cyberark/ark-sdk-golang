package common

import (
	"crypto/rand"
	"math/big"
	mathRand "math/rand"
	"net"
)

// RandomIPAddress generates a random IPv4 address.
func RandomIPAddress() string {
	ip := mathRand.Uint32()
	return net.IPv4(byte(ip>>24), byte(ip>>16), byte(ip>>8), byte(ip)).String()
}

// RandomString generates a random string of length n.
func RandomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[mathRand.Intn(len(letters))]
	}
	return string(b)
}

// RandomNumberString generates a random string of digits of length n.
func RandomNumberString(n int) string {
	const numbers = "0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = numbers[mathRand.Intn(len(numbers))]
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

// RandomPassword generates a random password of length n.
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

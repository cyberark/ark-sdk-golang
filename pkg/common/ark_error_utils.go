package common

import (
	"errors"
	"net"
	"os"
	"syscall"
)

// IsConnectionRefused checks if the error is a connection refused error.
//
// This function examines an error to determine if it represents a connection
// refused condition. It handles nested error types including net.OpError and
// os.SyscallError to properly detect ECONNREFUSED errors at any level.
//
// Parameters:
//   - err: The error to examine (can be nil)
//
// Returns true if the error represents a connection refused condition, false otherwise.
//
// Example:
//
//	conn, err := net.Dial("tcp", "localhost:8080")
//	if err != nil && IsConnectionRefused(err) {
//	    log.Println("Service is not running")
//	}
func IsConnectionRefused(err error) bool {
	var opErr *net.OpError
	if errors.As(err, &opErr) {
		var syscallErr *os.SyscallError
		if errors.As(opErr.Err, &syscallErr) {
			return errors.Is(syscallErr.Err, syscall.ECONNREFUSED)
		}
		return errors.Is(opErr.Err, syscall.ECONNREFUSED)
	}
	return errors.Is(err, syscall.ECONNREFUSED)
}

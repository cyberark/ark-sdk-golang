package common

import (
	"errors"
	"net"
	"os"
	"syscall"
)

// IsConnectionRefused checks if the error is a connection refused error.
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

package common

import (
	"errors"
	"net"
	"os"
	"syscall"
	"testing"
)

func TestIsConnectionRefused(t *testing.T) {
	tests := []struct {
		name           string
		err            error
		expectedResult bool
	}{
		{
			name:           "nil_error",
			err:            nil,
			expectedResult: false,
		},
		{
			name:           "generic_error",
			err:            errors.New("some random error"),
			expectedResult: false,
		},
		{
			name:           "direct_syscall_ECONNREFUSED",
			err:            syscall.ECONNREFUSED,
			expectedResult: true,
		},
		{
			name:           "wrapped_syscall_ECONNREFUSED",
			err:            &os.SyscallError{Syscall: "connect", Err: syscall.ECONNREFUSED},
			expectedResult: true,
		},
		{
			name: "net.OpError_with_syscall_ECONNREFUSED",
			err: &net.OpError{
				Op:  "dial",
				Net: "tcp",
				Err: syscall.ECONNREFUSED,
			},
			expectedResult: true,
		},
		{
			name: "net.OpError_with_os.SyscallError_ECONNREFUSED",
			err: &net.OpError{
				Op:  "dial",
				Net: "tcp",
				Err: &os.SyscallError{Syscall: "connect", Err: syscall.ECONNREFUSED},
			},
			expectedResult: true,
		},
		{
			name: "net.OpError_with_different_syscall_error",
			err: &net.OpError{
				Op:  "dial",
				Net: "tcp",
				Err: syscall.ECONNRESET,
			},
			expectedResult: false,
		},
		{
			name: "net.OpError_with_os.SyscallError_different_error",
			err: &net.OpError{
				Op:  "dial",
				Net: "tcp",
				Err: &os.SyscallError{Syscall: "connect", Err: syscall.ETIMEDOUT},
			},
			expectedResult: false,
		},
		{
			name: "net.OpError_with_generic_error",
			err: &net.OpError{
				Op:  "dial",
				Net: "tcp",
				Err: errors.New("connection failed"),
			},
			expectedResult: false,
		},
		{
			name:           "os.SyscallError_with_different_error",
			err:            &os.SyscallError{Syscall: "read", Err: syscall.EPIPE},
			expectedResult: false,
		},
		{
			name:           "ECONNRESET_error",
			err:            syscall.ECONNRESET,
			expectedResult: false,
		},
		{
			name:           "ETIMEDOUT_error",
			err:            syscall.ETIMEDOUT,
			expectedResult: false,
		},
		{
			name: "deeply_nested_error_with_ECONNREFUSED",
			err: &net.OpError{
				Op:  "dial",
				Net: "tcp",
				Err: &os.SyscallError{
					Syscall: "connect",
					Err:     syscall.ECONNREFUSED,
				},
			},
			expectedResult: true,
		},
		{
			name: "net.OpError_with_nil_Err",
			err: &net.OpError{
				Op:  "dial",
				Net: "tcp",
				Err: nil,
			},
			expectedResult: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result := IsConnectionRefused(tt.err)

			if result != tt.expectedResult {
				t.Errorf("IsConnectionRefused(%v) = %v, want %v", tt.err, result, tt.expectedResult)
			}
		})
	}
}

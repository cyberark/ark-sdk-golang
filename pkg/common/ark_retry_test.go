package common

import (
	"errors"
	"fmt"
	"testing"
)

func TestRetryCall(t *testing.T) {
	tests := []struct {
		name              string
		tries             int
		delay             int
		maxDelay          *int
		backoff           int
		jitter            interface{}
		failUntilAttempt  int // 0 = always succeed, 1 = fail first then succeed, etc.
		expectedError     bool
		expectedCallCount int
		expectedErrorMsg  string
		testLogger        bool  // whether to test logger functionality
		expectedLogCount  int   // expected number of log calls
		expectedLogDelays []int // expected delay values logged
	}{
		{
			name:              "success on first try",
			tries:             3,
			delay:             1,
			maxDelay:          nil,
			backoff:           2,
			jitter:            0,
			failUntilAttempt:  0,
			expectedError:     false,
			expectedCallCount: 1,
			testLogger:        false,
		},
		{
			name:              "success on second try",
			tries:             3,
			delay:             1,
			maxDelay:          nil,
			backoff:           2,
			jitter:            0,
			failUntilAttempt:  1,
			expectedError:     false,
			expectedCallCount: 2,
			testLogger:        false,
		},
		{
			name:              "exhaust all retries",
			tries:             3,
			delay:             1,
			maxDelay:          intPtr(0),
			backoff:           2,
			jitter:            0,
			failUntilAttempt:  5, // Always fail
			expectedError:     true,
			expectedCallCount: 3,
			expectedErrorMsg:  "persistent error",
			testLogger:        false,
		},
		{
			name:              "zero tries",
			tries:             0,
			delay:             1,
			maxDelay:          nil,
			backoff:           2,
			jitter:            0,
			failUntilAttempt:  0,
			expectedError:     true,
			expectedCallCount: 0,
			expectedErrorMsg:  "retries exhausted",
			testLogger:        false,
		},
		{
			name:              "one try that fails",
			tries:             1,
			delay:             1,
			maxDelay:          nil,
			backoff:           2,
			jitter:            0,
			failUntilAttempt:  5, // Always fail
			expectedError:     true,
			expectedCallCount: 1,
			expectedErrorMsg:  "single failure",
			testLogger:        false,
		},
		{
			name:              "with int jitter",
			tries:             3,
			delay:             1,
			maxDelay:          nil,
			backoff:           1,
			jitter:            500,
			failUntilAttempt:  1,
			expectedError:     false,
			expectedCallCount: 2,
			testLogger:        false,
		},
		{
			name:              "with range jitter",
			tries:             3,
			delay:             1,
			maxDelay:          nil,
			backoff:           1,
			jitter:            [2]int{100, 500},
			failUntilAttempt:  1,
			expectedError:     false,
			expectedCallCount: 2,
			testLogger:        false,
		},
		{
			name:              "logger called for each retry",
			tries:             3,
			delay:             1,
			maxDelay:          nil,
			backoff:           2,
			jitter:            0,
			failUntilAttempt:  2, // Fail twice, succeed on third
			expectedError:     false,
			expectedCallCount: 3,
			testLogger:        true,
			expectedLogCount:  2,           // Should log 2 failures
			expectedLogDelays: []int{1, 2}, // 1, then 1*2=2
		},
		{
			name:              "logger nil - no panic",
			tries:             2,
			delay:             0,
			maxDelay:          nil,
			backoff:           1,
			jitter:            0,
			failUntilAttempt:  1,
			expectedError:     false,
			expectedCallCount: 2,
			testLogger:        false, // Will test with nil logger
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			callCount := 0
			var lastError error

			// Logger setup for testing
			var loggedErrors []error
			var loggedDelays []int
			var logger func(error, int)

			if tt.testLogger {
				logger = func(err error, delay int) {
					loggedErrors = append(loggedErrors, err)
					loggedDelays = append(loggedDelays, delay)
				}
			} // else logger remains nil

			fn := func() error {
				callCount++
				if tt.failUntilAttempt == 0 || callCount > tt.failUntilAttempt {
					return nil
				}

				// Use appropriate error message based on test case
				switch tt.expectedErrorMsg {
				case "persistent error":
					lastError = errors.New("persistent error")
				case "single failure":
					lastError = errors.New("single failure")
				default:
					lastError = fmt.Errorf("attempt %d failed", callCount)
				}
				return lastError
			}

			err := RetryCall(fn, tt.tries, tt.delay, tt.maxDelay, tt.backoff, tt.jitter, logger)

			// Check error expectation
			if tt.expectedError && err == nil {
				t.Errorf("Expected error, got nil")
			}
			if !tt.expectedError && err != nil {
				t.Errorf("Expected no error, got %v", err)
			}

			// Check specific error message for certain cases
			if tt.expectedError && tt.expectedErrorMsg != "" && err != nil {
				if tt.expectedErrorMsg == "retries exhausted" && err.Error() != "retries exhausted" {
					t.Errorf("Expected 'retries exhausted' error, got %v", err)
				} else if tt.expectedErrorMsg != "retries exhausted" && err != lastError {
					t.Errorf("Expected original error %v, got %v", lastError, err)
				}
			}

			// Check call count
			if callCount != tt.expectedCallCount {
				t.Errorf("Expected function to be called %d times, got %d", tt.expectedCallCount, callCount)
			}

			// Check logger functionality if enabled
			if tt.testLogger {
				if len(loggedErrors) != tt.expectedLogCount {
					t.Errorf("Expected %d logged errors, got %d", tt.expectedLogCount, len(loggedErrors))
				}

				if len(loggedDelays) != tt.expectedLogCount {
					t.Errorf("Expected %d logged delays, got %d", tt.expectedLogCount, len(loggedDelays))
				}

				// Check expected delay progression
				if len(tt.expectedLogDelays) > 0 {
					for i, expectedDelay := range tt.expectedLogDelays {
						if i < len(loggedDelays) && loggedDelays[i] != expectedDelay {
							t.Errorf("Expected logged delay[%d] = %d, got %d", i, expectedDelay, loggedDelays[i])
						}
					}
				}

				// Verify all delays respect maxDelay if set
				if tt.maxDelay != nil {
					for i, delay := range loggedDelays {
						if delay > *tt.maxDelay {
							t.Errorf("Logged delay[%d] = %d exceeds maxDelay %d", i, delay, *tt.maxDelay)
						}
					}
				}
			}
		})
	}
}

// Helper function to create int pointer
func intPtr(i int) *int {
	return &i
}

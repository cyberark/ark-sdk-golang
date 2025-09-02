package common

import (
	"fmt"
	"math/rand"
	"time"
)

// RetryCall executes a function with retry logic using exponential backoff.
//
// The function will be retried up to 'tries' times with an initial delay of 'delay' seconds.
// After each failed attempt, the delay is multiplied by 'backoff' and optional jitter is added.
// If maxDelay is specified, the delay will not exceed this value.
//
// Parameters:
//   - fn: The function to execute and retry on failure
//   - tries: Maximum number of attempts (must be > 0)
//   - delay: Initial delay between retries in seconds
//   - maxDelay: Optional maximum delay cap in seconds (nil for no limit)
//   - backoff: Multiplier applied to delay after each attempt
//   - jitter: Additional delay randomization - either fixed int or [2]int range
//   - logger: Optional callback to log retry attempts (receives error and current delay)
//
// Returns nil on success, or the last error encountered if all retries are exhausted.
//
// Example:
//
//	err := RetryCall(
//	    func() error { return someOperation() },
//	    3,     // max 3 attempts
//	    1,     // start with 1 second delay
//	    &10,   // cap at 10 seconds
//	    2,     // double delay each time
//	    [2]int{0, 500}, // add 0-500ms jitter
//	    func(err error, delay int) { log.Printf("Retry in %ds: %v", delay, err) },
//	)
func RetryCall(
	fn func() error,
	tries int,
	delay int,
	maxDelay *int,
	backoff int,
	jitter interface{},
	logger func(error, int),
) error {
	_tries, _delay := tries, delay
	for _tries != 0 {
		err := fn()
		if err == nil {
			return nil
		}

		_tries--
		if _tries == 0 {
			return err
		}

		if logger != nil {
			logger(err, _delay)
		}

		time.Sleep(time.Duration(_delay) * time.Second)
		_delay *= backoff

		switch j := jitter.(type) {
		case int:
			_delay += j
		case [2]int:
			_delay += rand.Intn(j[1]-j[0]) + j[0]
		}

		if maxDelay != nil && _delay > *maxDelay {
			_delay = *maxDelay
		}
	}
	return fmt.Errorf("retries exhausted")
}

package common

import (
	"fmt"
	"math/rand"
	"time"
)

// RetryCall retries a function call based on the provided parameters.
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

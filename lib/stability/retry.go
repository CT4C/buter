package stability

import (
	"fmt"
	"time"
)

type Caller func() (any, error)

/*
delay - In Milliseconds

Return nil as main value when retries is empty and
nothing returned from the Caller
*/
func Retry(caller Caller, retries int, delay int) (any, error) {
	result, err := caller()
	if err != nil && retries > 0 {
		// May be need to log here, when verbose needed
		fmt.Println("R", retries)
		time.Sleep(time.Duration(delay) * time.Millisecond)
		return Retry(caller, retries-1, delay)
	}

	if err != nil && retries <= 0 {
		return nil, err
	}

	return result, nil
}

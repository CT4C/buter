package requester

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/edpryk/buter/lib/stability"
)

var defaultHeaders = map[string]string{
	"Content-Type": "application/json",
	// Set diff in DOS mode
	"Connection": "close",
}

func AsyncRequestWithRetry(parameters RequestParameters, retries int, delay int) (<-chan CustomResponse, <-chan error) {
	resCh := make(chan CustomResponse, 1)
	errCh := make(chan error, 1)

	go func() {
		requestCaller := func() (any, error) {
			reader := strings.NewReader(parameters.Body.String())

			if parameters.Method == http.MethodPost {
				defaultHeaders["Content-Length"] = fmt.Sprintf("%d", len(parameters.Body.String()))
			}

			for key := range parameters.Header {
				defaultHeaders[key] = parameters.Header[key]
			}

			return Do(
				parameters.Method,
				parameters.Url,
				defaultHeaders,
				reader,
			)
		}
		startTime := time.Now()
		res, err := stability.Retry(requestCaller, retries, delay)
		if err != nil {
			errCh <- err
		} else {
			resCh <- CustomResponse{
				Response: res.(http.Response),
				Duration: time.Since(startTime),
				Payloads: parameters.Payloads,
			}
		}
	}()

	return resCh, errCh
}

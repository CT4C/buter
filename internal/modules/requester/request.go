package requester

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/edpryk/buter/lib/stability"
)

func AsyncRequestWithRetry(parameters RequestParameters, retries int, delay int) (<-chan CustomResponse, <-chan error) {
	resCh := make(chan CustomResponse, 1)
	errCh := make(chan error, 1)

	go func() {
		requestCaller := func() (any, error) {
			reader := strings.NewReader(parameters.Body.String())

			// TODO: move to separated func
			defaultHeaders := make(map[string]string)
			defaultHeaders["Connection"] = "close"

			if parameters.Method == http.MethodPost {
				defaultHeaders["Content-Length"] = fmt.Sprintf("%d", len(parameters.Body.String()))
				defaultHeaders["Content-Type"] = "application/json"
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

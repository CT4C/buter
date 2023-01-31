package requester

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/edpryk/buter/lib/stability"
)

func AsyncRequestWitnRetry(parameters ReuqestParameters, retries int, delay int) (<-chan CustomResponse, <-chan error) {
	resCh := make(chan CustomResponse)
	errCh := make(chan error)

	go func() {
		requstCaller := func() (any, error) {
			reader := strings.NewReader(parameters.Body.String())

			// TODO: move to seprated func
			defaultHeaders := make(map[string]string)

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
		res, err := stability.Retry(requstCaller, retries, delay)

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

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

			if parameters.Method == http.MethodPost {
				parameters.Header["Content-Length"] = fmt.Sprintf("%d", len(parameters.Body.String()))
			}

			for key := range parameters.Header {
				parameters.Header[key] = parameters.Header[key]
			}

			return Do(
				parameters.Method,
				parameters.Url,
				parameters.Header,
				reader,
			)
		}
		startTime := time.Now()
		res, err := stability.Retry(requestCaller, retries, delay)
		/*
			TODO: Need to realize where to close body
		*/
		defer res.(http.Response).Body.Close()

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

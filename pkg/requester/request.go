package requester

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/edpryk/buter/lib/stability"
)

const maxBodyLen = 10 * 1024

func updateHeaders(URL string, method string, body string, headers map[string]string) error {
	u, err := url.Parse(URL)
	if err != nil {
		return err
	}

	if method == http.MethodPost || method == http.MethodPatch || method == http.MethodPut {
		headers["Content-Length"] = fmt.Sprintf("%d", len(body))
	}

	headers["Host"] = u.Host

	return nil
}

func AsyncRequestWithRetry(parameters RequestParameters, retries int, delay int) (<-chan CustomResponse, <-chan error) {
	resCh := make(chan CustomResponse, 1)
	errCh := make(chan error, 1)

	updateHeaders(parameters.Url, parameters.Method, parameters.Body, parameters.Header)

	go func() {
		requestCaller := func() (any, error) {
			reader := strings.NewReader(parameters.Body)

			return Do(
				parameters.Method,
				parameters.Url,
				parameters.Header,
				reader,
			)
		}
		startTime := time.Now()
		res, err := stability.Retry(requestCaller, retries, delay)

		if err != nil {
			errCh <- err
		} else {
			defer res.(http.Response).Body.Close()

			data := make([]byte, maxBodyLen)
			n, err := res.(http.Response).Body.Read(data)
			if err != nil && err != io.EOF {
				errCh <- err
				return
			}

			resCh <- CustomResponse{
				Response: res.(http.Response),
				Duration: time.Since(startTime),
				Payloads: parameters.Payloads,
				Body:     data[:n],
			}
		}
	}()

	return resCh, errCh
}

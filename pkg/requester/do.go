package requester

import (
	"io"
	"net/http"
	"sync"
)

func Do(method string, url string, headers map[string]string, body io.Reader) (http.Response, error) {
	transport := &http.Transport{}

	client := http.Client{
		Transport: transport,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return http.Response{}, err
	}

	mu := sync.Mutex{}

	mu.Lock()
	for key := range headers {
		req.Header.Set(key, headers[key])
	}
	mu.Unlock()

	res, err := client.Do(req)
	if err != nil {
		return http.Response{}, err
	}

	return *res, nil
}

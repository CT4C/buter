package requester

import (
	"io"
	"net/http"
)

func Do(method string, url string, headers map[string]string, body io.Reader) (http.Response, error) {
	transport := &http.Transport{}

	client := http.Client{
		Transport: transport,
	}
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return http.Response{}, err
	}

	for key := range headers {
		req.Header.Set(key, headers[key])
	}

	res, err := client.Do(req)
	if err != nil {
		return http.Response{}, err
	}

	return *res, nil
}

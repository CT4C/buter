package requester

import (
	"log"
	"net/http"
)

func Do(method string, url string) http.Response {

	transport := &http.Transport{}

	client := http.Client{
		Transport: transport,
	}

	res, err := client.Get(url)
	if err != nil {
		log.Println(err)
		return http.Response{}
	}

	return *res
}

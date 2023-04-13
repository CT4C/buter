package reporter

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/edpryk/buter/lib/lists"
	"github.com/edpryk/buter/pkg/requester"
)

type Filters interface {
	Status() []string
	Length() []string
	Duration() []int
}

type Stopper interface {
	Status() []string
}

/*
	Apply interface
*/
type Response interface {
	ContentLength() int
	Status() int
	Duration() time.Duration
}

type Reporter struct{}

func (r Reporter) StartWorker(responseQ chan requester.CustomResponse, filters Filters, stopper Stopper, sopSig chan int) {
	requestNumber := 0

	for res := range responseQ {
		requestNumber++
		if len(filters.Length()) > 0 {
			if lists.Contain(filters.Length(), fmt.Sprint(len(res.Body))) {
				continue
			}
		}

		if len(filters.Status()) > 0 {
			if lists.Contain(filters.Status(), fmt.Sprint(res.StatusCode)) {
				continue
			}
		}

		if len(filters.Duration()) > 0 {
			if !lists.IntGreaterEq(filters.Duration(), int(res.Duration.Milliseconds())) {
				continue
			}
		}

		line := make(reportLine, 0)

		line.add("Req", requestNumber)

		for number, payload := range res.Payloads {
			line.add(fmt.Sprintf("P_%d", number+1), payload)
		}

		line.add("Status", res.StatusCode)
		line.add("Duration", res.Duration/time.Millisecond)
		line.add("Length", len(res.Body))

		if res.StatusCode == http.StatusFound {
			location, err := res.Location()
			if err == nil {
				line.add("Location", location.Path)
			}
		}

		log.Println(line.string())

		if len(stopper.Status()) > 0 {
			if lists.Contain(stopper.Status(), fmt.Sprint(res.StatusCode)) {
				sopSig <- 1
			}
		}
	}
}

func New() Reporter {
	return Reporter{}
}

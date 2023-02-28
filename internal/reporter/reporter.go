package reporter

import (
	"fmt"
	"log"
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
			if lists.Contain(filters.Length(), fmt.Sprint(res.ContentLength)) {
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

		report := fmt.Sprintf("%-3s %7d", "Req:", requestNumber)

		duration := res.Duration
		code := res.StatusCode
		payloads := ""
		for number, p := range res.Payloads {
			payloads += fmt.Sprintf("%-2sP_%d: %21s", " ", number+1, p)
		}

		report += payloads

		report += fmt.Sprintf(" Status %-5d", code)
		report += fmt.Sprintf("Duration %5dms", duration/time.Millisecond)
		report += fmt.Sprintf(" Length %5d", res.ContentLength)

		log.Println(report)

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

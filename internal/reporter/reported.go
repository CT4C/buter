package reporter

import (
	"fmt"
	"log"
	"time"

	"github.com/edpryk/buter/internal/requester"
	"github.com/edpryk/buter/lib/lists"
)

type Filters interface {
	Status() []int
	Length() []int
}

type Stopper interface {
	Status() []int
}

type Reporter struct{}

func (r Reporter) StartWorker(responseQ chan requester.CustomResponse, filters Filters, stopper Stopper, sopSig chan int) {
	counter := 0

	for res := range responseQ {
		counter++

		if len(filters.Length()) > 0 {
			if !lists.Contain(filters.Length(), int(res.ContentLength)) {
				continue
			}
		}

		if len(filters.Status()) > 0 {
			if !lists.Contain(filters.Status(), res.StatusCode) {
				continue
			}
		}

		report := fmt.Sprintf("%-3s %7d", "Req:", counter)

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
			if lists.Contain(stopper.Status(), res.StatusCode) {
				sopSig <- 1
			}
		}
	}
}

func New() Reporter {
	return Reporter{}
}

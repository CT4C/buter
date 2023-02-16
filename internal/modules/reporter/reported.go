package reporter

import (
	"fmt"
	"log"
	"time"

	"github.com/edpryk/buter/internal/modules/requester"
)

type Reporter struct{}

type Filters interface {
	Url() bool
	Status() bool
	Length() bool
	Duration() bool
}

func (r Reporter) StartWorker(responseQ chan requester.CustomResponse, filters Filters) {
	counter := 1

	for res := range responseQ {
		report := fmt.Sprintf("%3s Req: %7d", "", counter)

		duration := res.Duration
		code := res.StatusCode
		payloads := ""
		for number, p := range res.Payloads {
			payloads += fmt.Sprintf("%-2sP_%d: %21s", " ", number+1, p)
		}

		report += payloads

		if filters != nil {
		}

		report += fmt.Sprintf("%1sStatus %-5d", " ", code)
		report += fmt.Sprintf("Duration %5dms", duration/time.Millisecond)

		log.Println(report)
		counter++
	}
}

func New() Reporter {
	return Reporter{}
}

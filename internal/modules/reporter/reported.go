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
	for res := range responseQ {
		// if res.StatusCode > 300 {
		// 	continue
		// }
		report := fmt.Sprintf("%3s", "")

		duration := res.Duration
		// url := res.Request.URL
		code := res.StatusCode
		payloads := ""
		for number, p := range res.Payloads {
			payloads += fmt.Sprintf("%-5sP_%d: %32s", " ", number+1, p)
		}

		report += payloads

		if filters != nil {
		}

		report += fmt.Sprintf("%5sStatus %-5d", " ", code)
		report += fmt.Sprintf("Duration %5dms", duration/time.Millisecond)

		log.Println(report)
	}
}

func New() Reporter {
	return Reporter{}
}

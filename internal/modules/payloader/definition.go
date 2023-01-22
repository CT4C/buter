package payloader

import (
	"regexp"
)

var (
	rePayloadPosition = regexp.MustCompile("(![^!]+!)")
)

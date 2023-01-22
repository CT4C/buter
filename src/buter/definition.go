package buter

import (
	"regexp"
)

var (
	rePayloadPosition = regexp.MustCompile("(![^!]+!)")
)

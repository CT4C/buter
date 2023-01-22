package buter

import "errors"

var (
	errLessPayloadsThanPositions = errors.New("Less Payloads Than Positions")
	errLessPositionsThanPayloads = errors.New("Less Positions Then Payloads")
	errPayloadErr                = errors.New("Payload processing error")
	errAttackNotSupported        = errors.New("This attack type doesn't supported")
	errInvalidTotalPayload       = errors.New("Total Payloads Value is Invalid")
)

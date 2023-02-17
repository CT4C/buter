package buter

import "errors"

var (
	errPayloadErr                = errors.New("payload processing error")
	errAttackNotSupported        = errors.New("This attack type doesn't supported")
	errInvalidTotalPayload       = errors.New("Total Payloads Value is Invalid")
	errLessPositionsThanPayloads = errors.New("Less Positions Then Payloads")
	errLessPayloadsThanPositions = errors.New("Less Payloads Than Positions")
)

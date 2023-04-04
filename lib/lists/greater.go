package lists

func IntGreaterEq(list []int, value int) bool {
	isGreater := false

	for _, item := range list {
		if isGreater {
			break
		}

		isGreater = value >= item
	}

	return isGreater
}

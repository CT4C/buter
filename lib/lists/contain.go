package lists

func Contain[V comparable](list []V, value V) bool {
	exists := false

	for _, v := range list {
		if v == value {
			exists = true
			break
		}
	}

	return exists
}

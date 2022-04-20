package helpers

func Contains(array []string, element string) bool {
	for _, el := range array {
		if (el == element) {
			return true
		}
	}

	return false
}

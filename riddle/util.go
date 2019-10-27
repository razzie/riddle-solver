package riddle

func hasDuplicates(items []string) bool {
	count := len(items)
	for i := 0; i < count-1; i++ {
		val := items[i]
		for j := i + 1; j < count; j++ {
			if items[j] == val {
				return true
			}
		}
	}

	return false
}

func contains(items []string, value string) bool {
	for _, item := range items {
		if item == value {
			return true
		}
	}

	return false
}

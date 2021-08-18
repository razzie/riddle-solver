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

func copySlice(slice []string) []string {
	result := make([]string, 0, len(slice))
	for _, val := range slice {
		result = append(result, val)
	}
	return result
}

func uniqueAppend(slice []string, values []string) []string {
	for _, val := range values {
		if !contains(slice, val) {
			slice = append(slice, val)
		}
	}
	return slice
}

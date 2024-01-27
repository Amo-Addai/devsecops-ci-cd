package utils

// Unique - Strips out duplicate elements
func Unique(input []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range input {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

// Concat - Combine slices.
func Concat(one []string, two []string) []string {
	channelNames := make([]string, len(one) + len(two))
	for i, item := range one {
		channelNames[i] = item
	}
	for i, item := range two {
		channelNames[i + len(one)] = item
	}
	return channelNames
}
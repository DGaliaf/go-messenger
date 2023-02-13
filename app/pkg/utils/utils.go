package utils

func ContainDuplicates[T comparable](data []T) bool {
	counts := make(map[T]int)

	for _, d := range data {
		counts[d]++
	}

	for _, val := range counts {
		if val > 1 {
			return true
		}
	}

	return false
}

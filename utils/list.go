package utils

// FilterInt64 filters a slice of int64 elements and returns only those
func FilterInt64(slice []int64, filter func(int64) bool) (res []int64) {
	res = slice[:0]
	for _, x := range slice {
		if filter(x) {
			res = append(res, x)
		}
	}
	for i := len(res); i < len(slice); i++ {
		slice[i] = 0
	}
	return res
}

// ContainsInt64 checks if a slice of int64s contains an element
func ContainsInt64(haystack []int64, needle int64) bool {
	for _, element := range haystack {
		if element == needle {
			return true
		}
	}
	return false
}

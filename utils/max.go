package utils

// MaxInt returns the maximum value of a list of integers
func MaxInt(values ...int) (max int) {
	// if there are no elements, return
	if len(values) == 0 {
		return
	}

	// iterate over the values
	max = values[0]
	for _, v := range values[1:] {
		if v > max {
			max = v
		}
	}

	return
}

// MaxInt64 returns the maximum value of a list of integers
func MaxInt64(values ...int64) (max int64) {
	// if there are no elements, return
	if len(values) == 0 {
		return
	}

	// iterate over the values
	max = values[0]
	for _, v := range values[1:] {
		if v > max {
			max = v
		}
	}

	return
}

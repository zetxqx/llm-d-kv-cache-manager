package utils

// SliceMap applies a function to each element of a slice and returns a new
// slice with the results.
func SliceMap[Domain, Range any](slice []Domain, fn func(Domain) Range) []Range {
	if slice == nil {
		return nil
	}

	ans := make([]Range, len(slice))
	for idx, elt := range slice {
		ans[idx] = fn(elt)
	}

	return ans
}

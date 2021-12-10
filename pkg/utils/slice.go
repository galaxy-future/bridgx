package utils

func Filter(source []int64, filter func(int64) bool) []int64 {
	result := make([]int64, 0, len(source))
	for _, v := range source {
		a := v
		if filter(a) {
			result = append(result, a)
		}
	}
	return result
}

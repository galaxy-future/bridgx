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

// StringSliceDiff s1-s2
func StringSliceDiff(s1, s2 []string) []string {
	if len(s1) == 0 {
		return nil
	}
	if len(s2) == 0 {
		return s1
	}

	tmpMap := make(map[string]int, len(s1)+len(s2))
	diff := make([]string, 0, len(s1))
	for _, v := range s1 {
		tmpMap[v] += 1
	}
	for _, v := range s2 {
		tmpMap[v] += 2
	}
	for k, v := range tmpMap {
		if v == 1 {
			diff = append(diff, k)
		}
	}
	return diff
}

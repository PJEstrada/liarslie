package consensus

// FindMajorityValue finds the number that is repeated the most among the list of integers in values param.
func FindMajorityValue(values []int) int {
	valuesCount := map[int]int{}
	for _, val := range values {
		valuesCount[val] += 1
	}
	percentages := map[int]float64{}
	max := -1.0
	maxVal := -1
	for key, val := range valuesCount {
		percentages[key] += float64(float64(val) / float64(len(values)))
		if percentages[key] > float64(max) {
			max = percentages[key]
			maxVal = key
		}
	}
	return maxVal
}

// FindMajorityValuePercent finds the number that is repeated at least the minPercent ratio on the given int values.
func FindMajorityValuePercent(values []int, minPercent float32) int {
	valuesCount := map[int]int{}
	for _, val := range values {
		valuesCount[val] += 1
	}
	percentages := map[int]float64{}
	maxVal := -1
	for key, val := range valuesCount {
		percentages[key] += float64(float64(val) / float64(len(values)))
		if percentages[key] > float64(minPercent) {
			maxVal = key
		}
	}
	return maxVal
}

package sliceops

//import "math"

// Returns the minimum value in the slice and the index of
// the minimum value. If the input slice is empty, zero is returned
// as the minimum value and -1 is returned as the index.
// Use: val,ind := sliceops.Min(slice)
func Min(s []float64) (min float64, ind int) {
	if len(s) == 0 {
		return min, -1 // Ind is -1 to make clear it's not the zeroth index.
	}
	min = s[0]
	ind = 0
	for i, val := range s {
		if val < min {
			min = val
			ind = i
		}
	}
	return min, ind
}

// Returns the maximum value in the slice and the location of
// the maximum value. If the input slice is empty, zero is returned
// as the minimum value and -1 is returned as the index.
// Use: val,ind := sliceops.Max(slice)
func Max(s []float64) (max float64, ind int) {
	if len(s) == 0 {
		return max, -1 // Ind is -1 to make clear it's not the zeroth index.
	}
	max = s[0]
	ind = 0
	for i, val := range s {
		if val > max {
			max = val
			ind = i
		}
	}
	return max, ind
}

// Returns the sum of the elements of the slice
func Sum(s []float64) (sum float64) {
	for _, val := range s {
		sum += val
	}
	return
}

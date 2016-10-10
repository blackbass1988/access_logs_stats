package core

type Float64Data []float64

func (a Float64Data) Len() int           { return len(a) }
func (a Float64Data) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a Float64Data) Less(i, j int) bool { return a[i] < a[j] }

func (a Float64Data) Min() float64 {

	if len(a) == 0 {
		return 0
	}

	return a[0]
}

func (a Float64Data) Max() float64 {

	if len(a) == 0 {
		return 0
	}

	return a[len(a)-1]
}

func (a Float64Data) Avg() float64 {

	if len(a) == 0 {
		return 0
	}

	var sum float64
	sum = 0
	for _, num := range a {
		sum += num
	}

	return float64(sum) / float64(len(a))
}
func (a Float64Data) Percentile(cent float64) float64 {
	if len(a) == 0 {
		return 0
	}

	sliceSize := int(float64(len(a)) * cent / 100)

	if sliceSize == 0 {
		return a[0]
	}

	slice := a[sliceSize-1:]
	return float64(slice[0])
}

func (a Float64Data) ItemsPerSeconds(seconds float64) float64 {
	return float64(len(a)) / seconds
}

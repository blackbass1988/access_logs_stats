package pkg

//Float64Data extends []float64 with data analyze functions
type Float64Data []float64

//Len returns Len of Float64Data slice
func (a Float64Data) Len() int { return len(a) }

//Swap swaps 2 values in slice
func (a Float64Data) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

//Less returns true if i < j
func (a Float64Data) Less(i, j int) bool { return a[i] < a[j] }

//Min returns min value in Floatdata
func (a Float64Data) Min() float64 {

	if len(a) == 0 {
		return 0
	}

	return a[0]
}

//Max returns max value in Floatdata
func (a Float64Data) Max() float64 {

	if len(a) == 0 {
		return 0
	}

	return a[len(a)-1]
}

//Avg returns Avg of values in Floatdata
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

//Percentile returns x percentile of values in Floatdata
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

//Sum returns sum of values in Floatdata
func (a Float64Data) Sum() float64 {
	var sum float64

	if len(a) == 0 {
		return 0
	}

	sum = 0
	for _, v := range a {
		sum += v
	}
	return sum
}

//ItemsPerSeconds returns items per second in Floatdata
func (a Float64Data) ItemsPerSeconds(seconds float64) float64 {
	return float64(len(a)) / seconds
}

package core_test

import (
	"github.com/blackbass1988/access_logs_stats/core"
	"sort"
	"testing"
	"time"
)

func TestFloatData(t *testing.T) {

	numbers := []float64{3, 1.0, 2}
	floatNumber := core.Float64Data(numbers)
	sort.Sort(floatNumber)

	if floatNumber[0] != 1 || floatNumber[1] != 2 || floatNumber[2] != 3 {
		t.Error("incorrect order")
	}

	if floatNumber.Len() != 3 {
		t.Errorf("incorrect len. must 3 but was %d", floatNumber.Len())
	}

	if floatNumber.Avg() != 2 {
		t.Errorf("incorrect avg. must 2 but was %f", floatNumber.Avg())
	}

	if floatNumber.Percentile(100) != 3 {
		t.Errorf("incorrect percentile(100). must 3 but was %f", floatNumber.Percentile(100))
	}

	if floatNumber.Percentile(50) != 1 {
		t.Errorf("incorrect percentile(50). must 3 but was %f", floatNumber.Percentile(50))
	}

	if floatNumber.Percentile(10) != 1 {
		t.Errorf("incorrect percentile(10). must 1 but was %f", floatNumber.Percentile(10))
	}

	if floatNumber.Min() != 1 {
		t.Errorf("incorrect Min(). must 1 but was %f", floatNumber.Min())
	}

	if floatNumber.Max() != 3 {
		t.Errorf("incorrect Max(). must 3 but was %f", floatNumber.Max())
	}

	if floatNumber.Sum() != 6 {
		t.Errorf("incorrect Max(). must 6 but was %f", floatNumber.Sum())
	}

}

type lendivpurCase struct {
	numbers  []float64
	duration string
	expected float64
}

func TestLenDivDuration(t *testing.T) {

	var tests = []lendivpurCase{
		{[]float64{3, 1.0, 2}, "1s", 3},
		{[]float64{3, 1.0, 2}, "2s", 1.5},
		{[]float64{3, 1.0, 2}, "3s", 1},
	}

	for _, test := range tests {
		duration, _ := time.ParseDuration(test.duration)
		numbers := test.numbers
		floatNumber := core.Float64Data(numbers)
		if floatNumber.ItemsPerSeconds(duration.Seconds()) != test.expected {
			t.Error("ItemsPerDuration must ", test.expected,
				" but was ", floatNumber.ItemsPerSeconds(duration.Seconds()))
		}
	}
}

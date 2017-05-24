package pkg_test

import (
	"fmt"
	"math/rand"
	"strconv"
	"testing"
	"time"
)

var size = 1000000

func getSlice() []string {
	slice := []string{}
	for i := 0; i < size; i++ {
		val := "a" + strconv.Itoa(i)
		slice = append(slice, val)
	}
	return slice
}

func getMap() map[string]string {
	m := map[string]string{}

	for i := 0; i < size; i++ {
		val := "a" + strconv.Itoa(i)
		m[val] = val
	}
	return m
}

func random(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min) + min
}

func BenchmarkSlice(b *testing.B) {

	slice := getSlice()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		index := "a" + strconv.Itoa(random(0, size-1))

		for _, val := range slice {
			if val == index {
				_ = fmt.Sprint(val)
				break
			}
		}

		_ = fmt.Sprint(random(0, size-1))
	}
}

func BenchmarkMap(b *testing.B) {

	m := getMap()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		index := "a" + strconv.Itoa(random(0, size-1))

		if val, ok := m[index]; ok {
			_ = fmt.Sprint(val)
		}

	}
}

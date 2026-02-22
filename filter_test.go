package dblchk_test

import (
	"crypto/rand"
	"fmt"
	"testing"

	"github.com/mattwiller/dblchk"
)

var BenchmarkResult any

func TestEmptyFilter(t *testing.T) {
	filter := dblchk.NewFilter(0)
	assertEq(t, filter.MayContain([]byte("hello, world!")), false)
}

func BenchmarkNewFilter(b *testing.B) {
	var filter dblchk.Filter
	for b.Loop() {
		filter = dblchk.NewFilter(0)
	}
	BenchmarkResult = filter
}

func BenchmarkReset(b *testing.B) {
	b.StopTimer()
	filter := dblchk.NewFilter(0)

	b.StartTimer()
	for b.Loop() {
		filter.Reset()
	}

	BenchmarkResult = filter
}

func TestAddElement(t *testing.T) {
	element := []byte("hello, world!")
	filter := dblchk.NewFilter(0)
	filter.Add(element)

	assertEq(t, filter.MayContain(element), true)
}

func BenchmarkAddElement(b *testing.B) {
	b.StopTimer()
	filter := dblchk.NewFilter(0)
	element := []byte("hello, world!")

	b.StartTimer()
	for b.Loop() {
		filter.Add(element)
	}
	BenchmarkResult = filter
}

func TestMayContain(t *testing.T) {
	filter := dblchk.NewFilter(0)
	var elements [][]byte
	for range 10_000 {
		element := addToFilter(filter)
		elements = append(elements, element)
	}

	for _, element := range elements {
		assertEq(t, filter.MayContain(element), true)
	}
	for range 100 {
		randomElement, err := getRandomBytes()
		assertEq(t, err == nil, true)
		result := filter.MayContain(randomElement)
		assertEq(t, result, false)
	}
}

func BenchmarkMayContain(b *testing.B) {
	b.StopTimer()
	filter := dblchk.NewFilter(0)

	for range 10_000 {
		addToFilter(filter)
	}

	key, _ := getRandomBytes()
	var result bool
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		result = filter.MayContain(key)
	}
	BenchmarkResult = result
}

func getRandomBytes() ([]byte, error) {
	element := make([]byte, 16)
	_, err := rand.Read(element)
	return element, err
}

func addToFilter(filter dblchk.Filter) []byte {
	element, _ := getRandomBytes()
	filter.Add(element)
	return element
}

func assertEq(t *testing.T, a, b any) {
	t.Helper()
	if a != b {
		t.Error(a, "did not match expected value:", b)
	}
}

func TestFalsePositiveRate(t *testing.T) {
	filter := dblchk.NewFilter(0)
	n := int(1e5)
	for range n {
		addToFilter(filter)
	}

	collisions := 0
	for range int(1e6) {
		randomElement, _ := getRandomBytes()
		collision := filter.MayContain(randomElement)
		if collision {
			collisions++
		}
	}
	fmt.Printf("observed dblcheck %.2f%% collision rate for m=256KiB, n=100k\n", float64(collisions)/1e6*100)
}

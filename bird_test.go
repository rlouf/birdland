package birdland

import (
	"io/ioutil"
	"log"
	"math/rand"
	"testing"
)

// Sampling X items from the incoming query
// ////////////////////////////////////////////////////////////////////////////

func benchmarkSampleItemsFromQuery(querySize, numItems int, b *testing.B) {
	log.SetFlags(0) // don't let logs pollute the benchmarks
	log.SetOutput(ioutil.Discard)

	query := make([]QueryItem, querySize)
	for i := 0; i < querySize; i++ {
		item := QueryItem{
			Item:   rand.Intn(numItems),
			Weight: 10 * rand.Float64(),
		}
		query[i] = item
	}

	itemsWeights := make([]float64, numItems)
	for i := 0; i < numItems; i++ {
		itemsWeights[i] = 10 * rand.Float64()
	}
	bird, err := NewBird(itemsWeights, nil, nil)
	if err != nil {
		panic(`BenchmarkSampleItems: Bird initialization raised an error 
			but shouldn't have. Check your test case`)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = bird.SampleItemsFromQuery(query)
	}
}

func BenchmarkSampleItemsFormQuery10Query2000000Items(b *testing.B) {
	benchmarkSampleItemsFromQuery(10, 2000000, b)
}

func BenchmarkSampleItemsFormQuery100Query2000000Items(b *testing.B) {
	benchmarkSampleItemsFromQuery(100, 2000000, b)
}

func BenchmarkSampleItemsFormQuery1000Query2000000Items(b *testing.B) {
	benchmarkSampleItemsFromQuery(1000, 2000000, b)
}

// One step of the recommendation process with random init
// ////////////////////////////////////////////////////////////////////////////

func benchmarkStep(querySize, numUsers, numItems int, b *testing.B) {
	log.SetFlags(0) // don't let logs pollute the benchmarks
	log.SetOutput(ioutil.Discard)

	itemsToUsers := make([][]int, numItems)
	for i := 0; i < numItems; i++ {
		itemsToUsers[i] = []int{1}
	}

	usersToItems := make([][]int, numUsers)
	for i := 0; i < numUsers; i++ {
		num := 1 + rand.Intn(100) // +1 so that num != 0
		items := make([]int, num)
		for j := 0; j < num; j++ {
			it := rand.Intn(numItems)
			items[j] = it
			itemsToUsers[it] = append(itemsToUsers[it], i)
		}
		usersToItems[i] = items
	}

	itemWeights := make([]float64, numItems)
	for i := 0; i < numItems; i++ {
		itemWeights[i] = 10 * rand.Float64()
	}
	bird, err := NewBird(itemWeights, usersToItems, itemsToUsers, Draws(1000), Depth(2))
	if err != nil {
		panic(`BenchmarkStep: Bird initialization raised an error
			but shouldn't have. Check your test case`)
	}

	query := make([]int, querySize)
	for i := 0; i < querySize; i++ {
		query[i] = rand.Intn(numItems)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _ = bird.Step(query)
	}
}

func BenchmarkStep2000000Items1000Users100Query(b *testing.B) {
	benchmarkStep(100, 1000, 2000000, b)
}
func BenchmarkStep2000000Items100000Users100Query(b *testing.B) {
	benchmarkStep(100, 100000, 2000000, b)
}
func BenchmarkStep2000000Items1000000Users100Query(b *testing.B) {
	benchmarkStep(100, 1000000, 2000000, b)
}
func BenchmarkStep2000000Items1000000Users1000Query(b *testing.B) {
	benchmarkStep(1000, 1000000, 2000000, b)
}
func BenchmarkStep2000000Items1000000Users10000Query(b *testing.B) {
	benchmarkStep(10000, 1000000, 2000000, b)
}
func BenchmarkStep2000000Items1000000Users100000Query(b *testing.B) {
	benchmarkStep(100000, 1000000, 2000000, b)
}
func BenchmarkStep2000000Items1000000Users200000Query(b *testing.B) {
	benchmarkStep(200000, 1000000, 2000000, b)
}
func BenchmarkStep2000000Items1000000Users300000Query(b *testing.B) {
	benchmarkStep(300000, 1000000, 2000000, b)
}
func BenchmarkStep2000000Items1000000Users400000Query(b *testing.B) {
	benchmarkStep(400000, 1000000, 2000000, b)
}
func BenchmarkStep2000000Items1000000Users600000Query(b *testing.B) {
	benchmarkStep(600000, 1000000, 2000000, b)
}
func BenchmarkStep2000000Items1000000Users700000Query(b *testing.B) {
	benchmarkStep(700000, 1000000, 2000000, b)
}
func BenchmarkStep2000000Items1000000Users800000Query(b *testing.B) {
	benchmarkStep(800000, 1000000, 2000000, b)
}
func BenchmarkStep2000000Items1000000Users900000Query(b *testing.B) {
	benchmarkStep(900000, 1000000, 2000000, b)
}
func BenchmarkStep2000000Items1000000Users1000000Query(b *testing.B) {
	benchmarkStep(1000000, 1000000, 2000000, b)
}

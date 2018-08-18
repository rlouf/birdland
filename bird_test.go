package birdland

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"testing"
)

type SampleItemCase struct {
	Name         string
	ItemWeights  []float64
	FromItems    []int
	ExpectedItem int
	Valid        bool
}

var sampleitem_table = []SampleItemCase{
	{
		Name:         "Uniform weights",
		ItemWeights:  []float64{1, 1, 1, 1, 1, 1},
		FromItems:    []int{3, 2, 1, 4, 5},
		ExpectedItem: 2,
		Valid:        true,
	},
}

func TestSampleItem(t *testing.T) {
	for _, ex := range sampleitem_table {
		b, err := NewBird(ex.ItemWeights, nil, nil)
		if err != nil {
			panic(fmt.Sprintf(`%s: Bird initialization raised an error 
							but shouldn't have. Check your test case`, ex.Name))
		}
		sampledItem, err := b.SampleItem(ex.FromItems)
		if err != nil {
			t.Errorf("SampleItem: '%s': unexpected error %v", ex.Name, err)
		}
		if sampledItem != ex.ExpectedItem {
			t.Errorf("SampleItem: '%s': expected %d, got %d", ex.Name, ex.ExpectedItem, sampledItem)
		}
	}
}

// Sampling a single item
// ////////////////////////////////////////////////////////////////////////////

func benchmarkSampleItem(numItems, numFrom int, b *testing.B) {
	log.SetFlags(0) // don't let logs pollute the benchmarks
	log.SetOutput(ioutil.Discard)

	itemWeights := make([]float64, numItems)
	for i := 0; i < numItems; i++ {
		itemWeights[i] = 10 * rand.Float64()
	}
	bird, _ := NewBird(itemWeights, nil, nil)

	fromItems := make([]int, numFrom)
	for i := 0; i < numFrom; i++ {
		fromItems[i] = rand.Intn(numItems)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = bird.SampleItem(fromItems)
	}
}

func BenchmarkSampleItem1000Items1000From(b *testing.B)    { benchmarkSampleItem(1000, 1000, b) }
func BenchmarkSampleItem10000Items1000From(b *testing.B)   { benchmarkSampleItem(10000, 1000, b) }
func BenchmarkSampleItem100000Items1000From(b *testing.B)  { benchmarkSampleItem(100000, 1000, b) }
func BenchmarkSampleItem1000000Items1000From(b *testing.B) { benchmarkSampleItem(1000000, 1000, b) }
func BenchmarkSampleItem1000000Items100From(b *testing.B)  { benchmarkSampleItem(1000000, 100, b) }
func BenchmarkSampleItem1000000Items10From(b *testing.B)   { benchmarkSampleItem(1000000, 10, b) }
func BenchmarkSampleItem2000000Items100From(b *testing.B)  { benchmarkSampleItem(2000000, 100, b) } // 2M on Spotify
func BenchmarkSampleItem2000000Items10From(b *testing.B)   { benchmarkSampleItem(2000000, 10, b) }

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
	bird, err := NewBird(itemWeights, usersToItems, itemsToUsers)
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

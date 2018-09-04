package birdland

import (
	"io/ioutil"
	"log"
	"math/rand"
	"testing"
)

type InitCase struct {
	Name         string
	ItemWeights  []float64
	UsersToItems [][]int
	Draws        int
	Depth        int
	Valid        bool
}

var init_table = []InitCase{
	{
		Name:         "Zero Depth",
		ItemWeights:  []float64{1, 1},
		UsersToItems: [][]int{[]int{0}, []int{1}},
		Depth:        0,
		Draws:        1,
		Valid:        false,
	},
	{
		Name:         "Negative Depth",
		ItemWeights:  []float64{1, 1},
		UsersToItems: [][]int{[]int{0}, []int{1}},
		Depth:        -1,
		Draws:        1,
		Valid:        false,
	},
	{
		Name:         "Zero Draws",
		ItemWeights:  []float64{1, 1},
		UsersToItems: [][]int{[]int{0}, []int{1}},
		Depth:        1,
		Draws:        0,
		Valid:        false,
	},
	{
		Name:         "Negative Draws",
		ItemWeights:  []float64{1, 1},
		UsersToItems: [][]int{[]int{0}, []int{1}},
		Depth:        1,
		Draws:        -1,
		Valid:        false,
	},
	{
		Name:         "Empty ItemWeights",
		ItemWeights:  []float64{},
		UsersToItems: [][]int{[]int{0}, []int{1}},
		Depth:        1,
		Draws:        1,
		Valid:        false,
	},
	{
		Name:         "Empty UsersToItems",
		ItemWeights:  []float64{1, 1},
		UsersToItems: [][]int{},
		Depth:        1,
		Draws:        1,
		Valid:        false,
	},
	{
		Name:         "More items in adjacency tables that weight list",
		ItemWeights:  []float64{0.1, 0.2, 0.4},
		UsersToItems: [][]int{[]int{0, 2}, []int{4}},
		Depth:        1,
		Draws:        1,
		Valid:        false,
	},
	{
		Name:         "Perfectly valid input",
		ItemWeights:  []float64{1, 1},
		UsersToItems: [][]int{[]int{0}, []int{1}},
		Depth:        1,
		Draws:        1,
		Valid:        true,
	},
}

func TestInitialization(t *testing.T) {
	for _, ex := range init_table {
		_, err := NewBird(ex.ItemWeights, ex.UsersToItems, Draws(ex.Draws), Depth(ex.Depth))
		if err != nil && ex.Valid {
			t.Errorf("Initialization: %s: Bird initialization should not have raised "+
				"an error but did: %v", ex.Name, err)
		}
		if err == nil && !ex.Valid {
			t.Errorf("Initialization: %s: Bird initialization should have raised "+
				"an error but did not", ex.Name)
		}
	}
}

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
	bird := Bird{
		ItemWeights: itemsWeights,
		RandSource:  rand.New(rand.NewSource(42)),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = bird.sampleItemsFromQuery(query)
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

func benchmarkStep(querySize, numUsers, numItems int, b *testing.B) {
	log.SetFlags(0) // don't let logs pollute the benchmarks
	log.SetOutput(ioutil.Discard)

	usersToItems := make([][]int, numUsers)
	for i := 0; i < numUsers; i++ {
		num := 1 + rand.Intn(100) // +1 so that num != 0
		items := make([]int, num)
		for j := 0; j < num; j++ {
			it := rand.Intn(numItems)
			items[j] = it
		}
		usersToItems[i] = items
	}

	itemWeights := make([]float64, numItems)
	for i := 0; i < numItems; i++ {
		itemWeights[i] = 10 * rand.Float64()
	}
	bird, err := NewBird(itemWeights, usersToItems)
	if err != nil {
		panic("BenchmarkStep: Bird initialization raised an error " +
			"but shouldn't have. Check your test case")
	}

	query := make([]int, querySize)
	for i := 0; i < querySize; i++ {
		query[i] = rand.Intn(numItems)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _ = bird.step(query)
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

func benchmarkProcess(numItems, numUsers, querySize, draws, depth int, b *testing.B) {
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
	bird, err := NewBird(itemWeights, usersToItems, Draws(draws), Depth(depth))
	if err != nil {
		panic("BenchmarkStep: Bird initialization raised an error " +
			"but shouldn't have. Check your test case")
	}

	query := make([]QueryItem, querySize)
	for i := 0; i < querySize; i++ {
		queryItem := QueryItem{Item: rand.Intn(numItems), Weight: rand.Float64()}
		query[i] = queryItem
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _ = bird.Process(query)
	}
}

// Increase the numbers of users to 10M
func BenchmarkProcess10KUsers(b *testing.B) {
	benchmarkProcess(2000000, 10000, 100, 100, 1, b)
}

func BenchmarkProcess100KUsers(b *testing.B) {
	benchmarkProcess(2000000, 100000, 100, 100, 1, b)
}

func BenchmarkProcess1MUsers(b *testing.B) {
	benchmarkProcess(2000000, 1000000, 100, 100, 1, b)
}

// Increase the numbers of draws to 100k with 1M users
func BenchmarkProcess100Draws(b *testing.B) {
	benchmarkProcess(2000000, 1000000, 100, 100, 1, b)
}

func BenchmarkProcess1KDraws(b *testing.B) {
	benchmarkProcess(2000000, 1000000, 100, 1000, 1, b)
}

func BenchmarkProcess10KDraws(b *testing.B) {
	benchmarkProcess(2000000, 1000000, 100, 10000, 1, b)
}

func BenchmarkProcess100KDraws(b *testing.B) {
	benchmarkProcess(2000000, 1000000, 100, 100000, 1, b)
}

// Increase the depth up to 10 with 1M users and 10K draws
func BenchmarkProcess1Depth(b *testing.B) {
	benchmarkProcess(2000000, 1000000, 100, 10000, 1, b)
}

func BenchmarkProcess2Depth(b *testing.B) {
	benchmarkProcess(2000000, 1000000, 100, 10000, 2, b)
}

func BenchmarkProcess3Depth(b *testing.B) {
	benchmarkProcess(2000000, 1000000, 100, 10000, 3, b)
}

func BenchmarkProcess4Depth(b *testing.B) {
	benchmarkProcess(2000000, 1000000, 100, 10000, 4, b)
}

func BenchmarkProcess5Depth(b *testing.B) {
	benchmarkProcess(2000000, 1000000, 100, 10000, 5, b)
}

func BenchmarkProcess10Depth(b *testing.B) {
	benchmarkProcess(2000000, 1000000, 100, 10000, 10, b)
}

package birdland

import (
	"math/rand"
	"testing"
)

type BirdInitCase struct {
	Name         string
	ItemWeights  []float64
	UsersToItems [][]int
	Draws        int
	Depth        int
	Valid        bool
}

var birdInitTable = []BirdInitCase{
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
		UsersToItems: [][]int{[]int{0, 1}, []int{2, 3}},
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

func TestBirdInitialization(t *testing.T) {
	for _, ex := range birdInitTable {
		cfg := NewBirdCfg()
		cfg.Depth = ex.Depth
		cfg.Draws = ex.Draws

		_, err := NewBird(cfg, ex.ItemWeights, ex.UsersToItems)
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


func benchmarkBirdSampleItemsFromQuery(querySize, numItems int, b *testing.B) {
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

	itemList := make([]int, numItems)
	for i :=0; i < numItems; i++ {
		itemList[i] = i
	}
	usersToItems := [][]int{itemList}

	bird,err := NewBird(NewBirdCfg(), itemsWeights, usersToItems)
	if err != nil {
		b.Error("Unable to initialize SampleItemsFromQuery benchmark")
	}  
	bird.RandSource = rand.New(rand.NewSource(42))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = bird.sampleItemsFromQuery(query)
	}
}

func BenchmarkBirdSampleItemsFormQuery10Query2000000Items(b *testing.B) {
	benchmarkBirdSampleItemsFromQuery(10, 2000000, b)
}

func BenchmarkBirdSampleItemsFormQuery100Query2000000Items(b *testing.B) {
	benchmarkBirdSampleItemsFromQuery(100, 2000000, b)
}

func BenchmarkBirdSampleItemsFormQuery1000Query2000000Items(b *testing.B) {
	benchmarkBirdSampleItemsFromQuery(1000, 2000000, b)
}

func benchmarkBirdStep(querySize, numUsers, numItems int, b *testing.B) {
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

	bird, err := NewBird(NewBirdCfg(), itemWeights, usersToItems)
	if err != nil {
		panic("BenchmarkBirdStep: Bird initialization raised an error " +
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

func BenchmarkBirdStep2000000Items1000Users100Query(b *testing.B) {
	benchmarkBirdStep(100, 1000, 2000000, b)
}

func BenchmarkBirdStep2000000Items100000Users100Query(b *testing.B) {
	benchmarkBirdStep(100, 100000, 2000000, b)
}

func BenchmarkBirdStep2000000Items1000000Users100Query(b *testing.B) {
	benchmarkBirdStep(100, 1000000, 2000000, b)
}

func BenchmarkBirdStep2000000Items1000000Users1000Query(b *testing.B) {
	benchmarkBirdStep(1000, 1000000, 2000000, b)
}

func BenchmarkBirdStep2000000Items1000000Users10000Query(b *testing.B) {
	benchmarkBirdStep(10000, 1000000, 2000000, b)
}

func BenchmarkBirdStep2000000Items1000000Users100000Query(b *testing.B) {
	benchmarkBirdStep(100000, 1000000, 2000000, b)
}

func BenchmarkBirdStep2000000Items1000000Users200000Query(b *testing.B) {
	benchmarkBirdStep(200000, 1000000, 2000000, b)
}

func BenchmarkBirdStep2000000Items1000000Users300000Query(b *testing.B) {
	benchmarkBirdStep(300000, 1000000, 2000000, b)
}

func BenchmarkBirdStep2000000Items1000000Users400000Query(b *testing.B) {
	benchmarkBirdStep(400000, 1000000, 2000000, b)
}

func BenchmarkBirdStep2000000Items1000000Users600000Query(b *testing.B) {
	benchmarkBirdStep(600000, 1000000, 2000000, b)
}

func BenchmarkBirdStep2000000Items1000000Users700000Query(b *testing.B) {
	benchmarkBirdStep(700000, 1000000, 2000000, b)
}

func BenchmarkBirdStep2000000Items1000000Users800000Query(b *testing.B) {
	benchmarkBirdStep(800000, 1000000, 2000000, b)
}

func BenchmarkBirdStep2000000Items1000000Users900000Query(b *testing.B) {
	benchmarkBirdStep(900000, 1000000, 2000000, b)
}

func BenchmarkBirdStep2000000Items1000000Users1000000Query(b *testing.B) {
	benchmarkBirdStep(1000000, 1000000, 2000000, b)
}

func benchmarkBirdProcess(numItems, numUsers, querySize, draws, depth int, b *testing.B) {
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

	cfg := NewBirdCfg()
	cfg.Depth = depth
	cfg.Draws = draws

	bird, err := NewBird(cfg, itemWeights, usersToItems)
	if err != nil {
		panic("BenchmarkBirdStep: Bird initialization raised an error " +
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
func BenchmarkBirdProcess10KUsers(b *testing.B) {
	benchmarkBirdProcess(2000000, 10000, 100, 100, 1, b)
}

func BenchmarkBirdProcess100KUsers(b *testing.B) {
	benchmarkBirdProcess(2000000, 100000, 100, 100, 1, b)
}

func BenchmarkBirdProcess1MUsers(b *testing.B) {
	benchmarkBirdProcess(2000000, 1000000, 100, 100, 1, b)
}

// Increase the numbers of draws to 100k with 1M users
func BenchmarkBirdProcess100Draws(b *testing.B) {
	benchmarkBirdProcess(2000000, 1000000, 100, 100, 1, b)
}

func BenchmarkBirdProcess1KDraws(b *testing.B) {
	benchmarkBirdProcess(2000000, 1000000, 100, 1000, 1, b)
}

func BenchmarkBirdProcess10KDraws(b *testing.B) {
	benchmarkBirdProcess(2000000, 1000000, 100, 10000, 1, b)
}

func BenchmarkBirdProcess100KDraws(b *testing.B) {
	benchmarkBirdProcess(2000000, 1000000, 100, 100000, 1, b)
}

// Increase the depth up to 10 with 1M users and 10K draws
func BenchmarkBirdProcess1Depth(b *testing.B) {
	benchmarkBirdProcess(2000000, 1000000, 100, 10000, 1, b)
}

func BenchmarkBirdProcess2Depth(b *testing.B) {
	benchmarkBirdProcess(2000000, 1000000, 100, 10000, 2, b)
}

func BenchmarkBirdProcess3Depth(b *testing.B) {
	benchmarkBirdProcess(2000000, 1000000, 100, 10000, 3, b)
}

func BenchmarkBirdProcess4Depth(b *testing.B) {
	benchmarkBirdProcess(2000000, 1000000, 100, 10000, 4, b)
}

func BenchmarkBirdProcess5Depth(b *testing.B) {
	benchmarkBirdProcess(2000000, 1000000, 100, 10000, 5, b)
}

func BenchmarkBirdProcess10Depth(b *testing.B) {
	benchmarkBirdProcess(2000000, 1000000, 100, 10000, 10, b)
}

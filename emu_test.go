package birdland

import (
	"math/rand"
	"testing"
)

type EmuInitCase struct {
	Name                 string
	ItemWeights          []float64
	UsersToWeightedItems []map[int]float64
	Draws                int
	Depth                int
	Valid                bool
}

var emuInitTable = []EmuInitCase{
	{
		Name:                 "Zero Depth",
		ItemWeights:          []float64{1, 1},
		UsersToWeightedItems: []map[int]float64{{0: 1.}, {1: 1.}},
		Depth:                0,
		Draws:                1,
		Valid:                false,
	},
	{
		Name:                 "Negative Depth",
		ItemWeights:          []float64{1, 1},
		UsersToWeightedItems: []map[int]float64{{0: 1.}, {1: 1.}},
		Depth:                -1,
		Draws:                1,
		Valid:                false,
	},
	{
		Name:                 "Zero Draws",
		ItemWeights:          []float64{1, 1},
		UsersToWeightedItems: []map[int]float64{{0: 1.}, {1: 1.}},
		Depth:                1,
		Draws:                0,
		Valid:                false,
	},
	{
		Name:                 "Negative Draws",
		ItemWeights:          []float64{1, 1},
		UsersToWeightedItems: []map[int]float64{{0: 1.}, {1: 1.}},
		Depth:                1,
		Draws:                -1,
		Valid:                false,
	},
	{
		Name:                 "Empty ItemWeights",
		ItemWeights:          []float64{},
		UsersToWeightedItems: []map[int]float64{{0: 1.}, {1: 1.}},
		Depth:                1,
		Draws:                1,
		Valid:                false,
	},
	{
		Name:                 "Empty UsersToWeightedItems",
		ItemWeights:          []float64{1, 1},
		UsersToWeightedItems: []map[int]float64{},
		Depth:                1,
		Draws:                1,
		Valid:                false,
	},
	{
		Name:                 "Negative weight in UsersToWeightedItems",
		ItemWeights:          []float64{1, 1},
		UsersToWeightedItems: []map[int]float64{{0: 1.}, {1: -1.}},
		Depth:                1,
		Draws:                1,
		Valid:                false,
	},
	{
		Name:                 "More items in adjacency tables that weight list",
		ItemWeights:          []float64{0.1, 0.2, 0.4},
		UsersToWeightedItems: []map[int]float64{{0: 1., 1: 1.}, {2: 1., 3: 1.}},
		Depth:                1,
		Draws:                1,
		Valid:                false,
	},
	{
		Name:                 "Perfectly valid input",
		ItemWeights:          []float64{1, 1},
		UsersToWeightedItems: []map[int]float64{{0: 1.}, {1: 1.}},
		Depth:                1,
		Draws:                1,
		Valid:                true,
	},
}

func TestEmuInitialization(t *testing.T) {
	for _, ex := range emuInitTable {
		cfg := NewBirdCfg()
		cfg.Depth = ex.Depth
		cfg.Draws = ex.Draws

		_, err := NewEmu(cfg, ex.ItemWeights, ex.UsersToWeightedItems)
		if err != nil && ex.Valid {
			t.Errorf("EmuInitialization: %s: Bird initialization should not have raised "+
				"an error but did: %v", ex.Name, err)
		}
		if err == nil && !ex.Valid {
			t.Errorf("EmuInitialization: %s: Bird initialization should have raised "+
				"an error but did not", ex.Name)
		}
	}
}

func benchmarkEmuStep(querySize, numUsers, numItems int, b *testing.B) {
	usersToWeightedItems := make([]map[int]float64, numUsers)
	for i := 0; i < numUsers; i++ {
		num := 1 + rand.Intn(100) // +1 so that num != 0
		weightedItems := make(map[int]float64, num)
		for j := 0; j < num; j++ {
			it := rand.Intn(numItems)
			weightedItems[it] = 10 * rand.Float64()
		}
		usersToWeightedItems[i] = weightedItems
	}

	itemWeights := make([]float64, numItems)
	for i := 0; i < numItems; i++ {
		itemWeights[i] = 10 * rand.Float64()
	}

	bird, err := NewEmu(NewBirdCfg(), itemWeights, usersToWeightedItems)
	if err != nil {
		panic("BenchmarkEmuStep: Bird initialization raised an error " +
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

func BenchmarkEmuStep2000000Items1000Users100Query(b *testing.B) {
	benchmarkEmuStep(100, 1000, 2000000, b)
}

func BenchmarkEmuStep2000000Items100000Users100Query(b *testing.B) {
	benchmarkEmuStep(100, 100000, 2000000, b)
}

func BenchmarkEmuStep2000000Items1000000Users100Query(b *testing.B) {
	benchmarkEmuStep(100, 1000000, 2000000, b)
}

func BenchmarkEmuStep2000000Items1000000Users1000Query(b *testing.B) {
	benchmarkEmuStep(1000, 1000000, 2000000, b)
}

func BenchmarkEmuStep2000000Items1000000Users10000Query(b *testing.B) {
	benchmarkEmuStep(10000, 1000000, 2000000, b)
}

func BenchmarkEmuStep2000000Items1000000Users100000Query(b *testing.B) {
	benchmarkEmuStep(100000, 1000000, 2000000, b)
}

func BenchmarkEmuStep2000000Items1000000Users200000Query(b *testing.B) {
	benchmarkEmuStep(200000, 1000000, 2000000, b)
}

func BenchmarkEmuStep2000000Items1000000Users300000Query(b *testing.B) {
	benchmarkEmuStep(300000, 1000000, 2000000, b)
}

func BenchmarkEmuStep2000000Items1000000Users400000Query(b *testing.B) {
	benchmarkEmuStep(400000, 1000000, 2000000, b)
}

func BenchmarkEmuStep2000000Items1000000Users600000Query(b *testing.B) {
	benchmarkEmuStep(600000, 1000000, 2000000, b)
}

func BenchmarkEmuStep2000000Items1000000Users700000Query(b *testing.B) {
	benchmarkEmuStep(700000, 1000000, 2000000, b)
}

func BenchmarkEmuStep2000000Items1000000Users800000Query(b *testing.B) {
	benchmarkEmuStep(800000, 1000000, 2000000, b)
}

func BenchmarkEmuStep2000000Items1000000Users900000Query(b *testing.B) {
	benchmarkEmuStep(900000, 1000000, 2000000, b)
}

func BenchmarkEmuStep2000000Items1000000Users1000000Query(b *testing.B) {
	benchmarkEmuStep(1000000, 1000000, 2000000, b)
}

func benchmarkEmuProcess(numItems, numUsers, querySize, draws, depth int, b *testing.B) {
	itemsToUsers := make([][]int, numItems)
	for i := 0; i < numItems; i++ {
		itemsToUsers[i] = []int{1}
	}

	usersToWeightedItems := make([]map[int]float64, numUsers)
	for i := 0; i < numUsers; i++ {
		num := 1 + rand.Intn(100) // +1 so that num != 0
		weightedItems := make(map[int]float64, num)
		for j := 0; j < num; j++ {
			it := rand.Intn(numItems)
			weightedItems[it] = 10 * rand.Float64()
			itemsToUsers[it] = append(itemsToUsers[it], i)
		}
		usersToWeightedItems[i] = weightedItems
	}

	itemWeights := make([]float64, numItems)
	for i := 0; i < numItems; i++ {
		itemWeights[i] = 10 * rand.Float64()
	}

	cfg := NewBirdCfg()
	cfg.Depth = depth
	cfg.Draws = draws

	bird, err := NewEmu(cfg, itemWeights, usersToWeightedItems)
	if err != nil {
		panic("BenchmarkEmuStep: Bird initialization raised an error " +
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
func BenchmarkEmuProcess10KUsers(b *testing.B) {
	benchmarkEmuProcess(2000000, 10000, 100, 100, 1, b)
}

func BenchmarkEmuProcess100KUsers(b *testing.B) {
	benchmarkEmuProcess(2000000, 100000, 100, 100, 1, b)
}

func BenchmarkEmuProcess1MUsers(b *testing.B) {
	benchmarkEmuProcess(2000000, 1000000, 100, 100, 1, b)
}

// Increase the numbers of draws to 100k with 1M users
func BenchmarkEmuProcess100Draws(b *testing.B) {
	benchmarkEmuProcess(2000000, 1000000, 100, 100, 1, b)
}

func BenchmarkEmuProcess1KDraws(b *testing.B) {
	benchmarkEmuProcess(2000000, 1000000, 100, 1000, 1, b)
}

func BenchmarkEmuProcess10KDraws(b *testing.B) {
	benchmarkEmuProcess(2000000, 1000000, 100, 10000, 1, b)
}

func BenchmarkEmuProcess100KDraws(b *testing.B) {
	benchmarkEmuProcess(2000000, 1000000, 100, 100000, 1, b)
}

// Increase the depth up to 10 with 1M users and 10K draws
func BenchmarkEmuProcess1Depth(b *testing.B) {
	benchmarkEmuProcess(2000000, 1000000, 100, 10000, 1, b)
}

func BenchmarkEmuProcess2Depth(b *testing.B) {
	benchmarkEmuProcess(2000000, 1000000, 100, 10000, 2, b)
}

func BenchmarkEmuProcess3Depth(b *testing.B) {
	benchmarkEmuProcess(2000000, 1000000, 100, 10000, 3, b)
}

func BenchmarkEmuProcess4Depth(b *testing.B) {
	benchmarkEmuProcess(2000000, 1000000, 100, 10000, 4, b)
}

func BenchmarkEmuProcess5Depth(b *testing.B) {
	benchmarkEmuProcess(2000000, 1000000, 100, 10000, 5, b)
}

func BenchmarkEmuProcess10Depth(b *testing.B) {
	benchmarkEmuProcess(2000000, 1000000, 100, 10000, 10, b)
}

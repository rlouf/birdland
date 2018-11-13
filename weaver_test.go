package birdland

import (
	"math/rand"
	"testing"
)

type WeaverInitCase struct {
	Name         string
	ItemWeights  []float64
	UsersToItems [][]int
	SocialCoef   []map[int]float64
	Draws        int
	Depth        int
	Valid        bool
}

var weaverInitTable = []WeaverInitCase{
	{
		Name:         "Zero Depth",
		ItemWeights:  []float64{1, 1},
		UsersToItems: [][]int{[]int{0}, []int{1}},
		SocialCoef:   []map[int]float64{{0: 1.}, {1: 1.}},
		Depth:        0,
		Draws:        1,
		Valid:        false,
	},
	{
		Name:         "Negative Depth",
		ItemWeights:  []float64{1, 1},
		UsersToItems: [][]int{[]int{0}, []int{1}},
		SocialCoef:   []map[int]float64{{0: 1.}, {1: 1.}},
		Depth:        -1,
		Draws:        1,
		Valid:        false,
	},
	{
		Name:         "Zero Draws",
		ItemWeights:  []float64{1, 1},
		UsersToItems: [][]int{[]int{0}, []int{1}},
		SocialCoef:   []map[int]float64{{0: 1.}, {1: 1.}},
		Depth:        1,
		Draws:        0,
		Valid:        false,
	},
	{
		Name:         "Negative Draws",
		ItemWeights:  []float64{1, 1},
		UsersToItems: [][]int{[]int{0}, []int{1}},
		SocialCoef:   []map[int]float64{{0: 1.}, {1: 1.}},
		Depth:        1,
		Draws:        -1,
		Valid:        false,
	},
	{
		Name:         "Empty ItemWeights",
		ItemWeights:  []float64{},
		UsersToItems: [][]int{[]int{0}, []int{1}},
		SocialCoef:   []map[int]float64{{0: 1.}, {1: 1.}},
		Depth:        1,
		Draws:        1,
		Valid:        false,
	},
	{
		Name:         "Empty UsersToItems",
		ItemWeights:  []float64{1, 1},
		UsersToItems: [][]int{},
		SocialCoef:   []map[int]float64{{0: 1.}, {1: 1.}},
		Depth:        1,
		Draws:        1,
		Valid:        false,
	},
	{
		Name:         "More items in adjacency tables that weight list",
		ItemWeights:  []float64{0.1, 0.2, 0.4},
		UsersToItems: [][]int{[]int{0, 2}, []int{4}},
		SocialCoef:   []map[int]float64{{0: 1.}, {1: 1.}},
		Depth:        1,
		Draws:        1,
		Valid:        false,
	},
	{
		Name:         "Negative social weight",
		ItemWeights:  []float64{0.1, 0.4},
		UsersToItems: [][]int{[]int{0}, []int{1}},
		SocialCoef:   []map[int]float64{{0: -1.}, {1: 1.}},
		Depth:        1,
		Draws:        1,
		Valid:        false,
	},
	{
		Name:         "Discrepancy in the number of users",
		ItemWeights:  []float64{0.1, 0.2, 0.4},
		UsersToItems: [][]int{[]int{0, 1}, []int{2}, []int{1}},
		SocialCoef:   []map[int]float64{{0: 1.}, {1: 1.}},
		Depth:        1,
		Draws:        1,
		Valid:        false,
	},
	{
		Name:         "Undefined user in the social graph",
		ItemWeights:  []float64{0.1, 0.2},
		UsersToItems: [][]int{[]int{0}, []int{1}},
		SocialCoef:   []map[int]float64{{0: 1.}, {1:1., 2: 1.}},
		Depth:        1,
		Draws:        1,
		Valid:        false,
	},
	{
		Name:         "Perfectly valid input",
		ItemWeights:  []float64{1, 1},
		UsersToItems: [][]int{[]int{0}, []int{1}},
		SocialCoef:   []map[int]float64{{0: 1.}, {1: 1.}},
		Depth:        1,
		Draws:        1,
		Valid:        true,
	},
}

func TestWeaverInitialization(t *testing.T) {
	for _, ex := range weaverInitTable {
		cfg := NewBirdCfg()
		cfg.Depth = ex.Depth
		cfg.Draws = ex.Draws

		_, err := NewWeaver(cfg, ex.ItemWeights, ex.UsersToItems, ex.SocialCoef)
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

func benchmarkWeaverStep(querySize, numUsers, numItems int, b *testing.B) {
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

	var numFriends int
	if numUsers < 100 {
		numFriends = numUsers
	} else {
		numFriends = 100
	}
	socialCoef := make([]map[int]float64, numUsers)
	for i := 0; i < numUsers; i++ {
		socialCoef[i] = make(map[int]float64, 0)
		for j := 0; j < numFriends; j++ {
			socialCoef[i][j] = 10 * rand.Float64()
		}
	}

	weaver, err := NewWeaver(NewBirdCfg(), itemWeights, usersToItems, socialCoef)
	if err != nil {
		panic("BenchmarkWeaverStep: Weaver initialization raised an error " +
			"but shouldn't have. Check your test case")
	}

	query := make([]int, querySize)
	for i := 0; i < querySize; i++ {
		query[i] = rand.Intn(numItems)
	}

	user := rand.Intn(numUsers)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _ = weaver.socialStep(query, user)
	}
}

func BenchmarkWeaverStep2000000Items1000Users100Query(b *testing.B) {
	benchmarkWeaverStep(100, 1000, 2000000, b)
}

func BenchmarkWeaverStep2000000Items100000Users100Query(b *testing.B) {
	benchmarkWeaverStep(100, 100000, 2000000, b)
}

func BenchmarkWeaverStep2000000Items1000000Users100Query(b *testing.B) {
	benchmarkWeaverStep(100, 1000000, 2000000, b)
}

func BenchmarkWeaverStep2000000Items1000000Users1000Query(b *testing.B) {
	benchmarkWeaverStep(1000, 1000000, 2000000, b)
}

func BenchmarkWeaverStep2000000Items1000000Users10000Query(b *testing.B) {
	benchmarkWeaverStep(10000, 1000000, 2000000, b)
}

func BenchmarkWeaverStep2000000Items1000000Users100000Query(b *testing.B) {
	benchmarkWeaverStep(100000, 1000000, 2000000, b)
}

func BenchmarkWeaverStep2000000Items1000000Users200000Query(b *testing.B) {
	benchmarkWeaverStep(200000, 1000000, 2000000, b)
}

func BenchmarkWeaverStep2000000Items1000000Users300000Query(b *testing.B) {
	benchmarkWeaverStep(300000, 1000000, 2000000, b)
}

func BenchmarkWeaverStep2000000Items1000000Users400000Query(b *testing.B) {
	benchmarkWeaverStep(400000, 1000000, 2000000, b)
}

func BenchmarkWeaverStep2000000Items1000000Users600000Query(b *testing.B) {
	benchmarkWeaverStep(600000, 1000000, 2000000, b)
}

func BenchmarkWeaverStep2000000Items1000000Users700000Query(b *testing.B) {
	benchmarkWeaverStep(700000, 1000000, 2000000, b)
}

func BenchmarkWeaverStep2000000Items1000000Users800000Query(b *testing.B) {
	benchmarkWeaverStep(800000, 1000000, 2000000, b)
}

func BenchmarkWeaverStep2000000Items1000000Users900000Query(b *testing.B) {
	benchmarkWeaverStep(900000, 1000000, 2000000, b)
}

func BenchmarkWeaverStep2000000Items1000000Users1000000Query(b *testing.B) {
	benchmarkWeaverStep(1000000, 1000000, 2000000, b)
}

func benchmarkWeaverProcess(numItems, numUsers, querySize, draws, depth int, b *testing.B) {
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

	var numFriends int
	if numUsers < 100 {
		numFriends = numUsers
	} else {
		numFriends = 100
	}
	socialCoef := make([]map[int]float64, numUsers)
	for i := 0; i < numUsers; i++ {
		socialCoef[i] = make(map[int]float64, 0)
		for j := 0; j < numFriends; j++ {
			socialCoef[i][j] = 10 * rand.Float64()
		}
	}

	cfg := NewBirdCfg()
	cfg.Depth = depth
	cfg.Draws = draws

	weaver, err := NewWeaver(cfg, itemWeights, usersToItems, socialCoef)
	if err != nil {
		panic("BenchmarkWeaverStep: Weaver initialization raised an error " +
			"but shouldn't have. Check your test case")
	}

	query := make([]QueryItem, querySize)
	for i := 0; i < querySize; i++ {
		queryItem := QueryItem{Item: rand.Intn(numItems), Weight: rand.Float64()}
		query[i] = queryItem
	}

	user := rand.Intn(numUsers)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _ = weaver.SocialProcess(query, user)
	}
}

// Increase the numbers of users to 10M
func BenchmarkWeaverProcess10KUsers(b *testing.B) {
	benchmarkWeaverProcess(2000000, 10000, 100, 100, 1, b)
}

func BenchmarkWeaverProcess100KUsers(b *testing.B) {
	benchmarkWeaverProcess(2000000, 100000, 100, 100, 1, b)
}

func BenchmarkWeaverProcess1MUsers(b *testing.B) {
	benchmarkWeaverProcess(2000000, 1000000, 100, 100, 1, b)
}

// Increase the numbers of draws to 100k with 1M users
func BenchmarkWeaverProcess100Draws(b *testing.B) {
	benchmarkWeaverProcess(2000000, 1000000, 100, 100, 1, b)
}

func BenchmarkWeaverProcess1KDraws(b *testing.B) {
	benchmarkWeaverProcess(2000000, 1000000, 100, 1000, 1, b)
}

func BenchmarkWeaverProcess10KDraws(b *testing.B) {
	benchmarkWeaverProcess(2000000, 1000000, 100, 10000, 1, b)
}

func BenchmarkWeaverProcess100KDraws(b *testing.B) {
	benchmarkWeaverProcess(2000000, 1000000, 100, 100000, 1, b)
}

// Increase the depth up to 10 with 1M users and 10K draws
func BenchmarkWeaverProcess1Depth(b *testing.B) {
	benchmarkWeaverProcess(2000000, 1000000, 100, 10000, 1, b)
}

func BenchmarkWeaverProcess2Depth(b *testing.B) {
	benchmarkWeaverProcess(2000000, 1000000, 100, 10000, 2, b)
}

func BenchmarkWeaverProcess3Depth(b *testing.B) {
	benchmarkWeaverProcess(2000000, 1000000, 100, 10000, 3, b)
}

func BenchmarkWeaverProcess4Depth(b *testing.B) {
	benchmarkWeaverProcess(2000000, 1000000, 100, 10000, 4, b)
}

func BenchmarkWeaverProcess5Depth(b *testing.B) {
	benchmarkWeaverProcess(2000000, 1000000, 100, 10000, 5, b)
}

func BenchmarkWeaverProcess10Depth(b *testing.B) {
	benchmarkWeaverProcess(2000000, 1000000, 100, 10000, 10, b)
}

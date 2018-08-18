package sampler

import (
	"math/rand"
	"testing"
)

type AccumulateCase struct {
	Name   string
	Values []float64
	Acc    []float64
	Valid  bool
}

type TowerSamplerCase struct {
	Name       string
	NumSamples int
	Weights    []float64
	Samples    []int
	Valid      bool
}

var accumulate_table = []AccumulateCase{
	{
		Name:   "Negative weight",
		Values: []float64{0, -1, 0},
		Acc:    []float64{},
		Valid:  false,
	},
	{
		Name:   "Zero length",
		Values: []float64{},
		Acc:    []float64{},
		Valid:  false,
	},
	{
		Name:   "Zero weights",
		Values: []float64{0, 0, 0, 0},
		Acc:    []float64{},
		Valid:  false,
	},
	{
		Name:   "Uniform weights",
		Values: []float64{1, 1, 1, 1},
		Acc:    []float64{0.25, 0.5, 0.75, 1},
		Valid:  true,
	},
	{
		Name:   "Uniform weights",
		Values: []float64{1, 1, 1, 1},
		Acc:    []float64{0.25, 0.5, 0.75, 1},
		Valid:  true,
	},
	{
		Name:   "All but two similar weights",
		Values: []float64{1, 2, 2, 5},
		Acc:    []float64{0.10, 0.3, 0.5, 1},
		Valid:  true,
	},
	{
		Name:   "One weight is zero",
		Values: []float64{1, 0, 4, 5},
		Acc:    []float64{0.10, 0.10, 0.5, 1},
		Valid:  true,
	},
}

var towersampler_table = []TowerSamplerCase{
	{
		Name:       "Zero length",
		NumSamples: 0,
		Weights:    []float64{},
		Samples:    []int{},
		Valid:      false,
	},
	{
		Name:       "Zero samples",
		NumSamples: 0,
		Weights:    []float64{1, 1, 1},
		Samples:    []int{},
		Valid:      true,
	},
	{
		Name:       "One weight is zero",
		NumSamples: 10,
		Weights:    []float64{1, 0, 1},
		Samples:    []int{0, 0, 2, 0, 0, 0, 2, 0, 0, 2}, // can only be 0 or 2
		Valid:      true,
	},
	{
		Name:       "All but first weight are zero",
		NumSamples: 10,
		Weights:    []float64{1, 0, 0, 0, 0},
		Samples:    []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		Valid:      true,
	},
	{
		Name:       "All but middle weight are zero",
		NumSamples: 10,
		Weights:    []float64{0, 0, 1, 0, 0},
		Samples:    []int{2, 2, 2, 2, 2, 2, 2, 2, 2, 2},
		Valid:      true,
	},
	{
		Name:       "All but last weight are zero",
		NumSamples: 10,
		Weights:    []float64{0, 0, 0, 0, 1},
		Samples:    []int{4, 4, 4, 4, 4, 4, 4, 4, 4, 4},
		Valid:      true,
	},
	{
		Name:       "Uniform weights",
		NumSamples: 10,
		Weights:    []float64{1, 1, 1, 1, 1},
		Samples:    []int{1, 0, 3, 1, 0, 1, 4, 1, 1, 3},
		Valid:      true,
	},
	{
		Name:       "Just some standard sampling",
		NumSamples: 10,
		Weights:    []float64{2, 3, 5}, // 0.2 | 0.5 | 1
		Samples:    []int{1, 0, 2, 1, 0, 1, 2, 1, 1, 2},
		Valid:      true,
	},
}

func TestAccumulate(t *testing.T) {
	for _, ex := range accumulate_table {
		cum, err := accumulate(ex.Values)
		if err != nil {
			if ex.Valid {
				t.Errorf(`accumulate: %s should not have raised an error,
						raised  %v instead`, ex.Name, err)
			}
		} else {
			if !ex.Valid {
				t.Errorf(`accumulate: %s should have raised an error,
						got none instead`, ex.Name)
			}
		}
		for i, c := range cum {
			if c != ex.Acc[i] {
				t.Errorf(`accumulate: %s does not return the correct slice:
						expected %v, got %v`, ex.Name, ex.Acc, cum)
			}
		}
	}
}

func TestSampling(t *testing.T) {
	for _, ex := range towersampler_table {
		r := rand.New(rand.NewSource(42))
		ts, err := NewTowerSampler(r, ex.Weights)
		if err != nil {
			if ex.Valid {
				t.Errorf(`tower sampler: init: %s should not have raised an error,
						raised  %v instead`, ex.Name, err)
			}
		} else {
			if !ex.Valid {
				t.Errorf(`tower sampler: init: %s should have raised an error,
						got none instead`, ex.Name)
			}
		}

		samples := ts.Sample(ex.NumSamples)
		if len(samples) != ex.NumSamples {
			t.Errorf("tower sampler: init: %s: expected %v samples, got %v instead",
				ex.Name, ex.NumSamples, len(samples))
			continue
		}
		for i, s := range samples {
			if s != ex.Samples[i] {
				t.Errorf(`tower sampler: sample: %s does not return the correct slice:
						expected %v, got %v`, ex.Name, ex.Samples, samples)
				continue
			}
		}
	}
}

// Benchmarks
// ////////////////////////////////////////////////////////////////////////////

func initWeightsForBenchmarks(numWeights int) []float64 {
	weights := make([]float64, numWeights)
	for j := 0; j < numWeights; j++ {
		weights[j] = rand.Float64()
	}

	return weights
}

func benchmarkTowerSamplerInit(numWeights int, b *testing.B) {

	b.StopTimer()
	weights := initWeightsForBenchmarks(numWeights)
	r := rand.New(rand.NewSource(42))
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		_, _ = NewTowerSampler(r, weights)
	}
}

func BenchmarkTowerSamplerInit100(b *testing.B)     { benchmarkTowerSamplerInit(100, b) }
func BenchmarkTowerSamplerInit1000(b *testing.B)    { benchmarkTowerSamplerInit(1000, b) }
func BenchmarkTowerSamplerInit10000(b *testing.B)   { benchmarkTowerSamplerInit(10000, b) }
func BenchmarkTowerSamplerInit100000(b *testing.B)  { benchmarkTowerSamplerInit(100000, b) }
func BenchmarkTowerSamplerInit1000000(b *testing.B) { benchmarkTowerSamplerInit(1000000, b) }

func benchmarkTowerSamplerSampling(numWeights int, numSamples int, b *testing.B) {

	b.StopTimer()
	weights := initWeightsForBenchmarks(numWeights)
	r := rand.New(rand.NewSource(42))
	ts, _ := NewTowerSampler(r, weights)
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		_ = ts.Sample(numSamples)
	}
}

func BenchmarkTowerSamplerSampling100(b *testing.B) {
	benchmarkTowerSamplerSampling(10000, 100, b)
}
func BenchmarkTowerSamplerSampling1000(b *testing.B) {
	benchmarkTowerSamplerSampling(10000, 1000, b)
}
func BenchmarkTowerSamplerSampling10000(b *testing.B) {
	benchmarkTowerSamplerSampling(10000, 10000, b)
}
func BenchmarkTowerSamplerSampling100000(b *testing.B) {
	benchmarkTowerSamplerSampling(10000, 100000, b)
}
func BenchmarkTowerSamplerSampling1000000(b *testing.B) {
	benchmarkTowerSamplerSampling(10000, 1000000, b)
}

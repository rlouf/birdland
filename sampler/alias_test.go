package sampler

import (
	"math/rand"

	"testing"
)

type NormalizeCase struct {
	Name   string
	Values []float64
	Acc    []float64
	Valid  bool
}

type AliasSamplerCase struct {
	Name       string
	NumSamples int
	Weights    []float64
	Samples    []int
	Valid      bool
}

var normalize_table = []NormalizeCase{
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

var aliassampler_table = []AliasSamplerCase{
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
		Samples:    []int{2, 2, 0, 0, 2, 0, 0, 2, 2, 2}, // can only be 0 or 2
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
		Samples:    []int{0, 3, 3, 2, 3, 4, 0, 2, 2, 0},
		Valid:      true,
	},
	{
		Name:       "Just some standard sampling",
		NumSamples: 10,
		Weights:    []float64{2, 3, 5},
		Samples:    []int{2, 2, 1, 0, 2, 1, 0, 2, 2, 2},
		Valid:      true,
	},
}

func TestAliasSampling(t *testing.T) {
	for _, ex := range aliassampler_table {
		r := rand.New(rand.NewSource(42))
		ts := AliasSampler{}
		err := ts.Init(r, ex.Weights)
		if err != nil {
			if ex.Valid {
				t.Errorf("tower sampler: init: %s should not have raised an error, raised  %v instead", ex.Name, err)
			}
		} else {
			if !ex.Valid {
				t.Errorf("tower sampler: init: %s should have raised an error, got none instead", ex.Name)
			}
		}

		samples := ts.Sample(ex.NumSamples)
		if len(samples) != ex.NumSamples {
			t.Errorf("tower sampler: init: %s: expected %v samples, got %v instead", ex.Name, ex.NumSamples, len(samples))
			continue
		}
		for i, s := range samples {
			if s != ex.Samples[i] {
				t.Errorf("tower sampler: sample: %s does not return the correct slice: expected %v, got %v", ex.Name, ex.Samples, samples)
				break
			}
		}
	}
}

func initWeightsForAliasBenchmarks(numWeights int) []float64 {
	weights := make([]float64, numWeights)
	for j := 0; j < numWeights; j++ {
		weights[j] = rand.Float64()
	}

	return weights
}

func benchmarkAliasSamplerInit(numWeights int, b *testing.B) {

	b.StopTimer()
	weights := initWeightsForAliasBenchmarks(numWeights)
	r := rand.New(rand.NewSource(42))
	ts := AliasSampler{}
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		_ = ts.Init(r, weights)
	}
}

func BenchmarkAliasSamplerInit1000(b *testing.B)    { benchmarkAliasSamplerInit(1000, b) }
func BenchmarkAliasSamplerInit10000(b *testing.B)   { benchmarkAliasSamplerInit(10000, b) }
func BenchmarkAliasSamplerInit100000(b *testing.B)  { benchmarkAliasSamplerInit(100000, b) }
func BenchmarkAliasSamplerInit1000000(b *testing.B) { benchmarkAliasSamplerInit(1000000, b) }

func benchmarkAliasSamplerSampling(numWeights int, numSamples int, b *testing.B) {

	b.StopTimer()
	weights := initWeightsForAliasBenchmarks(numWeights)
	r := rand.New(rand.NewSource(42))
	ts := AliasSampler{}
	_ = ts.Init(r, weights)
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		_ = ts.Sample(numSamples)
	}
}

func BenchmarkAliasSamplerSampling1000(b *testing.B)  { benchmarkAliasSamplerSampling(10000, 1000, b) }
func BenchmarkAliasSamplerSampling10000(b *testing.B) { benchmarkAliasSamplerSampling(10000, 10000, b) }
func BenchmarkAliasSamplerSampling100000(b *testing.B) {
	benchmarkAliasSamplerSampling(10000, 100000, b)
}
func BenchmarkAliasSamplerSampling1000000(b *testing.B) {
	benchmarkAliasSamplerSampling(10000, 1000000, b)
}

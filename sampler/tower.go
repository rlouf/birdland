package sampler

import (
	"fmt"
	"math/rand"
	"sort"

	"github.com/pkg/errors"
)

// TowerSampler implements the Tower Sampling algorithm used to sample from a discrete
// probability distribution.
type TowerSampler struct {
	CumulativeSum []float64
	Source        *rand.Rand
}

func (t *TowerSampler) Init(source *rand.Rand, weights []float64) error {

	if len(weights) == 0 {
		return fmt.Errorf("weights is an empty slice")
	}

	cumulative, err := accumulate(weights)
	if err != nil {
		return errors.Wrap(err, "cannot initialize the tower sampler")
	}

	t.CumulativeSum = cumulative
	t.Source = source

	return nil
}

func (t *TowerSampler) Sample(numSamples int) []int {
	samples := make([]int, numSamples)
	for i := 0; i < numSamples; i++ {
		x := t.Source.Float64()
		sample := sort.Search(len(t.CumulativeSum), func(j int) bool { return t.CumulativeSum[j] >= x })
		samples[i] = sample
	}

	return samples
}

// accumulate computes the cumulative sum of a slice normalized by
// the sum of all terms.
func accumulate(weights []float64) ([]float64, error) {
	var sum float64
	cumulativeSum := make([]float64, len(weights))
	for i, weight := range weights {
		if weight < 0 {
			return nil, fmt.Errorf("negative weight: %g", weight)
		}
		sum += weight
		cumulativeSum[i] = sum
	}

	if sum == 0 {
		return nil, fmt.Errorf("all weights are null")
	}

	for i, cumSum := range cumulativeSum {
		cumulativeSum[i] = cumSum / sum
	}

	return cumulativeSum, nil
}

package sampler

import (
	"fmt"
	"math/rand"

	"github.com/pkg/errors"
)

// AliasSampler implements the Alias Method to sample from a discrete
// probability distribution. Initialized with the Vose Method, the
// sampler takes O(n) to initialize and O(1) to sample.
type AliasSampler struct {
	ProbabilityTable []float64
	AliasTable       []int
	Source           *rand.Rand
}

// Init initalizes the sampler
func (t *AliasSampler) Init(source *rand.Rand, weights []float64) error {

	if len(weights) == 0 {
		return fmt.Errorf("weights is an empty slice")
	}

	ProbabilityTable, AliasTable, err := VoseInitialization(weights)
	if err != nil {
		return errors.Wrap(err, "cannot initialize the alias sampler")
	}
	t.ProbabilityTable = ProbabilityTable
	t.AliasTable = AliasTable
	t.Source = source

	return nil
}

// Sample generates a slice of items obtained by sampling the original distribution.
func (t *AliasSampler) Sample(numSamples int) []int {
	n := len(t.AliasTable)
	samples := make([]int, numSamples)

	if n == 0 {
		return samples
	}

	for i := 0; i < numSamples; i++ {
		k := t.Source.Intn(n)
		toss := t.Source.Float64()
		if toss < t.ProbabilityTable[k] {
			samples[i] = k
		} else {
			samples[i] = t.AliasTable[k]
		}
	}

	return samples
}

// VoseInitialization initialises the probability and alias tables using Vose's
// method. Vose's method runs in O(n) and is more numerically stable than
// alternatives. See http://www.keithschwarz.com/darts-dice-coins/ for more
// details.
func VoseInitialization(weights []float64) ([]float64, []int, error) {

	normalizedWeights, err := normalize(weights)
	if err != nil {
		return []float64{}, []int{}, errors.Wrap(err, "cannot normalize input weights")
	}

	small := make([]int, 0, len(normalizedWeights))
	large := make([]int, 0, len(normalizedWeights))
	for i, w := range normalizedWeights {
		if w < 1.0 {
			small = append(small, i)
		} else {
			large = append(large, i)
		}
	}

	AliasTable := make([]int, len(weights))
	ProbabilityTable := make([]float64, len(weights))
	var g, l int
	for (len(small) > 0) && (len(large) > 0) {
		l, small = small[0], small[1:]
		g, large = large[0], large[1:]

		AliasTable[l] = g
		ProbabilityTable[l] = normalizedWeights[l]

		normalizedWeights[g] = (normalizedWeights[g] + normalizedWeights[l]) - 1
		if normalizedWeights[g] < 1.0 {
			small = append(small, g)
		} else {
			large = append(large, g)
		}
	}

	for len(large) > 0 {
		g, large = large[0], large[1:]
		ProbabilityTable[g] = 1
	}
	for len(small) > 0 {
		l, small = small[0], small[1:]
		ProbabilityTable[g] = 1
	}

	return ProbabilityTable, AliasTable, nil
}

// normalize prepares the weights for the algorithm's initialization.
func normalize(weights []float64) ([]float64, error) {
	var sum float64
	n := len(weights)
	for _, w := range weights {
		if w < 0 {
			return []float64{}, fmt.Errorf("found negative weight %v", w)
		}
		sum += w
	}

	normalizedWeights := make([]float64, len(weights))
	for i, weight := range weights {
		normalizedWeights[i] = float64(n) * weight / sum
	}

	return normalizedWeights, nil
}

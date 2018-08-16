// Package sampler provides primitives for drawing samples from slices and arrays
// providing a list of weights.
//
// A sampler takes a list of weight as an input and returns a list of ids of
// the elements that were chosen.
//
// We deliberately choose to not take the list of elements as an input and defer the choice
// to the module that uses the sampler. Indeed, some applications may require to sample from
// an array of strings, while another from a slice of ints.
package sampler

type Sampler interface {
	Init([]float64) error
	Sample(int) []int
}

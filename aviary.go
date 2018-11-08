package birdland

import "github.com/pkg/errors"

// aviary is a simple mixture of Birds
type Aviary struct {
	Birds []Bird
}

func (a *Aviary) Process(query []QueryItem) ([]int, []int, error) {

	var items []int
	var referrers []int
	for i, bird := range a.Birds {
		stepItems, stepReferrers, err := bird.Process(query)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "cannot process query from bird %d", i)
		}
		items = append(items, stepItems...)
		referrers = append(referrers, stepReferrers...)
	}

	return items, referrers, nil
}
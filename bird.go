package birdland

import (
	"log"
	"math/rand"
	"time"

	"github.com/pkg/errors"
	"github.com/vimies/birdland/sampler"
)

type QueryItem struct {
	Item   int
	Weight float64
}

type Bird struct {
	ItemWeights  []float64
	UsersToItems [][]int
	ItemsToUsers [][]int
	Rand         *rand.Rand
	Draws        int
	Depth        int
}

func NewBird(itemWeights []float64, usersToItems [][]int, itemsToUsers [][]int) *Bird {
	b := Bird{
		Depth:        1,
		Draws:        1000,
		Rand:         rand.New(rand.NewSource(42)),
		ItemWeights:  itemWeights,
		UsersToItems: usersToItems,
		ItemsToUsers: itemsToUsers,
	}
	log.Printf("initialized Bird with %d draws and depth %d", b.Draws, b.Depth)
	return &b
}

func (b *Bird) setDepth(depth int) {
	b.Depth = depth
}

func (b *Bird) setDraws(depth int) {
	b.Draws = draws
}

// Process returns a slice of recommended items along with their referrer given
// a query consisting of a slice of items with their respective weights.
func (b *Bird) Process(query []QueryItem) ([]int, []int, error) {
	start := time.Now()

	stepItems, err := b.SampleItemsFromQuery(query)
	if err != nil {
		return nil, nil, errors.Wrap(err, "cannot process the query")
	}

	var items []int
	var referrers []int
	for d := 0; d < b.Depth; d++ {
		stepItems, stepReferrers, err := b.Step(stepItems)
		if err != nil {
			return nil, nil, errors.Wrap(err, "cannot process the query")
		}
		items = append(items, stepItems...)
		referrers = append(referrers, stepReferrers...)
	}

	elapsed := time.Since(start)
	log.Printf("processed query containing %d items in %v", len(query), elapsed)

	return items, referrers, nil
}

// Step transforms a slice of items into a slice of recommended items and a
// slice containing the corresponding referrers.
func (b *Bird) Step(items []int) ([]int, []int, error) {

	referrers := make([]int, len(items))
	for i, item := range items {
		relatedUsers := b.ItemsToUsers[item]
		referrers[i] = relatedUsers[b.Rand.Intn(len(relatedUsers))]
	}

	var err error
	newItems := make([]int, len(items))
	for j, user := range referrers {
		relatedItems := b.UsersToItems[user]
		newItems[j], err = b.SampleItem(relatedItems)
		if err != nil {
			return nil, nil, errors.Wrap(err, "cannot perform a processing step")
		}
	}

	return newItems, referrers, nil
}

// sampleItemsFromQuery takes a slice of queries and returns a list of items
// that have been sampled according to their respective weights as given by the
// weights in the query and the general item weight.
func (b *Bird) SampleItemsFromQuery(query []QueryItem) ([]int, error) {

	weights := make([]float64, len(query))
	items := make([]int, len(query))
	for i, q := range query {
		weights[i] = q.Weight * b.ItemWeights[q.Item]
		items[i] = q.Item
	}

	s, err := sampler.NewAliasSampler(b.Rand, weights)
	if err != nil {
		return nil, errors.Wrap(err, "cannot sample items from the query")
	}
	sampledItems := make([]int, b.Draws)
	for i, index := range s.Sample(b.Draws) {
		sampledItems[i] = items[index]
	}

	return sampledItems, nil
}

// sampleItem returns an item id sampled from a list of items.
func (b *Bird) SampleItem(from []int) (int, error) {
	weights := make([]float64, len(from))
	for i, f := range from {
		weights[i] = b.ItemWeights[f]
	}

	s, err := sampler.NewTowerSampler(b.Rand, weights) // 50% faster than Alias Sampler
	if err != nil {
		return 0, errors.Wrap(err, "cannot sample an item")
	}
	sampledItem := from[s.Sample(1)[0]]

	return sampledItem, nil
}

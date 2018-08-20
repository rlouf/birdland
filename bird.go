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
	ItemWeights       []float64
	UsersToItems      [][]int
	ItemsToUsers      [][]int
	UserItemsSamplers []sampler.AliasSampler
	Rand              *rand.Rand
	Draws             int
	Depth             int
}

// NewBird creates a new bird
// TODO(remi) check the validity of the input pre-initialization
func NewBird(itemWeights []float64, usersToItems [][]int, itemsToUsers [][]int, options ...func(*Bird) error) (*Bird, error) {

	randSource := rand.New(rand.NewSource(42))

	userItemsSampler, err := initUserItemsSamplers(randSource, itemWeights, usersToItems)
	if err != nil {
		return &Bird{}, errors.Wrap(err, "cannot initialize Bird")
	}

	b := Bird{
		Depth:             1,
		Draws:             1000,
		Rand:              randSource,
		ItemWeights:       itemWeights,
		UsersToItems:      usersToItems,
		ItemsToUsers:      itemsToUsers,
		UserItemsSamplers: userItemsSampler,
	}

	for _, option := range options {
		err := option(&b)
		if err != nil {
			return &b, errors.Wrap(err, "cannot initialize Bird")
		}
	}
	log.Printf("initialized Bird with %d draws and depth %d", b.Draws, b.Depth)

	return &b, nil
}

func Depth(depth int) func(*Bird) error {
	return func(t *Bird) error {
		return t.setDepth(depth)
	}
}

func Draws(draws int) func(*Bird) error {
	return func(t *Bird) error {
		return t.setDraws(draws)
	}
}

func (b *Bird) setDepth(depth int) error {
	if depth < 1 {
		return errors.New("the depth needs to be greater than 1")
	}
	b.Depth = depth
	return nil
}

func (b *Bird) setDraws(draws int) error {
	if draws < 1 {
		return errors.New("you need to set at least one draw")
	}
	b.Draws = draws
	return nil
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
		newItems[j], err = b.SampleItem(user)
		if err != nil {
			return nil, nil, errors.Wrap(err, "cannot perform a processing step")
		}
	}

	return newItems, referrers, nil
}

// sampleItem returns an item id sampled from a list of items.
func (b *Bird) SampleItem(user int) (int, error) {
	s := b.UserItemsSamplers[user]
	sampledItem := b.UsersToItems[user][s.Sample(1)[0]]

	return sampledItem, nil
}

// initUserItemsSamplers initializes the samplers used to sample from the items
// a given user has interacted with.
func initUserItemsSamplers(randSource *rand.Rand, itemWeights []float64, userToItems [][]int) ([]sampler.AliasSampler, error) {
	userItemsSamplers := make([]sampler.AliasSampler, len(userToItems))
	for i, userItems := range userToItems {

		weights := make([]float64, len(userItems))
		for j, item := range userItems {
			weights[j] = itemWeights[item]
		}

		userItemsSampler, err := sampler.NewAliasSampler(randSource, weights)
		if err != nil {
			return nil, errors.Wrap(err, "could not initialize the probability and alias tables")
		}

		userItemsSamplers[i] = *userItemsSampler
	}

	return userItemsSamplers, nil
}

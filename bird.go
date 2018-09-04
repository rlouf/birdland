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
	RandSource        *rand.Rand
	Draws             int
	Depth             int
}

// NewBird creates a new recommender from the input data. The number of draws
// and depth of random walks are respectively set to 1000 and 1, but can be
// changed by passing the functional options Draws() and Depth().
func NewBird(itemWeights []float64,
	usersToItems [][]int,
	itemsToUsers [][]int,
	options ...func(*Bird) error) (*Bird, error) {

	randSource := rand.New(rand.NewSource(time.Now().UnixNano()))

	err := validateBirdInputs(itemWeights, usersToItems, itemsToUsers)
	if err != nil {
		return &Bird{}, errors.Wrap(err, "cannot initialize Bird")
	}

	userItemsSampler, err := initUserItemsSamplers(randSource, itemWeights, usersToItems)
	if err != nil {
		return &Bird{}, errors.Wrap(err, "cannot initialize Bird")
	}

	b := Bird{
		Depth:             1,
		Draws:             1000,
		RandSource:        randSource,
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
		return errors.New("the depth needs to be at least 1")
	}
	b.Depth = depth
	return nil
}

func (b *Bird) setDraws(draws int) error {
	if draws < 1 {
		return errors.New("the number of draws needs to be at least 1")
	}
	b.Draws = draws
	return nil
}

// Process returns a slice of items that were visited during the random walks
// along with the users that referred these items.
func (b *Bird) Process(query []QueryItem) ([]int, []int, error) {
	start := time.Now()

	stepItems, err := b.sampleItemsFromQuery(query)
	if err != nil {
		return nil, nil, errors.Wrap(err, "cannot process the query")
	}

	var items []int
	var referrers []int
	for d := 0; d < b.Depth; d++ {
		stepItems, stepReferrers, err := b.step(stepItems)
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

// sampleItemsFromQuery takes a slice of query items and samples b.Draws items
// from it. Each item i is assigned a weight Q_i in the query---the number of
// listens, likes, shares, etc. The weight W_i used for sampleItemFromQuery is
// the product of Q_i with the item's weight I_i provided when Bird is
// initialized.
// sampleItemsFromQuery returns a slice of items that are the starting points
// of the subsequent random walks.
func (b *Bird) sampleItemsFromQuery(query []QueryItem) ([]int, error) {

	weights := make([]float64, len(query))
	items := make([]int, len(query))
	for i, q := range query {
		weights[i] = q.Weight * b.ItemWeights[q.Item]
		items[i] = q.Item
	}
	s, err := sampler.NewAliasSampler(b.RandSource, weights)
	if err != nil {
		return nil, errors.Wrap(err, "cannot sample items from the query")
	}

	sampledItems := make([]int, b.Draws)
	for i, index := range s.Sample(b.Draws) {
		sampledItems[i] = items[index]
	}

	return sampledItems, nil
}

// step performs one random walk step for each incoming item.
// step returns a slice of visited items along with the 'referrers', i.e. the
// users that were visited to reach these items.
func (b *Bird) step(items []int) ([]int, []int, error) {

	referrers := make([]int, len(items))
	for i, item := range items {
		relatedUsers := b.ItemsToUsers[item]
		referrers[i] = relatedUsers[b.RandSource.Intn(len(relatedUsers))]
	}

	var err error
	newItems := make([]int, len(items))
	for j, user := range referrers {
		newItems[j], err = b.sampleItem(user)
		if err != nil {
			return nil, nil, errors.Wrap(err, "cannot perform a processing step")
		}
	}

	return newItems, referrers, nil
}

// sampleItem returns an item sampled from a user's collection.
func (b *Bird) sampleItem(user int) (int, error) {
	s := b.UserItemsSamplers[user]
	sampledItem := b.UsersToItems[user][s.Sample(1)[0]]

	return sampledItem, nil
}

// initUserItemsSamplers initializes the samplers used to sample from a user's
// item collection. We use the alias sampling method which has proven sensibly
// better in benchmarks.
func initUserItemsSamplers(randSource *rand.Rand,
	itemWeights []float64,
	userToItems [][]int) ([]sampler.AliasSampler, error) {

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

// validateBirdInput checks the validity of the data fed to Bird.  It returns
// an error when it identifies a discrepancy that could make the processing
// algorithm crash.
// TODO(remi) check that userToItems and itemsToUsers are consistent.
func validateBirdInputs(itemWeights []float64,
	usersToItems [][]int,
	itemsToUsers [][]int) error {

	if len(itemWeights) == 0 {
		return errors.New("empty slice of item weights")
	}
	if len(usersToItems) == 0 {
		return errors.New("empty users to items adjacency table")
	}
	if len(itemsToUsers) == 0 {
		return errors.New("empty items to users adjacency table")
	}

	// Check that there is a weight for each item present in adjacency tables.
	numItems := len(itemWeights)
	var m int
	for _, userItems := range usersToItems {
		for _, item := range userItems {
			if item > m {
				m = item
			}
		}
	}
	if numItems < len(itemsToUsers) {
		return errors.New("there are more items in ItemsToUsers than there are weights")
	}
	if numItems < m {
		return errors.New("there are more items in UsersToItems than there are weights")
	}

	return nil
}

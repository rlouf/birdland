package birdland

import (
	"math/rand"
	"time"

	"github.com/pkg/errors"
	"github.com/vimies/birdland/sampler"
)

type QueryItem struct {
	Item   int
	Weight float64
}

type BirdCfg struct {
	Depth int `yaml:"depth"`
	Draws int `yaml:"draws"`
}

func NewBirdCfg() *BirdCfg {
	cfg := BirdCfg{
		Depth: 1,
		Draws: 1000,
	}

	return &cfg
}

type Bird struct {
	Cfg               *BirdCfg
	ItemWeights       []float64
	UsersToItems      [][]int
	ItemsToUsers      [][]int
	UserItemsSamplers []sampler.AliasSampler
	RandSource        *rand.Rand
}

// NewBird creates a new recommender from input data.
func NewBird(cfg *BirdCfg, itemWeights []float64, usersToItems [][]int) (*Bird, error) {
	if cfg.Depth < 1 {
		return nil, errors.New("depth must be greater or equal to 1")
	}

	if cfg.Draws < 1 {
		return nil, errors.New("number of draws must be greater or equal to 1")
	}

	randSource := rand.New(rand.NewSource(time.Now().UnixNano()))

	err := validateBirdInputs(itemWeights, usersToItems)
	if err != nil {
		return &Bird{}, errors.Wrap(err, "invalid input")
	}

	userItemsSampler, err := initUserItemsSamplers(randSource, itemWeights, usersToItems)
	if err != nil {
		return &Bird{}, errors.Wrap(err, "cannot initialize samplers")
	}

	itemsToUsers := permuteAdjacencyList(len(itemWeights), usersToItems)

	b := Bird{
		Cfg:               cfg,
		RandSource:        randSource,
		ItemWeights:       itemWeights,
		UsersToItems:      usersToItems,
		ItemsToUsers:      itemsToUsers,
		UserItemsSamplers: userItemsSampler,
	}

	return &b, nil
}

// Process returns a slice of items that were visited during the random walks
// along with the users that referred these items.
func (b *Bird) Process(query []QueryItem) ([]int, []int, error) {
	if len(query) == 0 {
		return nil, nil, errors.New("empty query")
	}

	stepItems, err := b.sampleItemsFromQuery(query)
	if err != nil {
		return nil, nil, errors.Wrap(err, "cannot sample items")
	}

	var items []int
	var referrers []int
	for d := 0; d < b.Cfg.Depth; d++ {
		stepItems, stepReferrers, err := b.step(stepItems)
		if err != nil {
			return nil, nil, errors.Wrap(err, "cannot step through items")
		}
		items = append(items, stepItems...)
		referrers = append(referrers, stepReferrers...)
	}

	return items, referrers, nil
}

// sampleItemsFromQuery takes a slice of query items and samples b.Cfg.Draws
// items from it. Each item i is assigned a weight Q_i in the query---the
// number of listens, likes, shares, etc. The weight W_i used for
// sampleItemFromQuery is the product of Q_i with the item's weight I_i
// provided when Bird is initialized.
//
// sampleItemsFromQuery returns a slice of items that are the starting points
// of the subsequent random walks. If the query refers to an item that has no
// record in ItemsToUsers, the item is ignored.
func (b *Bird) sampleItemsFromQuery(query []QueryItem) ([]int, error) {

	weights := make([]float64, len(query))
	items := make([]int, len(query))
	for i, q := range query {
		weights[i] = q.Weight * b.ItemWeights[q.Item]
		items[i] = q.Item
	}
	s, err := sampler.NewAliasSampler(b.RandSource, weights)
	if err != nil {
		return nil, errors.Wrap(err, "cannot create sampler")
	}

	sampledItems := make([]int, b.Cfg.Draws)
	for i, iid := range s.Sample(b.Cfg.Draws) {
		if len(b.ItemsToUsers[iid]) == 0 {
			continue
		}
		sampledItems[i] = items[iid]
	}

	if len(sampledItems) == 0 {
		return nil, errors.New("no items left after sampling," +
			"check that the query refers to actual items.")
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
		if len(relatedUsers) == 0 {
			return nil, nil, errors.New("the item refers to an item no one has interacted with")
		}
		referrers[i] = relatedUsers[b.RandSource.Intn(len(relatedUsers))]
	}

	newItems := make([]int, len(items))
	for j, user := range referrers {
		newItems[j] = b.sampleItem(user)
	}

	return newItems, referrers, nil
}

// sampleItem returns an item sampled from a user's collection.
func (b *Bird) sampleItem(user int) int {
	s := b.UserItemsSamplers[user]
	sampledItem := b.UsersToItems[user][s.Sample(1)[0]]

	return sampledItem
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
func validateBirdInputs(itemWeights []float64, usersToItems [][]int) error {

	if len(itemWeights) == 0 {
		return errors.New("empty slice of item weights")
	}
	if len(usersToItems) == 0 {
		return errors.New("empty users to items adjacency table")
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
	if numItems < m {
		return errors.New("there are more items in UsersToItems than there are weights")
	}

	return nil
}

// permuteAdjacencyList transforms the UsersToItems adjacency list into the complementary
// ItemsToUsers adjacency list.
func permuteAdjacencyList(numItems int, usersToItems [][]int) [][]int {

	itemsToUsers := make([][]int, numItems)
	for iid := 0; iid < numItems; iid++ {
		itemsToUsers[iid] = make([]int, 0)
	}

	for uid, userItems := range usersToItems {
		for _, iid := range userItems {
			itemsToUsers[iid] = append(itemsToUsers[iid], uid)
		}
	}

	return itemsToUsers
}

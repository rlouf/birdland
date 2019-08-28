package birdland

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/rlouf/birdland/sampler"
)

type WeaverCfg struct {
	DefaultWeight float64 `json:"default_weight"`
	*BirdCfg
}

func NewWeaverCfg() *WeaverCfg {
	cfg := WeaverCfg{
		DefaultWeight: 1,
		BirdCfg:       NewBirdCfg(),
	}

	return &cfg
}

// Weaver is a combination of a Bird, and a list of maps that represents
// a weighted user-user matrix.
// To avoid storing the full social graph, a mostly empty matrix, we store
// for each user a map that associates each connection to a weight. Each
// user that is not connected is attributed a DefaultWeight.
type Weaver struct {
	Cfg         *WeaverCfg
	SocialGraph []map[int]float64
	*Bird
}

// NewWeaver creates a new recommender from input data.
// Unlike Bird, users related to an item are not sampled uniformly, but according to socialCoef[user],
// making the recommendation dependent on both the query and the user being served.
func NewWeaver(cfg *WeaverCfg, itemWeights []float64, usersToItems [][]int,
	socialGraph []map[int]float64) (*Weaver, error) {

	err := validateWeaverInputs(itemWeights, usersToItems, socialGraph)
	if err != nil {
		return &Weaver{}, errors.Wrap(err, "invalid input")
	}

	bird, err := NewBird(cfg.BirdCfg, itemWeights, usersToItems)
	if err != nil {
		return &Weaver{}, errors.Wrap(err, "couldn't create new bird")
	}

	b := Weaver{
		cfg,
		socialGraph,
		bird,
	}

	return &b, nil
}

// Process returns a slice of items that were visited during the random walks
// along with the users that referred these items.
func (b *Weaver) Process(query []QueryItem, user int) ([]int, []int, error) {
	if len(query) == 0 {
		return nil, nil, errors.New("the input query is empty")
	}

	stepItems, err := b.sampleItemsFromQuery(query)
	if err != nil {
		return nil, nil, errors.Wrap(err, "cannot sample items from the query")
	}

	var items []int
	var referrers []int
	for d := 0; d < b.Cfg.Depth; d++ {
		stepItems, stepReferrers, err := b.step(stepItems, user)
		if err != nil {
			return nil, nil, errors.Wrap(err, "cannot step through items")
		}
		items = append(items, stepItems...)
		referrers = append(referrers, stepReferrers...)
	}

	return items, referrers, nil
}

// step performs one random walk step for each incoming item.
// it returns a slice of visited items along with the 'referrers', i.e. the
// users that were visited to reach these items.
func (b *Weaver) step(items []int, user int) ([]int, []int, error) {

	if user >= len(b.SocialGraph) {
		return nil, nil, fmt.Errorf("user %d does not belong to the social graph", user)
	}

	referrers := make([]int, len(items))
	itemUserSamplers := make(map[int]*sampler.AliasSampler)

	for i, item := range items {
		relatedUsers := b.ItemsToUsers[item]

		if len(relatedUsers) == 0 {
			return nil, nil, errors.New("the item refers to an item no one has interacted with")
		}

		// for each item, create a sampler of related users weighted by socialCoef
		// (with default weight value 1)
		if _, ok := itemUserSamplers[item]; !ok {
			weightedRelatedUsers := make([]float64, len(relatedUsers))
			for j, u := range relatedUsers {
				if w, ok := b.SocialGraph[user][u]; ok {
					weightedRelatedUsers[j] = w
				} else {
					weightedRelatedUsers[j] = b.Cfg.DefaultWeight
				}
			}
			itemUserSampler, err := sampler.NewAliasSampler(b.RandSource, weightedRelatedUsers)
			itemUserSamplers[item] = itemUserSampler
			if err != nil {
				return nil, nil, errors.Wrapf(err, "could not initialize users' sampler for user %d and item %d", user, item)
			}
		}
		referrers[i] = relatedUsers[itemUserSamplers[item].Sample(1)[0]]
	}

	newItems := make([]int, len(items))
	for j, user := range referrers {
		newItems[j] = b.sampleItem(user)
	}

	return newItems, referrers, nil
}

// validateWeaverInput checks the validity of the data fed to Weaver.  It returns
// an error when it identifies a discrepancy that could make the processing
// algorithm crash.
func validateWeaverInputs(itemWeights []float64, usersToItems [][]int, socialGraph []map[int]float64) error {

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
	if numItems <= m {
		return fmt.Errorf("there are more items (%d) in UsersToItems than there are weights (%d)", m, numItems)
	}

	if len(socialGraph) != len(usersToItems) {
		return errors.New("UsersToItems and the social graph don't contain the same number of users")
	}

	numUsers := len(socialGraph)
	m = 0
	for _, friendsCoef := range socialGraph {
		for user, w := range friendsCoef {
			if user > m {
				m = user
			}
			if w < 0 {
				return errors.New("weights in the social graph must be positive")
			}
		}
	}
	if numUsers <= m {
		return errors.New("some users mentioned in the connections are otherwise absent from the graph")
	}

	return nil
}

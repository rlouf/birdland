package birdland

import (
	"github.com/pkg/errors"
	"github.com/rlouf/birdland/sampler"
)

type Weaver struct {
	SocialCoef []map[int]float64
	*Bird
}

// NewWeaver creates a new recommender from input data.
// Unlike Bird, users related to an item are not sampled uniformly, but according to socialCoef[user],
// making the recommendation dependent on both the query and the user being served.
func NewWeaver(cfg *BirdCfg, itemWeights []float64, usersToItems [][]int,
	socialCoef []map[int]float64) (*Weaver, error) {

	err := validateWeaverInputs(itemWeights, usersToItems, socialCoef)
	if err != nil {
		return &Weaver{}, errors.Wrap(err, "invalid input")
	}

	bird, err := NewBird(cfg, itemWeights, usersToItems)
	if err != nil {
		return &Weaver{}, errors.Wrap(err, "couldn't create new bird")
	}

	b := Weaver{
		socialCoef,
		bird,
	}

	return &b, nil
}

// SocialProcess returns a slice of items that were visited during the random walks
// along with the users that referred these items.
func (b *Weaver) SocialProcess(query []QueryItem, user int) ([]int, []int, error) {
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
		stepItems, stepReferrers, err := b.socialStep(stepItems, user)
		if err != nil {
			return nil, nil, errors.Wrap(err, "cannot step through items")
		}
		items = append(items, stepItems...)
		referrers = append(referrers, stepReferrers...)
	}

	return items, referrers, nil
}

// socialStep performs one random walk step for each incoming item.
// socialStep returns a slice of visited items along with the 'referrers', i.e. the
// users that were visited to reach these items.
func (b *Weaver) socialStep(items []int, user int) ([]int, []int, error) {

	if user >= len(b.SocialCoef) {
		return nil, nil, errors.New("the user does not belong to the social graph")
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
				if w, ok := b.SocialCoef[user][u]; ok {
					weightedRelatedUsers[j] = w
				} else {
					weightedRelatedUsers[j] = 1.
				}
			}
			itemUserSampler, err := sampler.NewAliasSampler(b.RandSource, weightedRelatedUsers)
			itemUserSamplers[item] = itemUserSampler
			if err != nil {
				return nil, nil, errors.Wrapf(err, "couldn't initialize users' sampler for user %d and item %d", user, item)
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
func validateWeaverInputs(itemWeights []float64, usersToItems [][]int, socialCoef []map[int]float64) error {

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
		return errors.New("there are more items in UsersToItems than there are weights")
	}

	if len(socialCoef) != len(usersToItems) {
		return errors.New("UsersToItems and the social graph don't contain the same number of users")
	}

	numUsers := len(socialCoef)
	m = 0
	for _, friendsCoef := range socialCoef {
		for user, w := range friendsCoef {
			if user > m {
				m = user
			}
			if w < 0 {
				return errors.New("negative weight in the social graph")
			}
		}
	}
	if numUsers <= m {
		return errors.New("there are undefined users in the social graph")
	}

	return nil
}

package birdland

import (
	"math/rand"
	"time"

	"github.com/pkg/errors"
	"github.com/vimies/birdland/sampler"
)

// NewEmu creates a new recommender from input data where the users-to-items graph is weighted.
func NewEmu(cfg *BirdCfg, itemWeights []float64, usersToWeightedItems []map[int]float64) (*Bird, error) {
	if cfg.Depth < 1 {
		return nil, errors.New("depth must be greater or equal to 1")
	}

	if cfg.Draws < 1 {
		return nil, errors.New("number of draws must be greater or equal to 1")
	}

	randSource := rand.New(rand.NewSource(time.Now().UnixNano()))

	err := validateEmuInputs(itemWeights, usersToWeightedItems)
	if err != nil {
		return &Bird{}, errors.Wrap(err, "invalid input")
	}

	userItemsSampler, usersToItems, err := initUserWeightedItemsSamplers(randSource, usersToWeightedItems)
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

// initUserItemsSamplers initializes the samplers used to sample from a user's
// item collection. We use the alias sampling method which has proven sensibly
// better in benchmarks.
func initUserWeightedItemsSamplers(randSource *rand.Rand,
	usersToWeightedItems []map[int]float64) ([]sampler.AliasSampler, [][]int, error) {

	usersToItems := make([][]int, len(usersToWeightedItems))
	userItemsSamplers := make([]sampler.AliasSampler, len(usersToWeightedItems))
	for i, userItems := range usersToWeightedItems {

		usersToItems[i] = make([]int, len(userItems))
		weights := make([]float64, len(userItems))
		j := 0
		for item, w := range userItems {
			usersToItems[i][j] = item
			weights[j] = w
			j++
		}

		userItemsSampler, err := sampler.NewAliasSampler(randSource, weights)
		if err != nil {
			return nil, nil, errors.Wrap(err, "could not initialize the probability and alias tables")
		}

		userItemsSamplers[i] = *userItemsSampler
	}

	return userItemsSamplers, usersToItems, nil
}

// validateEmuInput checks the validity of the data fed to a weighted Bird.  It returns
// an error when it identifies a discrepancy that could make the processing
// algorithm crash.
func validateEmuInputs(itemWeights []float64, usersToWeightedItems []map[int]float64) error {

	if len(itemWeights) == 0 {
		return errors.New("empty slice of item weights")
	}
	if len(usersToWeightedItems) == 0 {
		return errors.New("empty users to items adjacency table")
	}

	// Check that there is a weight for each item present in adjacency tables.
	numItems := len(itemWeights)
	var m int
	for _, userItems := range usersToWeightedItems {
		for item, w := range userItems {
			if w < 0 {
				return errors.New("there is a negative weight in usersToWeightedItems")
			}
			if item > m {
				m = item
			}
		}
	}
	if numItems <= m {
		return errors.New("there are more items in UsersToItems than there are weights")
	}

	return nil
}

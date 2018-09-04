// Package birdland provides algorithms to recommend music. It is also the name
// of a legendary Jazz Club.
//
// Bird
//
// Most recommenders are based on collaborative filtering. Traditionally,
// algorithms compute a factorization of the user-item matrix into
// low-dimension, latent vectors for items and users. Recommendations are based
// on a proximity measure between these vectors.
//
// Bird (from Charlie "Bird" Parker) is a recommendation algorithm that is
// based on random walks on the user-item graph. It relies on a (currently
// unweighted) user-item graph which represents interactions of any kind
// between each user and some of the items, and on a list of item weights. To
// initialize Bird with sensible defaults, you can write:
//
//  charlie, err := NewBird(itemWeights, usersToItems)
// 	if err != nil {
// 		log.Errorf("failed to initialize Bird")
// 	}
//
// Bird draws `Draws` items from the input query which are starting points for
// random walks of depth `Depth`. It is possible to specify the values of these
// variables with
//
//  charlie := NewBird(itemWeights, usersToItems, Draws(10000), Depth(5))
// 	if err != nil {
// 		log.Errorf("failed to initialize Bird")
// 	}
//
// Once Bird initialized you can process a query (a slice of `QueryItem`) with
//
// 	visitedItems, referrers, err := charlie.Process(query)
//
// the visitedItems and referrers are then used to produce recommendations.
// Note that Bird does not support empty queries.
//
// It is possible (although not desirable) that an item in the query refers to
// an item with which no one has interacted. We ignore said for the rest of the
// calculations.
package birdland

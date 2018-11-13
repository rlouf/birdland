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
//  charlie, err := NewBird(NewBirdCfg(), itemWeights, usersToItems) if err !=
//  nil { log.Errorf("failed to initialize Bird") }
//
// Bird draws `Draws` items from the input query which are starting points for
// random walks of depth `Depth`. It is possible to specify the values of these
// variables with
//
//  cfg := NewBirdCfg() cfg.Draws = 10000 cfg.Depth = 5
//
//  charlie := NewBird(cfg, itemWeights, usersToItems) if err != nil {
//  log.Errorf("failed to initialize Bird") }
//
// Once Bird initialized you can process a query (a slice of `QueryItem`) with
//
// 	visitedItems, referrers, err := charlie.Process(query)
//
// the visitedItems and referrers are then used to produce recommendations.
// Note that Bird does not support empty queries.
//
// It is possible (although not desirable) that an item in the query refers to
// an item no one has interacted with. We ignore said item for the rest of the
// calculations.
//
// Use cases are recommendations based on a item/container bipartite graph. For
// instance: - Recommend new artists/songs based on user-item relationships; -
// Recommend users based on the same data; - Recommend new songs for a
// playlist/radio; - Recommend playlists/radio
//
//
// Emu
//
// The emu is the heaviest bird there is out there. Although Bird is very
// efficient, it considers that all user-item interactions are created equal.
// However, users are likely to listen to certain artists more. Liking songs
// from an artist once is more meaningful than listening to one of their songs
// once. It is important to take these differences into account in the
// recommender.
//
// The only difference between Emu and Bird is that Emu relies on a weighted
// user-item graph. Typically, the weights will be a combination of the number
// of plays, of shares, and likes. We still need to input a global list of
// weights for the items. To initialize Emu with sensible defaults you can
// write:
//
//  mingus, err := NewEmu(NewBirdCfg(), itemWeights, usersToItems) if err !=
//  nil { log.Errorf("failed to initialize Bird") }
//
// where `usersToItems` is now of type `[]map[int]float64`. `mingus` is now of
// type `*Bird`. The only technical difference between Bird and Emu in the code
// is the way we initialize samplers for each users; Bird and Emu share
// everything else.
package birdland

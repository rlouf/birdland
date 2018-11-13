# Birdland

A famous Jazz club. Also a recommendation library. The library is composed of different elements.

**samplers**
- `tower_sampler.go` implements the tower sampling algorithm to sample from a discrete distribution;
- `alias_sampler.go` implements the alias sampling algorithm to sample from a discrete distribution;

**processors**
- `bird.go` implements a simple recommender engine based on a user-item graph;
- `emu.go` is a recommender engine based on a user-item weighted graph;

**recommenders**
- `recommend.go` contains the functions used to produce recommendations from the engines.


## Recommenders

### Bird

You first need to initialize the recommender with a list of item weights, and two adjacency tables. Assuming
that we are interested in recommending artists or users based on listening data:

```golang
package main
import "github.com/vimies/birdland"

artistWeights := make([]float64, numArtists} // For instance the inverse popularity of artist.
usersToArtists := make([][]int) // For each user, the artists they listened to (liked, followed, etc.)
cfg := NewBirdCfg() // Default of 1000 draws and depth 1

bird, err := birdland.NewBird(cfg, artistWeights, usersToArtists)
```

which initializes the recommender. The recommender processes a query---a list of (artist, weight) pairs---and
outputs a list of artists and their corresponding referrers.

```golang

query := []QueryItem{} // QueryItem{Item int, Weight float64}
items, referrers, err := bird.Process(query)
```

The query, typically, will be a list of the artists a user has listened to along with the number of times they
have listened to it. It can also be, during the onboarding, a list of artists with equal weights. 

Using `items` and `referrers` we can recommend either artists or referrers.

### Emu

The way Emu works is very similar to Bird. The only difference lies in the initialization; instead of taking a
simple bipartite graph `[][]int` as an input, Emu takes a weighted bipartite graph `[]map[int]float64`.
Assuming we want to recommend users and artists based on plays, likes and shares related to artists:

```golang
package main
import "github.com/vimies/birdland"

artistWeights := make([]float64, numArtists} // For instance the inverse popularity of artist.
usersToWeightedArtists := make([]map[int]float64) // For each user, the artists they listened to (liked, followed, etc.)
cfg := NewBirdCfg() // Default of 1000 draws and depth 1

emu, err := birdland.NewEmu(cfg, artistWeights, usersToWeightedArtists)
```

Processing queries is done in the exact same way.

## Making recommendations

### Items

If we want to recommend artists based on the query

```golang
recommendedArtists := birdland.RecommendItems(items, referrers)
```

which produces an ordered `[]int` that contains the id of the recommended artists. 

### User

If we want to recommend users based on the query

```golang
recommendedUsers := birdland.RecommendUsers(items, referrers)
```

which produces an ordered `[]int` that contains the id of the recommended users. 

## Notes

The algorithm is in fact very general; we can think of the following applications:

- If the items are users and collections songs, it is possible to recommend new users based on a list of
  users.
- If the collection is a playlist, it is possible to recommend songs/artists based on their co-occurence in
  playlists. If the Item is a playlist and the collection a song, we can recommend similary playlists.

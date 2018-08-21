# Birdland

A famous Jazz club. Also a recommendation library. The library is composed of the following files:

- `tower_sampler.go` implements the tower sampling algorithm to sample from a discrete distribution;
- `alias_sampler.go` implements the alias sampling algorithm to sample from a discrete distribution;
- `bird.go` contains the simple recommender engine;
- `flock.go` contains the social recommender engine;
- `recommend.go` contains the functions used to produce recommendations from the engines.


## Recommendation

You first need to initialize the recommender with a list of item weights, and two adjacency tables. Assuming
that we are interested in recommending artists or users based on listening data:

```golang
package main
import "github.com/vimies/birdland"

artistWeights := make([]float64, numArtists} // For instance the inverse popularity of artist.
usersToArtists := make([][]int) // For each user, the artists they listened to (liked, followed, etc.)
artistsToUsers := make([][]int) // For each artist, the users who listened to (liked, followed, etc.) them.

bird, err := birdland.NewBird(artistWeights, usersToArtists, artistsToUsers)
```

which initializes the recommender. The recommender processes a query---a list of (artist, weight) pairs---and
outputs a list of artists and their corresponding referrers.

```golang

query := []QueryItem{} // QueryItem{Item int, Weight float64}
items, referrers := bird.Process(query)
```

The query, typically, will be a list of the artists a user has listened to along with the number of times they
have listened to it. It can also be, during the onboarding, a list of artists with equal weights. 

Using `items` and `referrers` we can recommend either artists or referrers.

### Recommending artists

If we want to recommend artists based on the query

```golang
recommendedArtists := bird.RecommendItems(items, referrers)
```

which produces an ordered `[]int` that contains the id of the recommended artists. 

### Recommending users

If we want to recommend users based on the query

```golang
recommendedUsers := bird.RecommendUsers(items, referrers)
```

which produces an ordered `[]int` that contains the id of the recommended users. 

## Notes

The algorithm is in fact very general; we can think of the following applications:

- If the items are users and collections songs, it is possible to recommend new users based on a list of
  users.
- If the collection is a playlist, it is possible to recommend songs/artists based on their co-occurence in
  playlists. If the Item is a playlist and the collection a song, we can recommend similary playlists.

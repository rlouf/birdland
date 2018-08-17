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

bird := birdland.NewBird(artistWeights, usersToArtists, artistsToUsers)
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

It would be possible to recommend music and other users based on a list of users along with weights by
permutation of the word "user" and "artist" in the examples above.Think about how to wrap the code up so that
it doesn't end being confusing.

We can recommend music based on playlists by permutation of the word "user" by "playlist". That way, it is
also be possible to suggest tracks/artists to complete a playlist.




# Social recommendation (WIP)

You first need to initialize the recommender with a list of item weights, and two adjacency tables for the
item-user bipartite graph, and one adjacency table for the user-user (directed or undirected) graph. Assuming
that we are interested in recommending artists or users based on listening data:

```golang
package main
import "github.com/vimies/birdland"

artistWeights := make([]float64, numArtists} // For instance the inverse popularity of artist.
usersToArtists := make([][]int) // For each user, the artists they listened to (liked, followed, etc.)
artistsToUsers := make([][]int) // For each artist, the users who listened to (liked, followed, etc.) them.
userAdjacency := make([][]int) // For each user, the use they follow

flock := birdland.NewFlock(artistWeights, usersToArtists, artistsToUsers, userAdjacency)
```

which initializes the recommender. The recommender processes a query---a list of (artist, weight) pairs along
with a user id---and, again, outputs a list of artists and referrers.
outputs a list of artists and their corresponding referrers.

```golang

query := []Query{} // Query{User int, Query: []QueryItem}
items, referrers := flock.Process(query)
```

The query, typically, will be a list of the artists a user has listened to along with the number of times they
have listened to it. It can also be, during the onboarding, a list of artists with equal weights. 

Using `items` and `referrers` we can recommend either artists or referrers.

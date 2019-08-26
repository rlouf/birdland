<p align="center">
  <img src="https://raw.githubusercontent.com/rlouf/birdland/master/media/birdland.png?token=AA5UP5EFQWUPLZYDB3E2WYK46JAL2">
</p>

#

Birdland is a famous Jazz club. It is also a recommendation library.

Recommending with Birdland is a two-step process: exploring the possibilities
then extracting recommendations from these possibilities. To explore the
possibilities, the algorithm performs a random walk on the (biaised) user-items
bipartite graph starting from a list of items provided as an input; Birdland is
a collaborative filtering method. This random walk generates a list of (user,
item) pairs that are processed by the recommender which returns a list of
recommended items.

Birdland has some advantages over other collaborative filtering algorithms:

- *It requires no pretraining.*  
  Most collaborative algorithm come with hidden costs. Not only do
  you need to maintain an extra service and a database, you also need to
  solve an extra problem: the N-nearest-neighbour search. 
- *It is fast.*  
  We achieved sub-millisecond performance on an API serving recommendation
  of millions of items for a million users.
- *It is simple to reason about, thus to customize*  
  To build `Bird`, the original algorithm, we started from the simple question:
  how would I look for new music to listen? Back in the LastFM days I would
  look for users who had listened to similar artists, what they've listened to
  etc. and trust more users who had very similar tastes. `Bird` does exactly
  that, but a million times faster than you would.
  There is something you do not like about this story? Well, you can adapt
  `Bird`, or use `Emu`.
- *It generalizes to a social recommender.*  
  `Weaver` uses the social network between users to inform recommendations.
- *It solves the long-tail problem for a specific set of parameters.*  
  Blog post to come. Now, whether this is desirable or not is another debate.
- *It is ready for production as is.*  
  Birdland has been tested succesfully in production. Import `birdland` in the
  service that will implement the API that serves the recommendations, plug
  in the data and you're all set.
 
 The codebase is organized around the following components:
  
**samplers**
- `tower_sampler.go` implements the tower sampling algorithm to sample from a discrete distribution;
- `alias_sampler.go` implements the alias sampling algorithm to sample from a
  discrete distribution.

**explorers**
- `bird.go` implements a simple recommender engine based on a user-item graph;
- `emu.go` is a recommender engine based on a user-item weighted graph;
- `weaver.go` is a recommender engine based on the user-item bipartite graph and
  the user-user social graph.
  
**recommenders**
- `recommend.go` contains the functions used to produce recommendations from the engines.


## Engines

### Bird

Named after Charlie "Bird" Parker.

You first need to initialize the recommender with a list of item weights, and two adjacency tables. Assuming
that we are interested in recommending artists or users based on listening data:

```golang
package main
import "github.com/rlouf/birdland"

artistWeights := make([]float64, numArtists} // For instance the inverse popularity of artist.
usersToArtists := make([][]int) // For each user, the artists they listened to (liked, followed, etc.)
cfg := NewBirdCfg() // Default of 1000 draws and depth 1

bird, err := birdland.NewBird(cfg, artistWeights, usersToArtists)
```

which initializes the engine. The engine processes queries---lists of (artist, weight) pairs---and
outputs a list of artists and their referrers.

```golang
query := []QueryItem{} // QueryItem{Item int, Weight float64}
items, referrers, err := bird.Process(query)
```

We can then user `items` and `referrers` to recommend either artists or
referrers (see the "Recommenders" section below). All algorithms use two
parameters:

- the depth of the random walk;
- the number of draws of random walks that are performed.

they can be tuned by initializing the configuration passed to `NewBird` by hand:

```
cfg = BirdCfg{Depth: 2, Draws: 10000}
```

### Emu

The emu is a heavy bird.

The way Emu works is very similar to Bird. The only difference lies in the
initialization; instead of taking a simple bipartite graph `[][]int` as an
input, Emu takes a weighted bipartite graph `[]map[int]float64`. The weights can
be anything from the number of plays, to the number of likes or a score given by
the user.

```golang
package main
import "github.com/rlouf/birdland"

artistWeights := make([]float64, numArtists} // For instance the inverse popularity of artist.
usersToWeightedArtists := make([]map[int]float64) // For each user, the artists they listened to (liked, followed, etc.)
cfg := NewBirdCfg() // Default of 1000 draws and depth 1

emu, err := birdland.NewEmu(cfg, artistWeights, usersToWeightedArtists)
```

Processing queries is done in the exact same way.

### Weaver (cleaning)

Weaver is allegedly the most social bird.

The same way Emu attributes different weighs to each item, Weaver attributes
different weights to each user. Indeed, in a social network you do not weigh
recommendations by strangers and by acquaintances the same way.

```golang
package main
import "github.com/rlouf/birdland"

cfg := NewBirdCfg{}
weaver, err := birdland.NewWeaver(cfg, itemWeights, usersToItems, socialWeights) 
```

## Recommenders

Since the engines traverse both users and items, we can recommend one or the 
other (or both) indifferently *within the same query*. Birdland provides
two functions to produce recommendations from the engines' outputs.

These two functions were used to provide a stable interface for the services
that use Birdland and so strategies could be swapped without affecting said
services. You can consult `recommend.go` to see the available strategies.

```golang
recommendedArtists := birdland.RecommendItems(items, referrers)
```

Produces an ordered `[]int` that contains the id of the recommended artists. 

```golang
recommendedUsers := birdland.RecommendUsers(items, referrers)
```

Produces an ordered `[]int` that contains the id of the recommended users. 


## Contribute

Questions, Issues or PRs are very welcome! Please read the `CONTRIBUTING.md` file
first, and happy forking.

## Credits

The icon was made by <a href="https://www.freepik.com/?__hstc=57440181.3c24109fd911bedc6428debe60ee2cde.1558556981649.1558556981649.1558556981649.1&__hssc=57440181.6.1558556981649&__hsfp=4016125896" title="Freepik">Freepik</a> from <a href="https://www.flaticon.com/" 			    title="Flaticon">www.flaticon.com</a> is licensed by <a href="http://creativecommons.org/licenses/by/3.0/" 			    title="Creative Commons BY 3.0" target="_blank">CC 3.0 BY

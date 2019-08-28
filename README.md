<p align="center">
  <img src="https://raw.githubusercontent.com/rlouf/birdland/master/media/birdland.png?token=AA5UP5EFQWUPLZYDB3E2WYK46JAL2">
</p>

#

[![CircleCI](https://circleci.com/gh/rlouf/birdland.svg?style=svg)](https://circleci.com/gh/rlouf/birdland)

Birdland is a famous Jazz club. It is also a recommendation library.

Birdland is a collaborative filtering algorithm in two steps: exploration and
recommendation. To explore possibilities the algorithm performs a random walk on
the (biaised) user-items bipartite graph starting from a list of items provided
as an input. This random walk generates a list of (user, item) pairs that are
processed by the recommender which returns a list of recommended items.

Birdland has some advantages over other collaborative filtering algorithms:

- *It requires no pretraining.*  
  Most collaborative algorithm come with hidden costs. Not only do
  you need to maintain an extra service and database, you also need to
  solve an additional problem: find the [N nearest neighbors](https://en.wikipedia.org/wiki/Nearest_neighbor_search) 
  of a vector (spoiler: it is not an easy problem). 
- *It is fast.*  
  We achieved performance of the order of the millisecond on an API serving recommendation
  of millions of items for a million users.
- *It is simple to reason about, thus to customize.*  
  To build `Bird` we started from the simple question: how would I look for new
  music to listen? Back in the LastFM days I would look for users who had
  listened to similar artists, what they've listened to etc. and trust more
  users who had very similar tastes. `Bird` does exactly that, but a million
  times faster than you would.  There is something you do not like about this
  story? Well, you can adapt `Bird`, or use `Emu`.
- *It generalizes to a social recommender.*  
  `Weaver` uses the social network
  between users to inform recommendations.
- *It recommends both items and users in one pass.*
  No need to find the N nearest neighbors again.
- *It solves the long-tail problem for a specific set of parameters.*  
  (Blog post to come) Now, whether this is desirable or not is another debate.
- *It is ready for production.*  
  Birdland has been tested succesfully in production. Import `birdland` in the
  service that implements the recommendation API, plug in the data and
  you're all set.
 
The codebase is organized around the following components:
  
**samplers**
- `tower_sampler.go` implements the tower sampling algorithm to sample from a
  discrete distribution;
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

Named after [Charlie "Bird" Parker](https://www.youtube.com/watch?v=LphuCadyQi0).

The very first step is to map the users and items to sets of consecutive
integers (starting with 0). This avoids working with maps, which substantially
improves performance.

Initialize the engine with a list of item weights, and the (user, item)
adjacency table: 

```golang
package main
import "github.com/rlouf/birdland"

artistWeights := make([]float64, numArtists) // global weight attributed to each artist
usersToArtists := make([][]int, numUsers) // for each user the list of artists they listened to (liked, followed, etc.)
cfg := NewBirdCfg()

bird, err := birdland.NewBird(cfg, artistWeights, usersToArtists)
```

This needs to be done only once (provided your data do not change). The engine
processes queries---lists of (artist_id, weight) pairs---and outputs a list of
artists and their referrers:

```golang
query := []QueryItem{} // QueryItem{Item int, Weight float64}
items, referrers, err := bird.Process(query)
```

We can then use `items` and `referrers` to recommend either artists or
referrers (see the "Recommenders" section below). All engines depend 
on two parameters:

- the depth of the random walk;
- the number of random walks that are performed (number of samples drawn from the
  query).

they can be tuned by initializing the configuration passed to `NewBird` by hand:

```
cfg = BirdCfg{Depth: 2, Draws: 10000}
```

### Emu

The emu is a heavy bird ([the 5th heaviest](https://en.wikipedia.org/wiki/List_of_largest_birds#Table_of_heaviest_living_bird_species)).

Emu works very similarly to Bird. The only difference lies in the
initialization; instead of taking a simple bipartite graph `[][]int` as an
input, Emu takes a weighted bipartite graph `[]map[int]float64`. In the context
of music recommendation, the weight can for instance be the number of times
the user played tracks from an artist.

```golang
package main
import "github.com/rlouf/birdland"

artistWeights := make([]float64, numArtists)
usersToWeightedArtists := make([]map[int]float64, numUsers)
cfg := NewBirdCfg() // Default of 1000 draws and depth 1

emu, err := birdland.NewEmu(cfg, artistWeights, usersToWeightedArtists)
```

Everything else is exactly the same.

### Weaver (cleaning)

Weavers are allegedly [very sociable birds](https://en.wikipedia.org/wiki/Sociable_weaver).

The same way Emu attributes different weighs to each item, Weaver attributes
different weights to each user. This follows from the idea that you would not
weigh recommendations by strangers and by acquaintances the same way.

```golang
package main
import "github.com/rlouf/birdland"

cfg := NewWeaverCfg()
itemWeights := make([]float64, numItems)
usersToItems := make([][]int, numUsers)

weaver, err := birdland.NewWeaver(cfg, itemWeights, usersToItems, socialGraph) 
```

We give to users who are not connected to the current user a default weight of 1.
This default behavior can be changed by initializing the configuration by hand:

```golang
cfg := WeaverCfg{DefaultWeight: 0, BirdCfg: NewBirdCfg()}
```

which would only consider recommendations coming from direct connections.

## Recommenders

Since the engines traverse both users and items, we can recommend one or the 
other (or both) indifferently *within the same query*. Birdland provides
several functions to produce recommendations from the engines' outputs.

Two functions were defined to provide a stable interface for the services
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
first, then happy forking.

## Credits

The icon was made by <a href="https://www.freepik.com/?__hstc=57440181.3c24109fd911bedc6428debe60ee2cde.1558556981649.1558556981649.1558556981649.1&__hssc=57440181.6.1558556981649&__hsfp=4016125896" title="Freepik">Freepik</a> from <a href="https://www.flaticon.com/" 			    title="Flaticon">www.flaticon.com</a> is licensed by <a href="http://creativecommons.org/licenses/by/3.0/" 			    title="Creative Commons BY 3.0" target="_blank">CC 3.0 BY

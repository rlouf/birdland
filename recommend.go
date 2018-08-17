package birdland

import (
	"sort"
)

type Pair struct {
	Object     int
	Occurences int // number of occurences
}

type PairList []Pair // necessary evil to sort map by value

func (p PairList) Len() int           { return len(p) }
func (p PairList) Less(i, j int) bool { return p[i].Occurences < p[j].Occurences }
func (p PairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

// RecommendItems and RecommendUsers contain the current methods that should be
// used to recommend items and users in production.  Used as an interface so
// backend developpers do not need to worry about the zoology of recommending
// methods.
func RecommendItems(items, referrers []int) []int { return RecommendMostVisited(items) }
func RecommendUsers(items, referrers []int) []int { return RecommendMostVisited(referrers) }

// RecommendMostVisited recommends the items in descending order of the number
// of visits by the processing algorithm. It is a very naive approach and
// probably should not be used in production. Works indifferently to
// recommend users or items.
func RecommendMostVisited(items []int) []int {
	countItems := make(map[int]int)
	for _, item := range items {
		countItems[item] += 1
	}

	pairList := PairList{}
	for item, count := range countItems {
		pairList = append(pairList, Pair{item, count})
	}

	sort.Sort(sort.Reverse(pairList))
	recommendedItems := make([]int, len(countItems))
	for i, pair := range pairList {
		recommendedItems[i] = pair.Object
	}

	return recommendedItems
}

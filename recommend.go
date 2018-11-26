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

// RecommendConsensus recommends the item by descending order of the number of
// unique referrers. With the data currently available, it is only possible to
// recommend items this way.
func RecommendConsensus(items, referrers []int) []int {

	if len(items) != len(referrers) {
		panic("items and referrers do not have the same number of elements")
	}

	mapUniqueReferrers := make(map[int]map[int]bool)
	for i, item := range items {
		if _, ok := mapUniqueReferrers[item]; !ok {
			mapUniqueReferrers[item] = map[int]bool{}
		}
		mapUniqueReferrers[item][referrers[i]] = true
	}

	countUniqueReferrers := PairList{}
	for item, referrersMap := range mapUniqueReferrers {
		countUniqueReferrers = append(countUniqueReferrers, Pair{item, len(referrersMap)})
	}

	sort.Sort(sort.Reverse(countUniqueReferrers))
	recommendedItems := make([]int, 0, len(items))
	for _, pair := range countUniqueReferrers {
		recommendedItems = append(recommendedItems, pair.Object)
	}

	return recommendedItems
}

// RecommendTrust recommends items based on how much we can trust their referrers.
// The algorithm begins with attributing a weight to the refferers proportional to
// the number of times it traversed them. The tracks are recommended by descending
// order of the cumulated weights.
func RecommendTrust(items, referrers []int) []int {

	if len(items) != len(referrers) {
		panic("items and referrers do not have the same number of elements")
	}

	countReferrerTraversals := make(map[int]int)
	for _, referrer := range referrers {
		countReferrerTraversals[referrer] += 1
	}

	itemWeights := make(map[int]int)
	for i, item := range items {
		itemWeights[item] += countReferrerTraversals[referrers[i]]
	}

	pairList := PairList{}
	for item, count := range itemWeights {
		pairList = append(pairList, Pair{item, count})
	}

	sort.Sort(sort.Reverse(pairList))
	recommendedItems := make([]int, len(itemWeights))
	for i, pair := range pairList {
		recommendedItems[i] = pair.Object
	}

	return recommendedItems
}

package birdland

import "testing"

type MostVisitedCase struct {
	Name     string
	Input    []int
	Expected []int
}

var mostVisited_table = []MostVisitedCase{
	{
		Name:     "Empty input",
		Input:    []int{},
		Expected: []int{},
	},
	{
		Name:     "Typical input",
		Input:    []int{1, 2, 2, 2, 3, 3, 0, 0, 0, 0, 0, 5, 5, 5, 5, 5, 5},
		Expected: []int{5, 0, 2, 3, 1},
	},
}

type ConsensusCase struct {
	Name      string
	Items     []int
	Referrers []int
	Expected  []int
}

var consensus_table = []ConsensusCase{
	{
		Name:      "Empty input",
		Items:     []int{},
		Referrers: []int{},
		Expected:  []int{},
	},
	{
		Name:      "Typical input",
		Items:     []int{1, 2, 2, 2, 1, 1, 1, 3, 3},
		Referrers: []int{1, 3, 4, 5, 1, 1, 1, 2, 1},
		Expected:  []int{2, 3, 1},
	},
}

type TrustCase struct {
	Name      string
	Items     []int
	Referrers []int
	Expected  []int
}

var trust_table = []TrustCase{
	{
		Name:      "Typical input",
		Items:     []int{1, 1, 1, 2, 5, 5, 5, 4},
		Referrers: []int{1, 1, 1, 1, 2, 3, 4, 5},
		Expected:  []int{1, 2, 5, 4},
	},
}

func TestRecommendMostVisited(t *testing.T) {
	for _, ex := range mostVisited_table {
		recommended := RecommendMostVisited(ex.Input)
		if len(recommended) != len(ex.Expected) {
			t.Errorf("RecommendMostVisited: %s: discrepancy in the length of the recommendations: expected %d, got %d", ex.Name, len(ex.Expected), len(recommended))
		}
		for i, r := range recommended {
			if r != ex.Expected[i] {
				t.Errorf("RecommendMostVisited: %s: expected %d, got %d", ex.Name, ex.Expected, recommended)
				break
			}
		}
	}
}

func TestRecommendConsensus(t *testing.T) {
	for _, ex := range consensus_table {
		recommended := RecommendConsensus(ex.Items, ex.Referrers)
		if len(recommended) != len(ex.Expected) {
			t.Errorf("RecommendConsensus: %s: discrepancy in the length of the recommendations: expected %d, got %d", ex.Name, len(ex.Expected), len(recommended))
		}
		for i, r := range recommended {
			if r != ex.Expected[i] {
				t.Errorf("RecommendConsensus: %s: expected %d, got %d", ex.Name, ex.Expected, recommended)
				break
			}
		}
	}
}

func TestRecommendTrust(t *testing.T) {
	for _, ex := range trust_table {
		recommended := RecommendTrust(ex.Items, ex.Referrers)
		if len(recommended) != len(ex.Expected) {
			t.Errorf("RecommendTrust: %s: discrepancy in the length of the recommendations: expected %d, got %d", ex.Name, len(ex.Expected), len(recommended))
		}
		for i, r := range recommended {
			if r != ex.Expected[i] {
				t.Errorf("RecommendTrust: %s: expected %d, got %d", ex.Name, ex.Expected, recommended)
				break
			}
		}
	}
}

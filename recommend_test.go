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

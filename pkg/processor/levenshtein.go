package processor

import (
	levenshtein2 "github.com/psykhi/uno/pkg/levenshtein"
	"math"
)

type levenshtein struct {
	seen         [][]string
	maxDiffRatio float64
}

func newLevenshtein(maxDiffRatio float64) *levenshtein {
	seen := make([][]string, 0)
	return &levenshtein{seen: seen, maxDiffRatio: maxDiffRatio}
}

func (le *levenshtein) process(in Line) Line {
	in.IsNew = true
	for i, l := range le.seen {
		maxDiff := int(math.Ceil(float64(len(in.Tokens)) * le.maxDiffRatio))
		lenDiff := len(in.Tokens) - len(l)
		if lenDiff < 0 {
			lenDiff = -lenDiff
		}
		if lenDiff > maxDiff {
			continue
		}
		d := levenshtein2.LevenshteinDistanceK(in.Tokens, l, nil, nil, maxDiff)
		if d <= maxDiff && d >= 0 {
			if i > 0 {
				le.seen[i], le.seen[i-1] = le.seen[i-1], le.seen[i]
			}
			in.IsNew = false
			return in
		}
	}
	le.seen = append(le.seen, in.Tokens)
	return in
}

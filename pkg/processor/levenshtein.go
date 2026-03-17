package processor

import (
	levenshtein2 "github.com/psykhi/uno/pkg/levenshtein"
	"math"
)

type levenshtein struct {
	seen         [][]string
	maxDiffRatio float64
	v0           []int
	v1           []int
}

func newLevenshtein(maxDiffRatio float64) *levenshtein {
	seen := make([][]string, 0)
	return &levenshtein{seen: seen, maxDiffRatio: maxDiffRatio}
}

func (le *levenshtein) process(in Line) Line {
	in.IsNew = true
	maxDiff := int(math.Ceil(float64(len(in.Tokens)) * le.maxDiffRatio))
	for i, l := range le.seen {
		lenDiff := len(in.Tokens) - len(l)
		if lenDiff < 0 {
			lenDiff = -lenDiff
		}
		if lenDiff > maxDiff {
			continue
		}
		maxLen := len(in.Tokens)
		if len(l) > maxLen {
			maxLen = len(l)
		}
		if cap(le.v0) < maxLen+1 {
			le.v0 = make([]int, maxLen+1)
			le.v1 = make([]int, maxLen+1)
		}
		d := levenshtein2.LevenshteinDistanceK(in.Tokens, l, le.v0[:maxLen+1], le.v1[:maxLen+1], maxDiff)
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

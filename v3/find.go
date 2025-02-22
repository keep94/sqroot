package sqroot

import (
	"iter"
	"slices"

	"github.com/keep94/itertools"
)

// Matches returns all the 0 based positions in s where pattern is found.
func Matches(s Sequence, pattern []int) iter.Seq[int] {
	return matches(s, slices.Clone(pattern))
}

// BackwardMatches returns all the 0 based positions in s where pattern is
// found from last to first.
func BackwardMatches(s FiniteSequence, pattern []int) iter.Seq[int] {
	if len(pattern) == 0 {
		return zeroPattern(s.Backward())
	}
	return kmp(s.Backward(), patternReverse(pattern), true)
}

// Find returns a function that returns the next zero based index of the
// match for pattern in s. If s has a finite number of digits and there
// are no more matches for pattern, the returned function returns -1.
// Pattern is a sequence of digits between 0 and 9.
//
// Deprecated: Use Matches instead.
func Find(s Sequence, pattern []int) func() int {
	if len(pattern) == 0 {
		return zeroPatternOld(s.Iterator())
	}
	return kmpOld(s.Iterator(), slices.Clone(pattern), false)
}

// FindFirst finds the zero based index of the first match of pattern in s.
// FindFirst returns -1 if pattern is not found only if s has a finite number
// of digits. If s has an infinite number of digits and pattern is not found,
// FindFirst will run forever. pattern is a sequence of digits between 0 and 9.
func FindFirst(s Sequence, pattern []int) int {
	return collectFirst(matches(s, pattern))
}

// Deprecated: Use golang iterators with the github.com/keep94/itertools
// library. This is equivalent to
// slices.Collect(itertools.Take(sqroot.Matches(s, pattern), n))
func FindFirstN(s Sequence, pattern []int, n int) []int {
	return slices.Collect(itertools.Take(Matches(s, pattern), n))
}

// Deprecated: This is equivalent to
// slices.Collect(sqroot.Matches(s, pattern))
func FindAll(s FiniteSequence, pattern []int) []int {
	return slices.Collect(Matches(s, pattern))
}

// FindLast finds the zero based index of the last match of pattern in s.
// FindLast returns -1 if pattern is not found in s. pattern is a sequence of
// digits between 0 and 9.
func FindLast(s FiniteSequence, pattern []int) int {
	return collectFirst(BackwardMatches(s, pattern))
}

// Deprecated: Use golang iterators with the github.com/keep94/itertools
// library. This is equivalent to
// slices.Collect(itertools.Take(sqroot.BackwardMatches(s, pattern), n))
func FindLastN(s FiniteSequence, pattern []int, n int) []int {
	return slices.Collect(itertools.Take(BackwardMatches(s, pattern), n))
}

// FindR returns a function that starts at the end of s and returns the
// previous zero based index of the match for pattern in s with each call.
// If there are no more matches for pattern, the returned function returns
// -1. pattern is a sequence of digits between 0 and 9.
//
// Deprecated: Use BackwardMatches instead.
func FindR(s FiniteSequence, pattern []int) func() int {
	if len(pattern) == 0 {
		return zeroPatternOld(s.Reverse())
	}
	return kmpOld(s.Reverse(), patternReverse(pattern), true)
}

func matches(s Sequence, pattern []int) iter.Seq[int] {
	if len(pattern) == 0 {
		return zeroPattern(s.All())
	}
	return kmp(s.All(), pattern, false)
}

func collectFirst(seq iter.Seq[int]) int {
	for index := range seq {
		return index
	}
	return -1
}

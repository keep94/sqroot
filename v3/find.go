package sqroot

import (
	"iter"
	"slices"
)

// Matches returns all the 0 based positions in s where pattern is found.
func Matches(s Sequence, pattern []int) iter.Seq[int] {
	return matches(s, slices.Clone(pattern))
}

// BackwardMatches returns all the 0 based positions in s where pattern is
// found from last to first.
func BackwardMatches(s FiniteSequence, pattern []int) iter.Seq[int] {
	patternInReverse := patternReverse(pattern)
	return func(yield func(index int) bool) {
		gen := findR(s, patternInReverse)
		for index := gen(); index != -1; index = gen() {
			if !yield(index) {
				return
			}
		}
	}
}

// Find returns a function that returns the next zero based index of the
// match for pattern in s. If s has a finite number of digits and there
// are no more matches for pattern, the returned function returns -1.
// Pattern is a sequence of digits between 0 and 9.
func Find(s Sequence, pattern []int) func() int {
	return find(s, slices.Clone(pattern))
}

// FindFirst finds the zero based index of the first match of pattern in s.
// FindFirst returns -1 if pattern is not found only if s has a finite number
// of digits. If s has an infinite number of digits and pattern is not found,
// FindFirst will run forever. pattern is a sequence of digits between 0 and 9.
func FindFirst(s Sequence, pattern []int) int {
	return collectFirst(matches(s, pattern))
}

// FindFirstN works like FindFirst but it finds the first n matches and
// returns the zero based index of each match. If s has a finite
// number of digits, FindFirstN may return fewer than n matches.
// Like FindFirst, FindFirstN may run forever if s has an infinite
// number of digits, and there are not n matches available.
// pattern is a sequence of digits between 0 and 9.
func FindFirstN(s Sequence, pattern []int, n int) []int {
	return collectFirstN(matches(s, pattern), n)
}

// FindAll finds all the matches of pattern in s and returns the zero based
// index of each match. pattern is a sequence of digits between 0 and 9.
func FindAll(s FiniteSequence, pattern []int) []int {
	return slices.Collect(matches(s, pattern))
}

// FindLast finds the zero based index of the last match of pattern in s.
// FindLast returns -1 if pattern is not found in s. pattern is a sequence of
// digits between 0 and 9.
func FindLast(s FiniteSequence, pattern []int) int {
	return collectFirst(BackwardMatches(s, pattern))
}

// FindLastN works like FindLast but it finds the last n matches and
// returns the zero based index of each match. Last matches come first
// in returned array. pattern is a sequence of digits between 0 and 9.
func FindLastN(s FiniteSequence, pattern []int, n int) []int {
	return collectFirstN(BackwardMatches(s, pattern), n)
}

// FindR returns a function that starts at the end of s and returns the
// previous zero based index of the match for pattern in s with each call.
// If there are no more matches for pattern, the returned function returns
// -1. pattern is a sequence of digits between 0 and 9.
func FindR(s FiniteSequence, pattern []int) func() int {
	return findR(s, patternReverse(pattern))
}

func findR(s FiniteSequence, patternInReverse []int) func() int {
	if len(patternInReverse) == 0 {
		return zeroPattern(s.Reverse())
	}
	return kmp(s.Reverse(), patternInReverse, true)
}

func find(s Sequence, pattern []int) func() int {
	if len(pattern) == 0 {
		return zeroPattern(s.Iterator())
	}
	return kmp(s.Iterator(), pattern, false)
}

func matches(s Sequence, pattern []int) iter.Seq[int] {
	return func(yield func(index int) bool) {
		gen := find(s, pattern)
		for index := gen(); index != -1; index = gen() {
			if !yield(index) {
				return
			}
		}
	}
}

func collectFirstN(seq iter.Seq[int], n int) []int {
	if n <= 0 {
		return nil
	}
	var result []int
	for index := range seq {
		result = append(result, index)
		if len(result) == n {
			break
		}
	}
	return result
}

func collectFirst(seq iter.Seq[int]) int {
	for index := range seq {
		return index
	}
	return -1
}

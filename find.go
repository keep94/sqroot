package sqroot

import (
	"github.com/keep94/consume2"
)

// Find returns a function that returns the next zero based index of the
// match for pattern in s. If s has a finite number of digits and there
// are no more matches for pattern, the returned function returns -1.
// Pattern is a sequence of digits between 0 and 9.
func Find(s Sequence, pattern []int) func() int {
	if len(pattern) == 0 {
		return zeroPattern(s.FullIterator())
	}
	patternCopy := make([]int, len(pattern))
	copy(patternCopy, pattern)
	return kmp(s.FullIterator(), patternCopy, false)
}

// FindFirst finds the zero based index of the first match of pattern in s.
// FindFirst returns -1 if pattern is not found only if s has a finite number
// of digits. If s has an infinite number of digits and pattern is not found,
// FindFirst will run forever. pattern is a sequence of digits between 0 and 9.
func FindFirst(s Sequence, pattern []int) int {
	iter := find(s, pattern)
	return iter()
}

// FindFirstN works like FindFirst but it finds the first n matches and
// returns the zero based index of each match. If s has a finite
// number of digits, FindFirstN may return fewer than n matches.
// Like FindFirst, FindFirstN may run forever if s has an infinite
// number of digits, and there are not n matches available.
// pattern is a sequence of digits between 0 and 9.
func FindFirstN(s Sequence, pattern []int, n int) []int {
	return asIntSlice(find(s, pattern), consume2.PSlice[int](0, n))
}

// FindAll finds all the matches of pattern in s and returns the zero based
// index of each match. If s has an infinite number of digits, FindAll will
// run forever. pattern is a sequence of digits between 0 and 9.
func FindAll(s Sequence, pattern []int) []int {
	return asIntSlice(find(s, pattern), consume2.Identity[int]())
}

// FindLast finds the zero based index of the last match of pattern in s.
// FindLast returns -1 if pattern is not found in s. If s has an infinite
// number of digits, FindLast will run forever. pattern is a sequence of
// digits between 0 and 9.
func FindLast(s Sequence, pattern []int) int {
	iter := FindR(s, pattern)
	return iter()
}

// FindLastN works like FindLast but it finds the last n matches and
// returns the zero based index of each match. Last matches come first
// in returned array. If s has an infinite number of digits, FindLastN
// will run forever. pattern is a sequence of digits between 0 and 9.
func FindLastN(s Sequence, pattern []int, n int) []int {
	return asIntSlice(FindR(s, pattern), consume2.PSlice[int](0, n))
}

// FindR returns a function that starts at the end of s and returns the
// previous zero based index of the match for pattern in s with each call.
// If there are no more matches for pattern, the returned function returns
// -1. If s has an infinite number of digits, FindR runs forever. pattern is
// a sequence of digits between 0 and 9.
func FindR(s Sequence, pattern []int) func() int {
	if len(pattern) == 0 {
		return zeroPattern(s.FullReverse())
	}
	return kmp(s.FullReverse(), patternReverse(pattern), true)
}

func find(s Sequence, pattern []int) func() int {
	if len(pattern) == 0 {
		return zeroPattern(s.FullIterator())
	}
	return kmp(s.FullIterator(), pattern, false)
}

func asIntSlice(
	gen func() int, pipeline consume2.Pipeline[int, int]) (result []int) {
	consume2.FromIntGenerator(gen, pipeline.AppendTo(&result))
	return
}

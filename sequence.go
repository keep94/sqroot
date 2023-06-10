package sqroot

import (
	"github.com/keep94/consume2"
)

// Sequence represents a sequence of possibly non contiguous digits.
// That is, a Sequence may have holes. For instance, a Sequence could
// be 375XXX695. The 0th digit is a 3; the 1st digit is a 7; the 2nd digit
// is a 5; the 3rd, 4th, and 5th digits are unknown; the 6th digit is a 6;
// the 7th digit is a 9; the 8th digit is a 5.
// Both Digits and Number pointers implement Sequence.
type Sequence interface {
	digitIter() func() (Digit, bool)
}

func fromSequenceWithPositions(
	s Sequence, p Positions, consumer consume2.Consumer[Digit]) {
	r, ok := s.(randomAccess)
	if ok && r.enabled() {

		// Optimization: Just choose what we want instead of iterating over
		// all of s. This can be several orders of magnitude faster if p is
		// small relative to s. However if p is close to the size of s and
		// highly fragmented, this can be an order of magnitude slower.
		// Overall, this optimization is a win because usually p will be
		// small or not highly fragmented.
		for _, pr := range p.ranges {
			consume2.FromGenerator[Digit](
				r.get(pr.Start, pr.End).digitIter(), consumer)
		}
	} else {
		consume2.FromGenerator[Digit](iterateWithPositions(s, p), consumer)
	}
}

func iterateWithPositions(s Sequence, p Positions) func() (Digit, bool) {
	filter := p.filter()
	iter := s.digitIter()
	d, ok := iter()
	limit := p.limit()
	return func() (result Digit, hasResult bool) {
		for ok && d.Position < limit && !hasResult {
			if filter.Includes(d.Position) {
				result = d
				hasResult = true
			}
			d, ok = iter()
		}
		return
	}
}

type randomAccess interface {
	get(start, end int) Sequence
	enabled() bool
}

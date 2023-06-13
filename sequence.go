package sqroot

import (
	"io"
	"os"
	"strings"

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

// Fprint prints digits of s to w. Fprint returns the number of bytes written
// and any error encountered. p contains the positions of the digits to print.
// For options, the default is 50 digits per row, 5 digits per column,
// show digit count, and period (.) for missing digits.
func Fprint(w io.Writer, s Sequence, p Positions, options ...Option) (
	written int, err error) {
	settings := &printerSettings{
		digitsPerRow:    50,
		digitsPerColumn: 5,
		showCount:       true,
		missingDigit:    '.',
	}
	printer := newPrinter(w, p.limit(), mutateSettings(options, settings))
	fromSequenceWithPositions(s, p, printer)
	return printer.byteCount, printer.err
}

// Sprint works like Fprint and prints digits of s to a string.
func Sprint(s Sequence, p Positions, options ...Option) string {
	var builder strings.Builder
	Fprint(&builder, s, p, options...)
	return builder.String()
}

// Print works like Fprint and prints digits of s to stdout.
func Print(s Sequence, p Positions, options ...Option) (
	written int, err error) {
	return Fprint(os.Stdout, s, p, options...)
}

func fromSequenceWithPositions(
	s Sequence, p Positions, consumer consume2.Consumer[Digit]) {
	sr, ok := s.(subRangeSequence)
	if ok && sr.canSubRange() {

		// Optimization: Just choose what we want instead of iterating over
		// all of s. This can be several orders of magnitude faster if p is
		// small relative to s. However if p is close to the size of s and
		// highly fragmented, this can be an order of magnitude slower.
		// Overall, this optimization is a win because usually p will be
		// small or not highly fragmented.
		for _, pr := range p.ranges {
			consume2.FromGenerator[Digit](
				sr.subRange(pr.Start, pr.End).digitIter(), consumer)
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

type subRangeSequence interface {
	subRange(start, end int) Sequence
	canSubRange() bool
}

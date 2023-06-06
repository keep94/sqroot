package sqroot

import (
	"io"

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
	reverseDigitIter() func() (Digit, bool)
	canReverse() bool
}

type part interface {
	digitIter() func() (Digit, bool)
	limit() int
}

type lazyPart struct {
	sequence  Sequence
	positions Positions
}

func newPart(sequence Sequence, positions Positions) part {
	return &lazyPart{sequence: sequence, positions: positions}
}

func (p *lazyPart) limit() int {
	return p.positions.limit()
}

func (p *lazyPart) digitIter() func() (Digit, bool) {
	filter := p.positions.filter()
	iter := p.sequence.digitIter()
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

func fprint(
	w io.Writer, part part, settings *printerSettings) (n int, err error) {
	p := newPrinter(w, part.limit(), settings)
	consume2.FromGenerator[Digit](part.digitIter(), p)
	return p.byteCount, p.err
}

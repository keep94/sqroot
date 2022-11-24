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
// Both Digits and Mantissa implement Sequence.
type Sequence interface {
	positDigitIter() func() positDigit
}

type part interface {
	positDigitIter() func() positDigit
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

func (p *lazyPart) positDigitIter() func() positDigit {
	filter := p.positions.filter()
	iter := p.sequence.positDigitIter()
	pd := iter()
	limit := p.limit()
	return func() positDigit {
		result := invalidPositDigit
		for pd.Valid() && pd.Posit < limit && !result.Valid() {
			if filter.Includes(pd.Posit) {
				result = pd
			}
			pd = iter()
		}
		return result
	}
}

func fprint(
	w io.Writer, part part, settings *printerSettings) (n int, err error) {
	p := newPrinter(w, part.limit(), settings)
	sendPositDigits(part, p)
	return p.byteCount, p.err
}

func sendPositDigits(s Sequence, consumer consume2.Consumer[positDigit]) {
	iter := s.positDigitIter()
	for consumer.CanConsume() {
		pd := iter()
		if !pd.Valid() {
			return
		}
		consumer.Consume(pd)
	}
}

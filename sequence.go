package sqroot

import (
	"io"
	"os"
	"strings"

	"github.com/keep94/consume2"
)

// Sequence represents a sequence of digits. Number pointers implement
// Sequence.
type Sequence interface {
	digitIter() func() (digit, bool)
	reverseDigitIter() func() (digit, bool)
	subRange(start, end int) Sequence
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
	s Sequence, p Positions, consumer consume2.Consumer[digit]) {
	for _, pr := range p.ranges {
		consume2.FromGenerator[digit](
			s.subRange(pr.Start, pr.End).digitIter(), consumer)
	}
}

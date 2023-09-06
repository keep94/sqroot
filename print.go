package sqroot

import (
	"io"
	"os"
	"strings"

	"github.com/keep94/consume2"
)

// Option represents an option for the Print, Fprint, and Sprint methods
type Option interface {
	mutate(p *printerSettings)
}

// DigitsPerRow sets the number of digits per row. Zero or negative means no
// separate rows.
func DigitsPerRow(count int) Option {
	return optionFunc(func(p *printerSettings) {
		p.digitsPerRow = count
	})
}

// DigitsPerColumn sets the number of digits per column. Zero or negative
// means no separate columns.
func DigitsPerColumn(count int) Option {
	return optionFunc(func(p *printerSettings) {
		p.digitsPerColumn = count
	})
}

// ShowCount shows the digit count in the left margin if on is true.
func ShowCount(on bool) Option {
	return optionFunc(func(p *printerSettings) {
		p.showCount = on
	})
}

// MissingDigit sets the character to represent a missing digit.
func MissingDigit(missingDigit rune) Option {
	return optionFunc(func(p *printerSettings) {
		p.missingDigit = missingDigit
	})
}

func bufferSize(size int) Option {
	return optionFunc(func(p *printerSettings) {
		p.bufferSize = size
	})
}

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
	printer := newPrinter(w, p.End(), mutateSettings(options, settings))
	fromSequenceWithPositions(s, p, printer)
	printer.Finish()
	return printer.BytesWritten(), printer.Err()
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
	iter := p.Ranges()
	for pr, ok := iter(); ok; pr, ok = iter() {
		consume2.FromGenerator(
			s.subRange(pr.Start, pr.End).digitIter(), consumer)
	}
}

type optionFunc func(p *printerSettings)

func (o optionFunc) mutate(p *printerSettings) {
	o(p)
}

func mutateSettings(
	options []Option, settings *printerSettings) *printerSettings {
	for _, option := range options {
		option.mutate(settings)
	}
	return settings
}

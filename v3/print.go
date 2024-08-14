package sqroot

import (
	"io"
	"iter"
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

// TrailingLF adds a trailing line feed to what is printed if on is true.
func TrailingLF(on bool) Option {
	return optionFunc(func(p *printerSettings) {
		p.trailingLineFeed = on
	})
}

// LeadingDecimal prints "0." before the first digit if on is true.
func LeadingDecimal(on bool) Option {
	return optionFunc(func(p *printerSettings) {
		p.leadingDecimal = on
	})
}

func bufferSize(size int) Option {
	return optionFunc(func(p *printerSettings) {
		p.bufferSize = size
	})
}

// Sequence represents a sequence of digits of either finite or infinite
// length within the mantissa of a real number. Although they can start
// and optionally end anywhere within a mantissa, Sequences must be
// contiguous. That is they can have no gaps in the middle.
type Sequence interface {

	// All returns the 0 based position of each digit in this Sequence
	// followed by the digit itself. The digit is always between 0 and 9.
	All() iter.Seq2[int, int]

	// Iterator returns a function that generates the digits in this
	// Sequence along with their zero based positions from beginning to end.
	// If there are no more digits, returned function returns false.
	Iterator() func() (Digit, bool)

	// WithStart returns a view of this Sequence that only has digits with
	// zero based positions greater than or equal to start.
	WithStart(start int) Sequence

	// WithEnd returns a view of this Sequence that only has digits with
	// zero based positions less than end.
	WithEnd(end int) FiniteSequence

	private()
}

// FiniteSequence represents a Sequence of finite length.
type FiniteSequence interface {
	Sequence

	// Backward returns the 0 based position of each digit in this
	// FiniteSequence followed by the digit itself from end to beginning.
	// The digit is always between 0 and 9.
	Backward() iter.Seq2[int, int]

	// Reverse returns a function that generates the digits in this
	// FiniteSequence along with their zero based positions from end to
	// beginning. When there are no more digits, returned function
	// returns false.
	Reverse() func() (Digit, bool)

	// FiniteWithStart works like WithStart except that it returns a
	// FiniteSequence.
	FiniteWithStart(start int) FiniteSequence
}

// Fprint prints digits of s to w. Unless using advanced functionality,
// prefer Fwrite, Write, and Swrite to Fprint, Print, and Sprint.
// Fprint returns the number of bytes written and any error encountered.
// p contains the positions of the digits to print.
// For options, the default is 50 digits per row, 5 digits per column,
// show digit count, period (.) for missing digits, don't write a trailing
// line feed, and show the leading decimal point.
func Fprint(w io.Writer, s Sequence, p Positions, options ...Option) (
	written int, err error) {
	settings := &printerSettings{
		digitsPerRow:    50,
		digitsPerColumn: 5,
		showCount:       true,
		missingDigit:    '.',
		leadingDecimal:  true,
	}
	printer := newPrinter(w, p.End(), mutateSettings(options, settings))
	fromSequenceWithPositions(s, p, printer)
	printer.Finish()
	return printer.BytesWritten(), printer.Err()
}

// Fwrite writes all the digits of s to w. Fwrite returns the number of bytes
// written and any error encountered. For options, the default is 50 digits
// per row, 5 digits per column, show digit count, period (.) for missing
// digits, write a trailing line feed, and don't show the leading decimal
// point.
func Fwrite(w io.Writer, s FiniteSequence, options ...Option) (
	written int, err error) {
	settings := &printerSettings{
		digitsPerRow:     50,
		digitsPerColumn:  5,
		showCount:        true,
		missingDigit:     '.',
		trailingLineFeed: true,
	}
	printer := newPrinter(w, endOf(s), mutateSettings(options, settings))
	consume2.FromGenerator[Digit](s.Iterator(), printer)
	printer.Finish()
	return printer.BytesWritten(), printer.Err()
}

// Sprint works like Fprint and prints digits of s to a string.
func Sprint(s Sequence, p Positions, options ...Option) string {
	var builder strings.Builder
	Fprint(&builder, s, p, options...)
	return builder.String()
}

// Swrite works like Fwrite and writes all the digits of s to returned string.
func Swrite(s FiniteSequence, options ...Option) string {
	var builder strings.Builder
	Fwrite(&builder, s, options...)
	return builder.String()
}

// Print works like Fprint and prints digits of s to stdout.
func Print(s Sequence, p Positions, options ...Option) (
	written int, err error) {
	return Fprint(os.Stdout, s, p, options...)
}

// Write works like Fwrite and writes all the digits of s to stdout.
func Write(s FiniteSequence, options ...Option) (
	written int, err error) {
	return Fwrite(os.Stdout, s, options...)
}

// DigitsToString returns all the digits in s as a string.
func DigitsToString(s FiniteSequence) string {
	var sb strings.Builder
	for _, digit := range s.All() {
		sb.WriteByte('0' + byte(digit))
	}
	return sb.String()
}

func endOf(s FiniteSequence) int {
	for index := range s.Backward() {
		return index + 1
	}
	return 0
}

func fromSequenceWithPositions(
	s Sequence, p Positions, consumer consume2.Consumer[Digit]) {
	for pr := range p.All() {
		consume2.FromGenerator(
			s.WithStart(pr.Start).WithEnd(pr.End).Iterator(), consumer)
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

func opaqueSequence(s Sequence) Sequence {
	if _, ok := s.(*opqSequence); ok {
		return s
	}
	return &opqSequence{Sequence: s}
}

type opqSequence struct {
	Sequence
}

func (s *opqSequence) WithStart(start int) Sequence {
	result := s.Sequence.WithStart(start)
	if result == s.Sequence {
		return s
	}
	return opaqueSequence(result)
}

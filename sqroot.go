// Package sqroot calculates square roots to arbitrary precision.
package sqroot

import (
	"fmt"
	"io"
	"math/big"
	"os"
	"strconv"
	"strings"

	"github.com/keep94/consume2"
)

var (
	oneHundred = big.NewInt(100)
	two        = big.NewInt(2)
	one        = big.NewInt(1)
	ten        = big.NewInt(10)
)

// Option represents a Print Mantissa option.
type Option interface {
	mutate(p *printerSettings)
}

// DigitsPerRow sets the number of digits per row. The default is
// zero, which means no separate rows.
func DigitsPerRow(count int) Option {
	return optionFunc(func(p *printerSettings) {
		p.digitsPerRow = count
	})
}

// DigitsPerColumn sets the number of digits per column. The default is
// zero, which means no separate columns.
func DigitsPerColumn(count int) Option {
	return optionFunc(func(p *printerSettings) {
		p.digitsPerColumn = count
	})
}

// ShowCount shows the digit count in the left margin if on is true. The
// default is false.
func ShowCount(on bool) Option {
	return optionFunc(func(p *printerSettings) {
		p.showCount = on
	})
}

// Mantissa represents the mantissa of a square root. Non-nil Mantissas are
// between 0.1 inclusive and 1.0 exclusive. A nil Mantissa means 0.
type Mantissa func(consumer consume2.Consumer[int])

// Format prints mantissas with the f verb. Format supports width, precision,
// and the '-' flag for left justification.
func (m Mantissa) Format(state fmt.State, verb rune) {
	if verb != 'f' {
		fmt.Fprintf(state, "%%!%c(mantissa)", verb)
		return
	}
	precision, precisionOk := state.Precision()
	if !precisionOk {
		precision = 16
	}
	width, widthOk := state.Width()
	if !widthOk {
		m.printWithPrecision(state, precision)
		return
	}
	var builder strings.Builder
	m.printWithPrecision(&builder, precision)
	field := builder.String()
	if !state.Flag('-') && len(field) < width {
		fmt.Fprintf(state, "%s", strings.Repeat(" ", width-len(field)))
	}
	fmt.Fprint(state, field)
	if state.Flag('-') && len(field) < width {
		fmt.Fprintf(state, "%s", strings.Repeat(" ", width-len(field)))
	}
}

// Send sends the digits to the right of decimal point to consumer. If this
// Mantissa is nil, which means 0, Send sends no digits to consumer.
func (m Mantissa) Send(consumer consume2.Consumer[int]) {
	if m != nil {
		m(consumer)
	}
}

// Print prints this Mantissa to stdout. Print returns the number of bytes
// written and any error encountered.
func (m Mantissa) Print(maxDigits int, options ...Option) (n int, err error) {
	return m.Fprint(os.Stdout, maxDigits, options...)
}

// Fprint prints this Mantissa to w. Fprint returns the number of bytes
// written and any error encountered.
func (m Mantissa) Fprint(w io.Writer, maxDigits int, options ...Option) (
	n int, err error) {
	if m == nil {
		return fmt.Fprint(w, "0")
	}
	p := newPrinter(w, maxDigits, options)
	m.Send(p)
	return p.byteCount, p.err
}

func (m Mantissa) printWithPrecision(w io.Writer, precision int) {
	if precision == 0 {
		fmt.Fprint(w, "0")
		return
	}
	if m == nil {
		fmt.Fprint(w, "0.")
		fmt.Fprint(w, strings.Repeat("0", precision))
		return
	}
	p := newPrinter(w, precision, nil)
	m.Send(p)
	digitCount := p.index
	fmt.Fprint(w, strings.Repeat("0", precision-digitCount))
}

// SquareRoot returns the square root of radican * 10^rexp. The return value
// is of the form mantissa * 10^exp. mantissa is between 0.1 inclusive
// and 1.0 exclusive. If radican is 0, SquareRoot returns 0 for exp and
// gives a single zero digit for the mantissa.
func SquareRoot(radican *big.Int, rexp int) (mantissa Mantissa, exp int) {
	if radican.Sign() < 0 {
		panic("radican must be non-negative")
	}
	if radican.Sign() == 0 {
		return
	}
	if rexp%2 != 0 {
		radican = new(big.Int).Mul(radican, ten)
		rexp--
	}
	radicanDigits, doubleZeroCount := base100(radican)
	exp = len(radicanDigits) + doubleZeroCount + rexp/2
	mantissa = func(consumer consume2.Consumer[int]) {
		squareRoot(radicanDigits, consumer)
	}
	return
}

type printer struct {
	writer          io.Writer
	maxDigits       int
	indentation     string
	digitCountSpec  string
	digitsPerRow    int
	digitsPerColumn int
	index           int
	indexInRow      int
	byteCount       int
	err             error
}

func newPrinter(
	writer io.Writer, maxDigits int, options []Option) *printer {
	settings := &printerSettings{}
	for _, option := range options {
		option.mutate(settings)
	}
	indentation, digitCountSpec := computeIndentation(
		settings.digitCountWidth(maxDigits))
	return &printer{
		writer:          writer,
		maxDigits:       maxDigits,
		indentation:     indentation,
		digitCountSpec:  digitCountSpec,
		digitsPerRow:    settings.digitsPerRow,
		digitsPerColumn: settings.digitsPerColumn,
	}
}

func (p *printer) CanConsume() bool {
	return p.err == nil && p.index < p.maxDigits
}

func (p *printer) Consume(digit int) {
	if !p.CanConsume() {
		return
	}
	if p.index == 0 {
		n, err := fmt.Fprintf(p.writer, "%s0.", p.indentation)
		if !p.updateByteCount(n, err) {
			return
		}
	} else if p.digitsPerRow > 0 && p.index%p.digitsPerRow == 0 {
		n, err := fmt.Fprintln(p.writer)
		if !p.updateByteCount(n, err) {
			return
		}
		if p.digitCountSpec != "" {
			n, err := fmt.Fprintf(p.writer, p.digitCountSpec, p.index)
			if !p.updateByteCount(n, err) {
				return
			}
		}
		n, err = fmt.Fprint(p.writer, "  ")
		if !p.updateByteCount(n, err) {
			return
		}
		p.indexInRow = 0
	} else if p.digitsPerColumn > 0 && p.indexInRow%p.digitsPerColumn == 0 {
		n, err := fmt.Fprint(p.writer, " ")
		if !p.updateByteCount(n, err) {
			return
		}
	}
	n, err := fmt.Fprintf(p.writer, "%d", digit)
	if !p.updateByteCount(n, err) {
		return
	}
	p.index++
	p.indexInRow++
}

func (p *printer) updateByteCount(n int, err error) bool {
	p.byteCount += n
	p.err = err
	return err == nil
}

func squareRoot(radicanDigits []*big.Int, consumer consume2.Consumer[int]) {
	radicanDigitsIdx := len(radicanDigits)
	incr := big.NewInt(1)
	remainder := big.NewInt(0)
	for consumer.CanConsume() {
		if radicanDigitsIdx == 0 && remainder.Sign() == 0 {
			return
		}
		remainder.Mul(remainder, oneHundred)
		if radicanDigitsIdx > 0 {
			radicanDigitsIdx--
			remainder.Add(remainder, radicanDigits[radicanDigitsIdx])
		}
		digit := 0
		for remainder.Cmp(incr) >= 0 {
			remainder.Sub(remainder, incr)
			digit++
			incr.Add(incr, two)
		}
		consumer.Consume(digit)
		incr.Sub(incr, one).Mul(incr, ten).Add(incr, one)
	}
}

func base100(radican *big.Int) (result []*big.Int, doubleZeroCount int) {
	radican = new(big.Int).Set(radican)
	trailingZeros := true
	for radican.Sign() > 0 {
		_, m := radican.DivMod(radican, oneHundred, new(big.Int))
		if trailingZeros && m.Sign() == 0 {
			doubleZeroCount++
		} else {
			result = append(result, m)
			trailingZeros = false
		}
	}
	return
}

type optionFunc func(p *printerSettings)

func (o optionFunc) mutate(p *printerSettings) {
	o(p)
}

type printerSettings struct {
	digitsPerRow    int
	digitsPerColumn int
	showCount       bool
}

func (p *printerSettings) digitCountWidth(maxDigits int) int {
	if !p.showCount || p.digitsPerRow <= 0 {
		return 0
	}
	if maxDigits <= p.digitsPerRow {
		return 0
	}
	maxCounter := ((maxDigits - 1) / p.digitsPerRow) * p.digitsPerRow
	return len(strconv.Itoa(maxCounter))
}

func computeIndentation(width int) (
	indentation string, digitCountSpec string) {
	if width <= 0 {
		return
	}
	indentation = strings.Repeat(" ", width)
	digitCountSpec = fmt.Sprintf("%%%dd", width)
	return
}

// Package sqroot calculates square roots to arbitrary precision.
package sqroot

import (
	"fmt"
	"io"
	"math/big"
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

// SquareRoot returns the square root of radican * 10^rexp. The return value
// is of the form mantissa * 10^exp. mantissa is between 0.1 inclusive
// and 1.0 exclusive. mantissa is actually a function that sends the digits
// of the mantissa to the passed consumer. If radican is 0, SquareRoot returns
// 0 for exp and gives a single zero digit for the mantissa.
func SquareRoot(radican *big.Int, rexp int) (
	mantissa func(consumer consume2.Consumer[int]), exp int) {
	if radican.Sign() < 0 {
		panic("radican must be non-negative")
	}
	if radican.Sign() == 0 {
		return zeroDigit, 0
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

// Option represents a Printer option.
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

// Printer is a Consumer[int] that prints out the mantissa of square roots.
type Printer struct {
	writer          io.Writer
	maxDigits       int
	indentation     string
	digitCountSpec  string
	digitsPerRow    int
	digitsPerColumn int
	index           int
	indexInRow      int
}

// NewPrinter creates a new Printer that sends digits to writer. maxDigits
// is the maximum number of digits to send.
func NewPrinter(writer io.Writer, maxDigits int, options ...Option) *Printer {
	settings := &printerSettings{}
	for _, option := range options {
		option.mutate(settings)
	}
	indentation, digitCountSpec := computeIndentation(
		settings.digitCountWidth(maxDigits))
	return &Printer{
		writer:          writer,
		maxDigits:       maxDigits,
		indentation:     indentation,
		digitCountSpec:  digitCountSpec,
		digitsPerRow:    settings.digitsPerRow,
		digitsPerColumn: settings.digitsPerColumn,
	}
}

// CanConsume returns true if this Printer can consume a digit.
func (p *Printer) CanConsume() bool {
	return p.index < p.maxDigits
}

// Consume sends a single digit to the underlying io.Writer of this Printer.
func (p *Printer) Consume(digit int) {
	if !p.CanConsume() {
		return
	}
	if p.index == 0 {
		fmt.Fprintf(p.writer, "%s0.", p.indentation)
	} else if p.digitsPerRow > 0 && p.index%p.digitsPerRow == 0 {
		fmt.Fprintln(p.writer)
		if p.digitCountSpec != "" {
			fmt.Fprintf(p.writer, p.digitCountSpec, p.index)
		}
		fmt.Fprint(p.writer, "  ")
		p.indexInRow = 0
	} else if p.digitsPerColumn > 0 && p.indexInRow%p.digitsPerColumn == 0 {
		fmt.Fprint(p.writer, " ")
	}
	fmt.Fprintf(p.writer, "%d", digit)
	p.index++
	p.indexInRow++
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

func zeroDigit(consumer consume2.Consumer[int]) {
	consumer.Consume(0)
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

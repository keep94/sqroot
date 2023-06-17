package sqroot

import (
	"fmt"
	"io"
	"strconv"
	"strings"
)

type digit struct {

	// The 0 based position of the digit.
	Position int

	// The value of the digit. Always between 0 and 9.
	Value int
}

type printer struct {
	rawPrinter
	missingDigit rune
}

func newPrinter(
	writer io.Writer, maxDigits int, settings *printerSettings) *printer {
	var result printer
	result.Init(writer, maxDigits, settings)
	result.missingDigit = settings.missingDigit
	return &result
}

func (p *printer) Consume(d digit) {
	if p.index < d.Position {
		if p.digitsPerRow > 0 && p.digitCountSpec != "" {
			p.skipRowsFor(d.Position)
		}
		for p.index < d.Position {
			p.rawPrinter.Consume(p.missingDigit)
		}
	}
	p.rawPrinter.Consume('0' + rune(d.Value))
}

func (p *printer) skipRowsFor(nextPosit int) {
	currentRow := p.index / p.digitsPerRow
	nextRow := nextPosit / p.digitsPerRow
	if p.index%p.digitsPerRow == 0 {
		p.skipRows(nextRow - currentRow)
	} else if nextRow > currentRow {
		p.skipRows(nextRow - currentRow - 1)
	}
}

type rawPrinter struct {
	writer          io.Writer
	indentation     string
	digitCountSpec  string
	digitsPerRow    int
	digitsPerColumn int
	index           int
	indexInRow      int
	byteCount       int
	err             error
}

func (p *rawPrinter) Init(
	writer io.Writer, maxDigits int, settings *printerSettings) {
	indentation, digitCountSpec := computeIndentation(
		settings.digitCountWidth(maxDigits))
	*p = rawPrinter{
		writer:          writer,
		indentation:     indentation,
		digitCountSpec:  digitCountSpec,
		digitsPerRow:    settings.digitsPerRow,
		digitsPerColumn: settings.digitsPerColumn,
	}
}

func (p *rawPrinter) CanConsume() bool {
	return p.err == nil
}

func (p *rawPrinter) Consume(digit rune) {
	if !p.CanConsume() {
		return
	}
	if p.index == 0 {
		n, err := fmt.Fprintf(p.writer, "%s0.", p.indentation)
		if !p.updateByteCount(n, err) {
			return
		}
	} else if p.digitsPerRow > 0 && p.index%p.digitsPerRow == 0 {
		if p.byteCount > 0 {
			n, err := fmt.Fprintln(p.writer)
			if !p.updateByteCount(n, err) {
				return
			}
		}
		if p.digitCountSpec != "" {
			n, err := fmt.Fprintf(p.writer, p.digitCountSpec, p.index)
			if !p.updateByteCount(n, err) {
				return
			}
		}
		n, err := fmt.Fprint(p.writer, "  ")
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
	n, err := fmt.Fprintf(p.writer, "%c", digit)
	if !p.updateByteCount(n, err) {
		return
	}
	p.index++
	p.indexInRow++
}

func (p *rawPrinter) skipRows(rowsToSkip int) {
	p.index += rowsToSkip * p.digitsPerRow
}

func (p *rawPrinter) updateByteCount(n int, err error) bool {
	p.byteCount += n
	p.err = err
	return err == nil
}

type printerSettings struct {
	digitsPerRow    int
	digitsPerColumn int
	showCount       bool
	missingDigit    rune
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

type formatter struct {
	writer          io.Writer
	sigDigits       int // invariant sigDigits >= exponent
	exponent        int
	exactDigitCount bool
	index           int
}

func newFormatter(
	w io.Writer, sigDigits, exponent int, exactDigitCount bool) *formatter {
	if sigDigits < exponent {
		panic("sigDigits must be >= exponent")
	}
	return &formatter{
		writer:          w,
		sigDigits:       sigDigits,
		exponent:        exponent,
		exactDigitCount: exactDigitCount,
	}
}

func (f *formatter) CanConsume() bool {
	return f.index < f.sigDigits
}

func (f *formatter) Consume(digit int) {
	if !f.CanConsume() {
		return
	}
	f.add(digit)
}

func (f *formatter) Finish() {
	maxDigits := f.sigDigits
	if !f.exactDigitCount {
		maxDigits = f.exponent
	}
	for f.index < maxDigits {
		f.add(0)
	}
	// If we haven't written anything yet
	if f.index == 0 {
		count := -f.exponent
		if f.exactDigitCount {
			count = f.sigDigits - f.exponent
		}
		f.addLeadingZeros(count)
	}
}

func (f *formatter) add(digit int) {
	if f.index == 0 && f.exponent <= 0 {
		f.addLeadingZeros(-f.exponent)
	}
	if f.index == f.exponent {
		fmt.Fprint(f.writer, ".")
	}
	fmt.Fprint(f.writer, digit)
	f.index++
}

func (f *formatter) addLeadingZeros(count int) {
	fmt.Fprint(f.writer, "0")
	if count <= 0 {
		return
	}
	fmt.Fprint(f.writer, ".")
	fmt.Fprint(f.writer, strings.Repeat("0", count))
}

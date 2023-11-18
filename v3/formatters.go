package sqroot

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// Digit represents a digit and a zero based position.
type Digit struct {

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

func (p *printer) Consume(d Digit) {
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
	cWriter         *countingWriter
	writer          *bufio.Writer
	indentation     string
	digitCountSpec  string
	digitsPerRow    int
	digitsPerColumn int
	index           int
	indexInRow      int
	err             error
}

func (p *rawPrinter) Init(
	writer io.Writer, maxDigits int, settings *printerSettings) {
	cWriter := &countingWriter{delegate: writer}
	var bWriter *bufio.Writer
	if settings.bufferSize <= 0 {
		bWriter = bufio.NewWriter(cWriter)
	} else {
		bWriter = bufio.NewWriterSize(cWriter, settings.bufferSize)
	}
	indentation, digitCountSpec := computeIndentation(
		settings.digitCountWidth(maxDigits))
	*p = rawPrinter{
		cWriter:         cWriter,
		writer:          bWriter,
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
		_, p.err = fmt.Fprintf(p.writer, "%s0.", p.indentation)
		if p.err != nil {
			return
		}
	} else if p.digitsPerRow > 0 && p.index%p.digitsPerRow == 0 {
		if p.BytesWritten()+p.bytesBuffered() > 0 {
			_, p.err = fmt.Fprintln(p.writer)
			if p.err != nil {
				return
			}
		}
		if p.digitCountSpec != "" {
			_, p.err = fmt.Fprintf(p.writer, p.digitCountSpec, p.index)
			if p.err != nil {
				return
			}
		}
		_, p.err = p.writer.WriteString("  ")
		if p.err != nil {
			return
		}
		p.indexInRow = 0
	} else if p.digitsPerColumn > 0 && p.indexInRow%p.digitsPerColumn == 0 {
		p.err = p.writer.WriteByte(' ')
		if p.err != nil {
			return
		}
	}
	_, p.err = p.writer.WriteRune(digit)
	if p.err != nil {
		return
	}
	p.index++
	p.indexInRow++
}

func (p *rawPrinter) Finish() {
	err := p.writer.Flush()
	if p.err == nil {
		p.err = err
	}
}

func (p *rawPrinter) BytesWritten() int {
	return p.cWriter.bytesWritten
}

func (p *rawPrinter) Err() error {
	return p.err
}

func (p *rawPrinter) bytesBuffered() int {
	return p.writer.Buffered()
}

func (p *rawPrinter) skipRows(rowsToSkip int) {
	p.index += rowsToSkip * p.digitsPerRow
}

type printerSettings struct {
	digitsPerRow    int
	digitsPerColumn int
	showCount       bool
	missingDigit    rune
	bufferSize      int
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
	writer          *bufio.Writer
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
		writer:          bufio.NewWriter(w),
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
	f.writer.Flush()
}

func (f *formatter) add(digit int) {
	if f.index == 0 && f.exponent <= 0 {
		f.addLeadingZeros(-f.exponent)
	}
	if f.index == f.exponent {
		f.writer.WriteByte('.')
	}
	f.writer.WriteByte('0' + byte(digit))
	f.index++
}

func (f *formatter) addLeadingZeros(count int) {
	f.writer.WriteByte('0')
	if count <= 0 {
		return
	}
	f.writer.WriteByte('.')
	for i := 0; i < count; i++ {
		f.writer.WriteByte('0')
	}
}

type countingWriter struct {
	delegate     io.Writer
	bytesWritten int
}

func (c *countingWriter) Write(p []byte) (n int, err error) {
	n, err = c.delegate.Write(p)
	c.bytesWritten += n
	return
}

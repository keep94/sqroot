package sqroot

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

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

func (p *printer) Consume(posit, digit int) {
	if p.index < posit {
		if p.digitsPerRow > 0 && p.rowStarter.CountOn() {
			p.skipRowsFor(posit)
		}
		for p.index < posit {
			p.rawPrinter.Consume(p.missingDigit)
		}
	}
	p.rawPrinter.Consume('0' + rune(digit))
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

type rowStarter interface {
	Start(w *bufio.Writer, index int) error
	CountOn() bool
}

type countOnStarter struct {
	zeroString    string
	nonZeroString string
}

func (c *countOnStarter) Start(w *bufio.Writer, index int) error {
	if index == 0 {
		_, err := w.WriteString(c.zeroString)
		return err
	}
	_, err := fmt.Fprintf(w, c.nonZeroString, index)
	return err
}

func (c *countOnStarter) CountOn() bool { return true }

type countOffStarter struct {
	zeroString    string
	nonZeroString string
}

func (c *countOffStarter) Start(w *bufio.Writer, index int) error {
	if index == 0 {
		_, err := w.WriteString(c.zeroString)
		return err
	}
	_, err := w.WriteString(c.nonZeroString)
	return err
}

func (c *countOffStarter) CountOn() bool { return false }

type rawPrinter struct {
	cWriter          *countingWriter
	writer           *bufio.Writer
	rowStarter       rowStarter
	digitsPerRow     int
	digitsPerColumn  int
	trailingLineFeed bool
	index            int
	indexInRow       int
	err              error
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
	*p = rawPrinter{
		cWriter:          cWriter,
		writer:           bWriter,
		rowStarter:       settings.computeRowStarter(maxDigits),
		digitsPerRow:     settings.digitsPerRow,
		digitsPerColumn:  settings.digitsPerColumn,
		trailingLineFeed: settings.trailingLineFeed,
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
		p.err = p.rowStarter.Start(p.writer, 0)
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
		p.err = p.rowStarter.Start(p.writer, p.index)
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
	if p.err == nil && p.trailingLineFeed {
		_, p.err = fmt.Fprintln(p.writer)
	}
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
	digitsPerRow     int
	digitsPerColumn  int
	showCount        bool
	missingDigit     rune
	bufferSize       int
	trailingLineFeed bool
	leadingDecimal   bool
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

func (p *printerSettings) computeRowStarter(maxDigits int) rowStarter {
	width := p.digitCountWidth(maxDigits)
	if width <= 0 {
		if p.leadingDecimal {
			return &countOffStarter{zeroString: "0.", nonZeroString: "  "}
		} else if p.showCount {
			return &countOffStarter{zeroString: "0  ", nonZeroString: "   "}
		} else {
			return &countOffStarter{}
		}
	}
	if p.leadingDecimal {
		return &countOnStarter{
			zeroString:    strings.Repeat(" ", width) + "0.",
			nonZeroString: fmt.Sprintf("%%%dd  ", width),
		}
	}
	return &countOnStarter{
		zeroString:    strings.Repeat(" ", width-1) + "0  ",
		nonZeroString: fmt.Sprintf("%%%dd  ", width),
	}
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

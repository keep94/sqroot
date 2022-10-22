// Package sqroot calculates square roots to arbitrary precision.
package sqroot

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math/big"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/keep94/consume2"
)

const (
	fPrecision                 = 6
	gPrecision                 = 16
	digitsBinaryVersion        = 187
	unmarshalTextUnexpectedEnd = "sqroot: UnmarshalText hit unexpected end of text"
)

// Option represents an option for the Print, Fprint, and Sprint methods of
// Mantissa and Digits.
type Option interface {
	mutate(p *printerSettings)
}

// DigitsPerRow sets the number of digits per row. Zero means no separate rows.
func DigitsPerRow(count int) Option {
	return optionFunc(func(p *printerSettings) {
		p.digitsPerRow = count
	})
}

// DigitsPerColumn sets the number of digits per column. Zero means no
// separate columns.
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

// Positions represents a set of zero based positions for which to fetch
// digits. The zero value contains no positions and is ready to use.
type Positions struct {
	ranges map[int]int
	limit  int
}

// Add adds a posit to this instance and returns this instance for chaining.
// Add panics if posit is negative.
func (p *Positions) Add(posit int) *Positions {
	return p.AddRange(posit, posit+1)
}

// AddRange adds a range of positions to this instance and returns this
// instance for chaining. The range is between start inclusive and end
// exclusive. If end <= start, AddRange is a no-op. AddRange panics if
// start is negative.
func (p *Positions) AddRange(start, end int) *Positions {
	if start < 0 {
		panic("start must be non-negative")
	}
	oldValue := p.ranges[start]
	if end-start <= oldValue {
		return p
	}
	if p.ranges == nil {
		p.ranges = make(map[int]int)
	}
	p.ranges[start] = end - start
	if end > p.limit {
		p.limit = end
	}
	return p
}

// Copy returns a new instance that is a copy of this instance.
func (p *Positions) Copy() *Positions {
	if len(p.ranges) == 0 {
		return &Positions{}
	}
	result := make(map[int]int, len(p.ranges))
	for k, v := range p.ranges {
		result[k] = v
	}
	return &Positions{ranges: result, limit: p.limit}
}

func (p *Positions) filter() *positionsFilter {
	starts := make([]int, 0, len(p.ranges))
	for k := range p.ranges {
		starts = append(starts, k)
	}
	sort.Ints(starts)
	return &positionsFilter{starts: starts, ranges: p.ranges}
}

// Sequence represents a sequence of possibly non contiguous digits.
// That is, a Sequence may have holes. For instance, a Sequence could
// be 375XXX695. The 0th digit is a 3; the 1st digit is a 7; the 2nd digit
// is a 5; the 3rd, 4th, and 5th digits are unknown; the 6th digit is a 6;
// the 7th digit is a 9; the 8th digit is a 5.
// Both Digits and Mantissa implement Sequence.
type Sequence interface {
	positDigitIter() func() positDigit
}

// Find returns a function that returns the next zero based index of the
// match for pattern in s. If s has a finite number of digits and there
// are no more matches for pattern, the returned function returns -1.
// Pattern is a sequence of digits between 0 and 9.
func Find(s Sequence, pattern []int) func() int {
	if len(pattern) == 0 {
		return zeroPattern(s.positDigitIter())
	}
	patternCopy := make([]int, len(pattern))
	copy(patternCopy, pattern)
	return kmp(s.positDigitIter(), patternCopy)
}

// FindFirst finds the zero based index of the first match of pattern in s.
// FindFirst returns -1 if pattern is not found only if s has a finite number
// of digits. If s has an infinite number of digits and pattern is not found,
// FindFirst will run forever. pattern is a sequence of digits between 0 and 9.
func FindFirst(s Sequence, pattern []int) int {
	iter := find(s, pattern)
	return iter()
}

// FindFirstN works like FindFirst but it finds the first n matches and
// returns the zero based index of each match. If s has a finite
// number of digits, FindFirstN may return fewer than n matches.
// Like FindFirst, FindFirstN may run forever if s has an infinite
// number of digits, and there are not n matches available.
// pattern is a sequence of digits between 0 and 9.
func FindFirstN(s Sequence, pattern []int, n int) []int {
	var result []int
	iter := find(s, pattern)
	for index := iter(); index != -1 && len(result) < n; index = iter() {
		result = append(result, index)
	}
	return result
}

// FindAll finds all the matches of pattern in s and returns the zero based
// index of each match. If s has an infinite number of digits, FindAll will
// run forever. pattern is a sequence of digits between 0 and 9.
func FindAll(s Sequence, pattern []int) []int {
	var result []int
	iter := find(s, pattern)
	for index := iter(); index != -1; index = iter() {
		result = append(result, index)
	}
	return result
}

// GetDigits gets the digits from s found at the zero based positions
// in p.
func GetDigits(s Sequence, p *Positions) Digits {
	return asDigits(newPart(s, p))
}

// Digits holds the digits found at selected positions of a Mantissa so
// that they can be quickly retrieved. The zero value is no digits.
// Digits implements Sequence.
type Digits struct {
	digits map[int]int
	posits []int
}

// MarshalBinary implements the encoding.BinaryMarshaler interface.
func (d Digits) MarshalBinary() ([]byte, error) {
	iter := d.Iterator()
	nextPosit := 0
	result := []byte{digitsBinaryVersion}
	state := 0
	pair := uint64(0)
	for posit := iter(); posit != -1; posit = iter() {
		delta := posit - nextPosit
		if delta > 0 {
			if state == 1 {
				result = binary.AppendUvarint(result, 100+pair)
				state = 0
				pair = 0
			}
			result = binary.AppendUvarint(result, uint64(delta)+109)
		}
		nextPosit = posit + 1
		pair = 10*pair + uint64(d.At(posit))
		if state == 1 {
			result = binary.AppendUvarint(result, pair)
			pair = 0
		}
		state = 1 - state
	}
	if state == 1 {
		result = binary.AppendUvarint(result, 100+pair)
	}
	return result, nil
}

// MarshalText implements the encoding.TextMarshaler interface.
func (d Digits) MarshalText() ([]byte, error) {
	iter := d.Iterator()
	nextPosit := 0
	var result []byte
	for posit := iter(); posit != -1; posit = iter() {
		if posit > nextPosit {
			result = append(result, byte('['))
			result = strconv.AppendInt(result, int64(posit), 10)
			result = append(result, byte(']'))
			nextPosit = posit
		}
		result = append(result, byte('0'+d.At(posit)))
		nextPosit++
	}
	return result, nil
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface.
func (d *Digits) UnmarshalBinary(b []byte) error {
	if len(b) == 0 || b[0] != digitsBinaryVersion {
		return errors.New("sqroot: Bad Digits Binary Version")
	}
	var builder digitsBuilder
	posit := 0
	reader := bytes.NewReader(b[1:])
	for reader.Len() > 0 {
		val, err := binary.ReadUvarint(reader)
		if err != nil {
			return err
		}
		if val >= 110 {
			posit += int(val - 109)
		} else if val >= 100 {
			if err := builder.AddDigit(posit, int(val-100)); err != nil {
				return err
			}
			posit++
		} else {
			if err := builder.AddDigit(posit, int(val/10)); err != nil {
				return err
			}
			posit++
			if err := builder.AddDigit(posit, int(val%10)); err != nil {
				return err
			}
			posit++
		}
	}
	*d = builder.Build()
	return nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (d *Digits) UnmarshalText(text []byte) error {
	var builder digitsBuilder
	posit := 0
	i := 0
	var err error
	for i < len(text) {
		if text[i] == '[' {
			posit, i, err = readPositiveInt(text, i+1)
			if err != nil {
				return err
			}
		}
		digit := int(text[i] - '0')
		if err := builder.AddDigit(posit, digit); err != nil {
			return err
		}
		posit++
		i++
	}
	*d = builder.Build()
	return nil
}

// At returns the digit between 0 and 9 at the given zero based position.
// If the digit at posit is unknown or if posit is negative, At returns -1.
func (d Digits) At(posit int) int {
	digit, ok := d.digits[posit]
	if !ok {
		return -1
	}
	return digit
}

// Iterator returns a function that generates all the zero based positions
// in this instance from lowest to highest. When there are no more positions,
// the returned function returns -1.
func (d Digits) Iterator() func() int {
	return d.IteratorAt(0)
}

// Reverse returns a function that generates all the zero based positions
// in this instance from highest to lowest. When there are no more positions,
// the returned function returns -1.
func (d Digits) Reverse() func() int {
	return d.ReverseAt(d.Max() + 1)
}

// IteratorAt returns a function that generates all the zero based
// positions in this instance from lowest to highest starting at posit.
// When there are no more positions, the returned function returns -1.
func (d Digits) IteratorAt(posit int) func() int {
	index := sort.Search(
		len(d.posits), func(x int) bool { return d.posits[x] >= posit })
	return func() int {
		if index == len(d.posits) {
			return -1
		}
		result := d.posits[index]
		index++
		return result
	}
}

// ReverseAt returns a function that generates all the zero based positions
// in this instance from highest to lowest that come before posit. When
// there are no more positions, the returned function returns -1.
func (d Digits) ReverseAt(posit int) func() int {
	index := sort.Search(
		len(d.posits), func(x int) bool { return d.posits[x] >= posit }) - 1
	return func() int {
		if index == -1 {
			return -1
		}
		result := d.posits[index]
		index--
		return result
	}
}

// Min returns the minimum position in this instance. If this instance
// is empty, Min returns -1.
func (d Digits) Min() int {
	if len(d.posits) == 0 {
		return -1
	}
	return d.posits[0]
}

// Max returns the maximum position in this instance. If this instance
// is empty, Max returns -1.
func (d Digits) Max() int {
	if len(d.posits) == 0 {
		return -1
	}
	return d.posits[len(d.posits)-1]
}

// Len returns the number of digits in this instance.
func (d Digits) Len() int {
	return len(d.posits)
}

// Print works like Fprint printing this instance to stdout.
func (d Digits) Print(options ...Option) (n int, err error) {
	return d.Fprint(os.Stdout, options...)
}

// Sprint works like Fprint printing this instance to the returned string.
func (d Digits) Sprint(options ...Option) string {
	var builder strings.Builder
	d.Fprint(&builder, options...)
	return builder.String()
}

// Fprint prints this instance to w and returns number of bytes written
// and any error encountered. For options, the default is 50 digits per
// row, 5 digits per column, and show digit count.
func (d Digits) Fprint(w io.Writer, options ...Option) (n int, err error) {
	settings := &printerSettings{
		digitsPerRow:    50,
		digitsPerColumn: 5,
		showCount:       true,
	}
	return fprint(w, d, mutateSettings(options, settings))
}

func (d Digits) limit() int {
	if len(d.posits) == 0 {
		return 0
	}
	return d.posits[len(d.posits)-1] + 1
}

func (d Digits) positDigitIter() func() positDigit {
	index := 0
	return func() positDigit {
		if index == len(d.posits) {
			return invalidPositDigit
		}
		posit := d.posits[index]
		result := positDigit{Posit: posit, Digit: d.digits[posit]}
		index++
		return result
	}
}

func readPositiveInt(text []byte, i int) (int, int, error) {
	result := 0
	for i < len(text) {
		if text[i] == ']' {
			if i+1 == len(text) {
				return 0, 0, errors.New(unmarshalTextUnexpectedEnd)
			}
			return result, i + 1, nil
		} else if text[i] >= '0' && text[i] <= '9' {
			result = result*10 + int(text[i]-'0')
		} else {
			return 0, 0, fmt.Errorf("sqroot: UnmarshalText unexpected character in text: %c", text[i])
		}
		i++
	}
	return 0, 0, errors.New(unmarshalTextUnexpectedEnd)
}

type digitsBuilder struct {
	digits map[int]int
	posits []int
}

func (d *digitsBuilder) AddDigit(posit int, digit int) error {
	if posit < 0 {
		return fmt.Errorf(
			"sqroot: posit must be non-negative but was %d", posit)
	}
	if digit < 0 || digit > 9 {
		return fmt.Errorf(
			"sqroot: digit must be between 0 and 9 but was %d", digit)
	}
	if len(d.posits) > 0 && d.posits[len(d.posits)-1] >= posit {
		return fmt.Errorf(
			"sqroot: posit must be ever increasing was %d now %d",
			d.posits[len(d.posits)-1],
			posit,
		)
	}
	if d.digits == nil {
		d.digits = make(map[int]int)
	}
	d.digits[posit] = digit
	d.posits = append(d.posits, posit)
	return nil
}

func (d *digitsBuilder) Build() Digits {
	result := Digits{digits: d.digits, posits: d.posits}
	d.digits = nil
	d.posits = nil
	return result
}

// Mantissa represents the mantissa of a square root. Non zero Mantissas are
// between 0.1 inclusive and 1.0 exclusive. The number of digits of a
// Mantissa can be infinite. The zero value for a Mantissa corresponds to 0.
// Mantissa implements Sequence.
type Mantissa struct {
	spec mantissaSpec
}

// WithSignificant returns a new Mantissa like this one that has no more
// than limit significant digits. WithSignificant rounds the returned
// Mantissa down toward zero when necessary.
func (m Mantissa) WithSignificant(limit int) Mantissa {
	return Mantissa{spec: withLimit(m.spec, limit)}
}

// Format prints this Mantissa with the f, F, g, G, e, E verbs. The verbs work
// in the usual way except that they always round down. Because Mantissas can
// have an infinite number of digits, g with no precision shows a max of 16
// significant digits. Format supports width, precision, and the '-' flag
// for left justification. The v verb is an alias for g.
func (m Mantissa) Format(state fmt.State, verb rune) {
	formatSpec, ok := newFormatSpec(state, verb, 0)
	if !ok {
		fmt.Fprintf(state, "%%!%c(mantissa=%s)", verb, m.String())
		return
	}
	formatSpec.PrintField(state, m, 0)
}

// String returns the decimal representation of m as generated by m.Sprint(16).
func (m Mantissa) String() string {
	return m.Sprint(gPrecision)
}

// Iterator returns the digits of this Mantissa as a function. The
// first call to returned function returns the first digit of Mantissa;
// the second call returns the second digit and so forth. If returned
// function runs out of Mantissa digits, it returns -1. If this
// Mantissa is zero, the returned function always returns -1.
func (m Mantissa) Iterator() func() int {
	if m.spec == nil {
		return func() int { return -1 }
	}
	return m.spec.Iterator()
}

// Print works like Fprint and prints this Mantissa to stdout.
func (m Mantissa) Print(maxDigits int, options ...Option) (n int, err error) {
	return m.Fprint(os.Stdout, maxDigits, options...)
}

// Sprint works like Fprint and prints this Mantissa to a string.
func (m Mantissa) Sprint(maxDigits int, options ...Option) string {
	var builder strings.Builder
	m.Fprint(&builder, maxDigits, options...)
	return builder.String()
}

// Fprint prints this Mantissa to w. Fprint returns the number of bytes
// written and any error encountered. For options, the default is no
// separate rows, no separate columns, and digit count turned off.
func (m Mantissa) Fprint(w io.Writer, maxDigits int, options ...Option) (
	n int, err error) {
	if m.spec == nil || maxDigits <= 0 {
		return fmt.Fprint(w, "0")
	}
	settings := &printerSettings{}
	return fprint(
		w,
		newPart(m, new(Positions).AddRange(0, maxDigits)),
		mutateSettings(options, settings))
}

// At returns the digit at the given 0 based position in this Mantissa. If
// this Mantissa has posit or fewer digits, At returns -1. If posit is
// negative, At returns -1. When fetching digits at multiple positions, it
// is more efficient to use the GetDigits function to get a Digits instance
// than it is to call At multiple times.
func (m Mantissa) At(posit int) int {
	if posit < 0 {
		return -1
	}
	return GetDigits(m, new(Positions).Add(posit)).At(posit)
}

func (m Mantissa) positDigitIter() func() positDigit {
	iter := m.Iterator()
	digit := iter()
	index := 0
	return func() positDigit {
		if digit == -1 {
			return invalidPositDigit
		}
		result := positDigit{Posit: index, Digit: digit}
		digit = iter()
		index++
		return result
	}
}

func (m Mantissa) send(consumer consume2.Consumer[int]) {
	iter := m.Iterator()
	for consumer.CanConsume() {
		digit := iter()
		if digit == -1 {
			return
		}
		consumer.Consume(digit)
	}
}

// Number represents a square root value. The zero value for Number
// corresponds to 0. A Number is of the form mantissa * 10^exponent where
// mantissa is between 0.1 inclusive and 1.0 exclusive. Like Mantissa, a
// Number instance can represent an infinite number of digits.
type Number struct {
	mantissa Mantissa
	exponent int
}

// WithSignificant returns a Number like this one that has no more than
// limit significant digits. WithSignificant rounds the returned Number
// down toward zero when necessary.
func (n Number) WithSignificant(limit int) Number {
	m := n.mantissa.WithSignificant(limit)
	if m.spec == nil {
		return Number{}
	}
	return Number{
		mantissa: m,
		exponent: n.exponent,
	}
}

// Mantissa returns the Mantissa of this Number.
func (n Number) Mantissa() Mantissa {
	return n.mantissa
}

// Exponent returns the exponent of this Number.
func (n Number) Exponent() int {
	return n.exponent
}

// Format prints this Number with the f, F, g, G, e, E verbs. The verbs work
// in the usual way except that they always round down. Because Number can
// have an infinite number of digits, g with no precision shows a max of 16
// significant digits. Format supports width, precision, and the '-' flag
// for left justification. The v verb is an alias for g.
func (n Number) Format(state fmt.State, verb rune) {
	formatSpec, ok := newFormatSpec(state, verb, n.exponent)
	if !ok {
		fmt.Fprintf(state, "%%!%c(number=%s)", verb, n.String())
		return
	}
	formatSpec.PrintField(state, n.mantissa, n.exponent)
}

// String returns the decimal representation of n using %g.
func (n Number) String() string {
	var builder strings.Builder
	fs := formatSpec{sigDigits: gPrecision, sci: bigExponent(n.exponent)}
	fs.PrintNumber(&builder, n.mantissa, n.exponent)
	return builder.String()
}

// Sqrt returns the square root of radican. Sqrt panics if radican is
// negative.
func Sqrt(radican int64) Number {
	return sqrtFrac(big.NewInt(radican), one)
}

// SqrtRat returns the square root of num / denom. denom must be positive,
// and num must be non-negative or else SqrtRat panics.
func SqrtRat(num, denom int64) Number {
	return sqrtFrac(big.NewInt(num), big.NewInt(denom))
}

// SqrtBigInt returns the square root of radican. SqrtBigInt panics if
// radican is negative.
func SqrtBigInt(radican *big.Int) Number {
	return sqrtFrac(radican, one)
}

// SqrtBigRat returns the square root of radican. The denominator of radican
// must be positive, and the numerator must be non-negative or else SqrtBigRat
// panics.
func SqrtBigRat(radican *big.Rat) Number {
	return sqrtFrac(radican.Num(), radican.Denom())
}

func sqrtFrac(num, denom *big.Int) Number {
	num = new(big.Int).Set(num)
	denom = new(big.Int).Set(denom)
	if denom.Sign() <= 0 {
		panic("Denominator must be positive")
	}
	if num.Sign() < 0 {
		panic("Numerator must be non-negative")
	}
	if num.Sign() == 0 {
		return Number{}
	}
	exp := 0
	for num.Cmp(denom) < 0 {
		exp--
		num.Mul(num, oneHundred)
	}
	if exp < 0 {
		exp++
		num.Div(num, oneHundred)
	}
	for num.Cmp(denom) >= 0 {
		exp++
		denom.Mul(denom, oneHundred)
	}
	spec := &sqrtSpec{}
	spec.num.Set(num)
	spec.denom.Set(denom)
	return Number{exponent: exp, mantissa: Mantissa{spec: spec}}
}

type formatSpec struct {
	sigDigits       int
	exactDigitCount bool
	sci             bool
	capital         bool
}

func newFormatSpec(state fmt.State, verb rune, exponent int) (
	formatSpec, bool) {
	precision, precisionOk := state.Precision()
	var sigDigits int
	var exactDigitCount bool
	var sci bool
	switch verb {
	case 'f', 'F':
		if !precisionOk {
			precision = fPrecision
		}
		sigDigits = precision + exponent
		exactDigitCount = true
		sci = false
	case 'g', 'G', 'v':
		if !precisionOk {
			precision = gPrecision
		}
		sigDigits = precision
		if sigDigits == 0 {
			sigDigits = 1
		}
		exactDigitCount = false
		sci = sigDigits < exponent || bigExponent(exponent)
	case 'e', 'E':
		if !precisionOk {
			precision = fPrecision
		}
		sigDigits = precision
		exactDigitCount = true
		sci = true
	default:
		return formatSpec{}, false
	}
	capital := verb == 'E' || verb == 'G'
	return formatSpec{
		sigDigits:       sigDigits,
		exactDigitCount: exactDigitCount,
		sci:             sci,
		capital:         capital}, true
}

func (f formatSpec) PrintField(state fmt.State, m Mantissa, exponent int) {
	width, widthOk := state.Width()
	if !widthOk {
		f.PrintNumber(state, m, exponent)
		return
	}
	var builder strings.Builder
	f.PrintNumber(&builder, m, exponent)
	field := builder.String()
	if !state.Flag('-') && len(field) < width {
		fmt.Fprint(state, strings.Repeat(" ", width-len(field)))
	}
	fmt.Fprint(state, field)
	if state.Flag('-') && len(field) < width {
		fmt.Fprint(state, strings.Repeat(" ", width-len(field)))
	}
}

func (f formatSpec) PrintNumber(w io.Writer, m Mantissa, exponent int) {
	if f.sci {
		sep := "e"
		if f.capital {
			sep = "E"
		}
		f.printSci(w, m, exponent, sep)
	} else {
		f.printFixed(w, m, exponent)
	}
}

func (f formatSpec) printFixed(w io.Writer, m Mantissa, exponent int) {
	formatter := newFormatter(w, f.sigDigits, exponent, f.exactDigitCount)
	m.send(formatter)
	formatter.Finish()
}

func (f formatSpec) printSci(
	w io.Writer, m Mantissa, exponent int, sep string) {
	f.printFixed(w, m, 0)
	fmt.Fprint(w, sep)
	fmt.Fprintf(w, "%+03d", exponent)
}

type optionFunc func(p *printerSettings)

func (o optionFunc) mutate(p *printerSettings) {
	o(p)
}

func bigExponent(exponent int) bool {
	return exponent < -3 || exponent > 6
}

func mutateSettings(
	options []Option, settings *printerSettings) *printerSettings {
	for _, option := range options {
		option.mutate(settings)
	}
	return settings
}

type part interface {
	positDigitIter() func() positDigit
	limit() int
}

type lazyPart struct {
	sequence  Sequence
	positions *Positions
}

func newPart(sequence Sequence, positions *Positions) part {
	return &lazyPart{sequence: sequence, positions: positions}
}

func (p *lazyPart) limit() int {
	return p.positions.limit
}

func (p *lazyPart) positDigitIter() func() positDigit {
	filter := p.positions.filter()
	iter := p.sequence.positDigitIter()
	pd := iter()
	return func() positDigit {
		result := invalidPositDigit
		for pd.Valid() && pd.Posit < p.positions.limit && !result.Valid() {
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

func asDigits(s Sequence) Digits {
	consumer := new(digitAt)
	sendPositDigits(s, consumer)
	return Digits{digits: consumer.digits, posits: consumer.posits}
}

func find(s Sequence, pattern []int) func() int {
	if len(pattern) == 0 {
		return zeroPattern(s.positDigitIter())
	}
	return kmp(s.positDigitIter(), pattern)
}

type positionsFilter struct {
	starts     []int
	ranges     map[int]int
	startIndex int
	limit      int
}

func (p *positionsFilter) Includes(posit int) bool {
	p.update(posit)
	return posit < p.limit
}

func (p *positionsFilter) update(posit int) {
	for p.startIndex < len(p.starts) && p.starts[p.startIndex] <= posit {
		start := p.starts[p.startIndex]
		localLimit := start + p.ranges[start]
		if localLimit > p.limit {
			p.limit = localLimit
		}
		p.startIndex++
	}
}

// Package sqroot calculates square roots to arbitrary precision.
package sqroot

import (
	"fmt"
	"io"
	"math/big"
	"os"
	"strings"

	"github.com/keep94/consume2"
)

const (
	fPrecision = 6
	gPrecision = 16
)

var (
	zeroMantissa = &Mantissa{}
	zeroNumber   = &Number{}
)

// Mantissa represents the mantissa of a square root. Non zero Mantissas are
// between 0.1 inclusive and 1.0 exclusive. The number of digits of a
// Mantissa can be infinite. The zero value for a Mantissa corresponds to 0.
// By default, a Mantissa instance computes its digits lazily on demand each
// time. Computing the first N digits of a Mantissa takes O(N^2) time. However
// a Mantissa can be set to memoize its digits. Mantissa pointers implement
// Sequence. Mantissa instances do not support assignment. Mantissa instances
// are safe to use with multiple goroutines.
type Mantissa struct {
	spec mantissaSpec
}

// WithSignificant returns a Mantissa like this one that has no more than
// limit significant digits. WithSignificant rounds the returned Mantissa
// down toward zero when necessary. WithSignificant panics if limit is
// negative. If m memoizes its digts, then the returned Mantissa will also
// memoize its digits. Moreover, m and the returned Mantissa will share the
// same memoization data. WithSignificant will return m, if it can determine
// that m already has limit or fewer significant digits.
func (m *Mantissa) WithSignificant(limit int) *Mantissa {
	return m.withSpec(withLimit(m.spec, limit))
}

// WithMemoize returns a Mantissa like this one that remembers all of its
// previously computed digits. WithMemoize will return m, if memoization is
// already enabled for m.
func (m *Mantissa) WithMemoize() *Mantissa {
	return m.withSpec(withMemoize(m.spec))
}

// Digits returns all the digits of this Mantissa. If this Mantissa has an
// infinite number of digits, Digits runs forever.
func (m *Mantissa) Digits() Digits {
	var builder digitsBuilder
	consume2.FromGenerator[Digit](m.digitIter(), &builder)
	return builder.Build()
}

// Format prints this Mantissa with the f, F, g, G, e, E verbs. The verbs work
// in the usual way except that they always round down. Because Mantissas can
// have an infinite number of digits, g with no precision shows a max of 16
// significant digits. Format supports width, precision, and the '-' flag
// for left justification. The v verb is an alias for g.
func (m *Mantissa) Format(state fmt.State, verb rune) {
	formatSpec, ok := newFormatSpec(state, verb, 0)
	if !ok {
		fmt.Fprintf(state, "%%!%c(mantissa=%s)", verb, m.String())
		return
	}
	formatSpec.PrintField(state, m, 0)
}

// String returns the decimal representation of m as generated by m.Sprint(16).
func (m *Mantissa) String() string {
	return m.Sprint(gPrecision)
}

// Iterator returns the digits of this Mantissa as a function. The
// first call to returned function returns the first digit of Mantissa;
// the second call returns the second digit and so forth. If returned
// function runs out of Mantissa digits, it returns -1. If this
// Mantissa is zero, the returned function always returns -1.
func (m *Mantissa) Iterator() func() int {
	if m.spec == nil {
		return func() int { return -1 }
	}
	return m.spec.Iterator()
}

// Print works like Fprint and prints this Mantissa to stdout.
func (m *Mantissa) Print(maxDigits int, options ...Option) (n int, err error) {
	return m.Fprint(os.Stdout, maxDigits, options...)
}

// Sprint works like Fprint and prints this Mantissa to a string.
func (m *Mantissa) Sprint(maxDigits int, options ...Option) string {
	var builder strings.Builder
	m.Fprint(&builder, maxDigits, options...)
	return builder.String()
}

// Fprint prints this Mantissa to w. Fprint returns the number of bytes
// written and any error encountered. For options, the default is no
// separate rows, no separate columns, and digit count turned off.
func (m *Mantissa) Fprint(w io.Writer, maxDigits int, options ...Option) (
	n int, err error) {
	if m.IsZero() || maxDigits <= 0 {
		return fmt.Fprint(w, "0")
	}
	settings := &printerSettings{missingDigit: '.'}
	return fprint(
		w,
		newPart(m, new(PositionsBuilder).AddRange(0, maxDigits).Build()),
		mutateSettings(options, settings))
}

// At returns the digit at the given 0 based position in this Mantissa. If
// this Mantissa has posit or fewer digits, At returns -1. If posit is
// negative, At returns -1. By default, At has to compute all prior digits,
// so computing the kth digit takes O(k^2) time best case. However with
// memoization enabled, computing the kth digit takes O(1) time best case, but
// memoization stores all computed digits in memory. GetDigits() is a good
// alternative when only a few digits need to be computed because it stores
// only the needed digits in memory while iterating through the digits of
// the mantissa one time. With GetDigits(), computing k digits always takes
// O(N^2) time where N is the largest digit position of the batch of digits
// to be computed.
func (m *Mantissa) At(posit int) int {
	if m.spec == nil {
		return -1
	}
	return m.spec.At(posit)
}

// IsMemoize returns true if this Mantissa memoizes its digits. If this
// Mantissa is zero, IsMemoize always returns true.
func (m *Mantissa) IsMemoize() bool {
	if m.spec == nil {
		return true
	}
	return m.spec.IsMemoize()
}

// IsZero returns true if this Mantissa is zero.
func (m *Mantissa) IsZero() bool {
	return m.spec == nil
}

func (m *Mantissa) withSpec(newSpec mantissaSpec) *Mantissa {
	if newSpec == m.spec {
		return m
	}
	if newSpec == nil {
		return zeroMantissa
	}
	return &Mantissa{spec: newSpec}
}

func (m *Mantissa) digitIter() func() (Digit, bool) {
	iter := m.Iterator()
	digit := iter()
	index := 0
	return func() (dt Digit, ok bool) {
		if digit == -1 {
			return
		}
		result := Digit{Position: index, Value: digit}
		digit = iter()
		index++
		return result, true
	}
}

// Number represents a square root value. The zero value for Number
// corresponds to 0. A Number is of the form mantissa * 10^exponent where
// mantissa is between 0.1 inclusive and 1.0 exclusive. Like Mantissa, a
// Number instance can represent an infinite number of digits. Number
// instances do not support assignment. Number instances are safe to use with
// multiple goroutines.
type Number struct {
	mantissa *Mantissa
	exponent int
}

// WithSignificant returns a Number like this one that has no more than
// limit significant digits. WithSignificant rounds the returned Number
// down toward zero when necessary. WithSignificant panics if limit is
// negative. If the mantissa of n memoizes its digits, then the mantissa of
// the returned Number will also memoize its digits. Moreover, the two will
// share the same memoization data. WithSignificant will return n, if it can
// determine that n already has limit or fewer significant digits.
func (n *Number) WithSignificant(limit int) *Number {
	return n.withMantissa(n.Mantissa().WithSignificant(limit))
}

// WithMemoize returns a Number like this one that has a Mantissa that
// remembers all of its previously computed digits. WithMemoize returns n, if
// the mantissa of n already memoizes its digits.
func (n *Number) WithMemoize() *Number {
	return n.withMantissa(n.Mantissa().WithMemoize())
}

// Mantissa returns the pointer to the Mantissa of this Number.
func (n *Number) Mantissa() *Mantissa {
	if n.mantissa == nil {
		return zeroMantissa
	}
	return n.mantissa
}

// Exponent returns the exponent of this Number.
func (n *Number) Exponent() int {
	return n.exponent
}

// Format prints this Number with the f, F, g, G, e, E verbs. The verbs work
// in the usual way except that they always round down. Because Number can
// have an infinite number of digits, g with no precision shows a max of 16
// significant digits. Format supports width, precision, and the '-' flag
// for left justification. The v verb is an alias for g.
func (n *Number) Format(state fmt.State, verb rune) {
	formatSpec, ok := newFormatSpec(state, verb, n.exponent)
	if !ok {
		fmt.Fprintf(state, "%%!%c(number=%s)", verb, n.String())
		return
	}
	formatSpec.PrintField(state, n.Mantissa(), n.exponent)
}

// String returns the decimal representation of n using %g.
func (n *Number) String() string {
	var builder strings.Builder
	fs := formatSpec{sigDigits: gPrecision, sci: bigExponent(n.exponent)}
	fs.PrintNumber(&builder, n.Mantissa(), n.exponent)
	return builder.String()
}

// IsZero returns true if this Number is zero.
func (n *Number) IsZero() bool {
	return n.Mantissa().IsZero()
}

func (n *Number) withMantissa(newMantissa *Mantissa) *Number {
	if newMantissa == n.Mantissa() {
		return n
	}
	if newMantissa.IsZero() {
		return zeroNumber
	}
	return &Number{mantissa: newMantissa, exponent: n.exponent}
}

// Sqrt returns the square root of radican. Sqrt panics if radican is
// negative.
func Sqrt(radican int64) *Number {
	return sqrtFrac(big.NewInt(radican), one)
}

// SqrtRat returns the square root of num / denom. denom must be positive,
// and num must be non-negative or else SqrtRat panics.
func SqrtRat(num, denom int64) *Number {
	return sqrtFrac(big.NewInt(num), big.NewInt(denom))
}

// SqrtBigInt returns the square root of radican. SqrtBigInt panics if
// radican is negative.
func SqrtBigInt(radican *big.Int) *Number {
	return sqrtFrac(radican, one)
}

// SqrtBigRat returns the square root of radican. The denominator of radican
// must be positive, and the numerator must be non-negative or else SqrtBigRat
// panics.
func SqrtBigRat(radican *big.Rat) *Number {
	return sqrtFrac(radican.Num(), radican.Denom())
}

func sqrtFrac(num, denom *big.Int) *Number {
	num = new(big.Int).Set(num)
	denom = new(big.Int).Set(denom)
	if denom.Sign() <= 0 {
		panic("Denominator must be positive")
	}
	if num.Sign() < 0 {
		panic("Numerator must be non-negative")
	}
	if num.Sign() == 0 {
		return zeroNumber
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
	return &Number{exponent: exp, mantissa: &Mantissa{spec: spec}}
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

func (f formatSpec) PrintField(state fmt.State, m *Mantissa, exponent int) {
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

func (f formatSpec) PrintNumber(w io.Writer, m *Mantissa, exponent int) {
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

func (f formatSpec) printFixed(w io.Writer, m *Mantissa, exponent int) {
	formatter := newFormatter(w, f.sigDigits, exponent, f.exactDigitCount)
	consumer := consume2.Map[Digit, int](
		formatter, func(d Digit) int { return d.Value })
	consume2.FromGenerator(m.digitIter(), consumer)
	formatter.Finish()
}

func (f formatSpec) printSci(
	w io.Writer, m *Mantissa, exponent int, sep string) {
	f.printFixed(w, m, 0)
	fmt.Fprint(w, sep)
	fmt.Fprintf(w, "%+03d", exponent)
}

func bigExponent(exponent int) bool {
	return exponent < -3 || exponent > 6
}

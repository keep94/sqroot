// Package sqroot computes square roots and cube roots to arbitrary precision.
//
// Number is the main type in this package. It represents a lazily evaluated
// non negative real number that generally has an infinite number of digits,
// but a Number can also have a finite number of digits.
//
// A FiniteNumber works like Number except that it always has a finite
// number of digits. A *FiniteNumber can be used anywhere a Number type
// is expected but not the other way around.
//
// A Sequence is a view of a contiguous subset of digits of a Number.
// For example, A Sequence could represent everything past the 1000th digit
// of the square root of 3. Because Sequences are views, they are cheap to
// create. Note that Number and *FiniteNumber can be used anywhere a Sequence
// type is expected. A Sequence can be either infinite or finite in length.
//
// A FiniteSequence works like Sequence except unlike Sequence, a
// FiniteSequence is always finite in length. A FiniteSequence can be used
// anywhere a Sequence is expected, and a *FiniteNumber can be used anywhere
// a FiniteSequence is expected. However, a Number or Sequence cannot be
// used where a FiniteSequence is expected because they can have an infinite
// number of digits. A FiniteSequence must have a finite number of digits.
package sqroot

import (
	"errors"
	"fmt"
	"io"
	"math"
	"math/big"
	"strings"

	"github.com/keep94/consume2"
)

const (
	fPrecision = 6
	gPrecision = 16
)

var (
	zeroNumber = &FiniteNumber{}
)

// Number is a reference to a non-negative real number.
// A non-zero Number is of the form mantissa * 10^exponent
// where mantissa is between 0.1 inclusive and 1.0 exclusive. A Number
// can represent either a finite or infinite number of digits. A Number
// computes the digits of its mantissa lazily on an as needed basis. To
// compute a given digit, a Number must compute all digits that come before
// that digit. A Number stores its computed digits so that they only have
// to be computed once. Number instances are safe to use with multiple
// goroutines.
//
// The Number factory functions such as the Sqrt and CubeRoot functions
// return new Number instances that contain no computed digits initially.
//
// Because Number instances store their computed digits, it is best to
// reuse a Number instance when possible. For example the code:
//
//	n := sqroot.Sqrt(6)
//	fmt.Println(sqroot.FindFirst(n, []int{0, 0, 0, 0, 0}))
//	fmt.Println(sqroot.FindFirst(n, []int{2, 2, 2, 2, 2}))
//
// runs faster than the code:
//
//	fmt.Println(sqroot.FindFirst(sqroot.Sqrt(6), []int{0, 0, 0, 0, 0}))
//	fmt.Println(sqroot.FindFirst(sqroot.Sqrt(6), []int{2, 2, 2, 2, 2}))
//
// In the first code block, the second line reuses digits computed in the
// first line, but in the second code block, no reuse is possible since
// sqroot.Sqrt(6) always returns a Number with no precomputed digits.
//
// A Number can be 0, in which case IsZero() returns true. Zero Numbers
// have an exponent of 0 and no digits in their mantissa. This means that
// calling At() on a zero Number always returns -1. Likewise calling
// Iterator() or Reverse() on a zero Number returns a function that yields
// no digits. However, calling String() on a zero Number returns "0" and
// printing a zero Number prints 0 according to the format specification
// used.
type Number interface {
	Sequence

	// At returns the significant digit of this Number at the given 0 based
	// position. If this Number has posit or fewer significant digits, At
	// returns -1. If posit is negative, At returns -1.
	At(posit int) int

	// WithSignificant returns a view of this Number that has no more than
	// limit significant digits. WithSignificant rounds the returned value
	// down toward zero. WithSignificant panics if limit is negative.
	WithSignificant(limit int) *FiniteNumber

	// Exponent returns the exponent of this Number.
	Exponent() int

	// Format prints this Number with the f, F, g, G, e, E verbs. The
	// verbs work in the usual way except that they always round down.
	// Because Number can have an infinite number of digits, g with no
	// precision shows a max of 16 significant digits. Format supports
	// width, precision, and the '-' flag for left justification. The v
	// verb is an alias for g.
	Format(state fmt.State, verb rune)

	// String returns the decimal representation of this Number using %g.
	String() string

	// IsZero returns true if this Number is zero.
	IsZero() bool

	withExponent(e int) Number
}

// Sqrt returns the square root of radican. Sqrt panics if radican is
// negative.
func Sqrt(radican int64) Number {
	return nRootFrac(big.NewInt(radican), one, newSqrtManager)
}

// SqrtRat returns the square root of num / denom. denom must be positive,
// and num must be non-negative or else SqrtRat panics.
func SqrtRat(num, denom int64) Number {
	return nRootFrac(big.NewInt(num), big.NewInt(denom), newSqrtManager)
}

// SqrtBigInt returns the square root of radican. SqrtBigInt panics if
// radican is negative.
func SqrtBigInt(radican *big.Int) Number {
	return nRootFrac(radican, one, newSqrtManager)
}

// SqrtBigRat returns the square root of radican. The denominator of radican
// must be positive, and the numerator must be non-negative or else SqrtBigRat
// panics.
func SqrtBigRat(radican *big.Rat) Number {
	return nRootFrac(radican.Num(), radican.Denom(), newSqrtManager)
}

// CubeRoot returns the cube root of radican. CubeRoot panics if radican is
// negative as Number can only hold positive results.
func CubeRoot(radican int64) Number {
	return nRootFrac(big.NewInt(radican), one, newCubeRootManager)
}

// CubeRootRat returns the cube root of num / denom. Because Number can only
// hold positive results, denom must be positive, and num must be non-negative
// or else CubeRootRat panics.
func CubeRootRat(num, denom int64) Number {
	return nRootFrac(big.NewInt(num), big.NewInt(denom), newCubeRootManager)
}

// CubeRootBigInt returns the cube root of radican. CubeRootBigInt panics if
// radican is negative as Number can only hold positive results.
func CubeRootBigInt(radican *big.Int) Number {
	return nRootFrac(radican, one, newCubeRootManager)
}

// CubeRootBigRat returns the cube root of radican. Because Number can only
// hold positive results, the denominator of radican must be positive, and the
// numerator must be non-negative or else CubeRootBigRat panics.
func CubeRootBigRat(radican *big.Rat) Number {
	return nRootFrac(radican.Num(), radican.Denom(), newCubeRootManager)
}

// NewNumberFromBigRat returns value as a Number. Because Number can only
// hold positive results, the denominator of value must be positive, and the
// numerator must be non-negative or else NewNumberFromBigRat panics.
// NewNumberFromBigRat can be used to create arbitrary Number instances for
// testing.
func NewNumberFromBigRat(value *big.Rat) Number {
	num := value.Num()
	denom := value.Denom()
	checkNumDenom(num, denom)
	if num.Sign() == 0 {
		return zeroNumber
	}
	return newNumber(newRatGenerator(num, denom))
}

// NewNumberForTesting creates an arbitrary Number for testing. fixed are
// digits between 0 and 9 representing the non repeating digits that come
// immediately after the decimal place of the mantissa. repeating are digits
// between 0 and 9 representing the repeating digits that follow the non
// repeating digits of the mantissa. exp is the exponent part of the
// returned Number. NewNumberForTesting returns an error if fixed or
// repeating contain values not between 0 and 9, or if the first digit of
// the mantissa would be zero since mantissas must be between 0.1 inclusive
// and 1.0 exclusive.
func NewNumberForTesting(fixed, repeating []int, exp int) (Number, error) {
	if len(fixed) == 0 && len(repeating) == 0 {
		return zeroNumber, nil
	}
	if !validDigits(fixed) || !validDigits(repeating) {
		return nil, errors.New("NewNumberForTesting: digits must be between 0 and 9")
	}
	gen := newRepeatingGenerator(fixed, repeating, exp)
	digits, _ := gen.Generate()
	if digits() == 0 {
		return nil, errors.New("NewNumberForTesting: leading zeros not allowed in digits")
	}
	return newNumber(gen), nil
}

// NewNumber returns a new Number based on g. Although g is expected to
// follow the contract of Generator, if g yields mantissa digits outside the
// range of 0 and 9, NewNumber regards that as a signal that there are no
// more mantissa digits. Also if g happens to yield 0 as the first digit
// of the mantissa, NewNumber will return zero.
func NewNumber(g Generator) Number {

	// gen is guaranteed to follow the Generator contract. if g yields
	// a digit outside the range of 0 and 9, gen will signal no more digits
	// in the mantissa. However, we still have to check that the first
	// mantissa digit yielded is not zero.
	gen := newValidDigits(g)

	digits, _ := gen.Generate()
	first := digits()
	if first == 0 || first == -1 {
		return zeroNumber
	}
	return newNumber(gen)
}

// FiniteNumber is a Number with a finite number of digits. FiniteNumber
// implements both Number and FiniteSequence. The zero value for FiniteNumber
// is 0.
//
// Pass FiniteNumber instances by reference not by value. Copying a
// FiniteNumber instance or overwriting a FiniteNumber instance with the
// assignment operator is not supported and may cause errors.
type FiniteNumber struct {
	mantissa mantissa
	exponent int
}

// WithStart comes from the Sequence interface.
func (n *FiniteNumber) WithStart(start int) Sequence {
	return n.FiniteWithStart(start)
}

// FiniteWithStart comes from the FiniteSequence interface.
func (n *FiniteNumber) FiniteWithStart(start int) FiniteSequence {
	if start <= 0 {
		return n
	}
	return &mantissaWithStart{
		mantissa: n.mantissa,
		start:    start,
	}
}

// WithEnd comes from the Sequence interface.
func (n *FiniteNumber) WithEnd(end int) FiniteSequence {
	return n.withMantissa(n.mantissa.WithLimit(end))
}

// At comes from the Number interface.
func (n *FiniteNumber) At(posit int) int {
	return n.mantissa.At(posit)
}

// WithSignificant comes from the Number interface.
func (n *FiniteNumber) WithSignificant(limit int) *FiniteNumber {
	if limit < 0 {
		panic("limit must be non-negative")
	}
	return n.withMantissa(n.mantissa.WithLimit(limit))
}

// Exponent comes from the Number interface.
func (n *FiniteNumber) Exponent() int {
	return n.exponent
}

// Format comes from the Number interface.
func (n *FiniteNumber) Format(state fmt.State, verb rune) {
	formatSpec, ok := newFormatSpec(state, verb, n.exponent)
	if !ok {
		fmt.Fprintf(state, "%%!%c(number=%s)", verb, n.String())
		return
	}
	formatSpec.PrintField(state, n)
}

// Exact works like String, but uses enough significant digits to return
// the exact representation of n.
func (n *FiniteNumber) Exact() string {
	var builder strings.Builder
	fs := formatSpecForG(endOf(n), n.exponent, false)
	fs.PrintNumber(&builder, n)
	return builder.String()
}

// String comes from the Number interface.
func (n *FiniteNumber) String() string {
	var builder strings.Builder
	fs := formatSpecForG(gPrecision, n.exponent, false)
	fs.PrintNumber(&builder, n)
	return builder.String()
}

// IsZero comes from the Number interface.
func (n *FiniteNumber) IsZero() bool {
	return n.mantissa.IsZero()
}

// Iterator comes from the Sequence interface.
func (n *FiniteNumber) Iterator() func() (Digit, bool) {
	return n.mantissa.IteratorAt(0)
}

// Reverse comes from the FiniteSequence interface.
func (n *FiniteNumber) Reverse() func() (Digit, bool) {
	return n.mantissa.ReverseTo(0)
}

func (n *FiniteNumber) withExponent(e int) Number {
	if e == n.exponent || n.IsZero() {
		return n
	}
	return &FiniteNumber{exponent: e, mantissa: n.mantissa}
}

func (n *FiniteNumber) withMantissa(newMantissa mantissa) *FiniteNumber {
	if newMantissa == n.mantissa {
		return n
	}
	if newMantissa.IsZero() {
		return zeroNumber
	}
	return &FiniteNumber{mantissa: newMantissa, exponent: n.exponent}
}

func (n *FiniteNumber) private() {
}

func nRootFrac(
	num, denom *big.Int, newManager func() rootManager) Number {
	checkNumDenom(num, denom)
	if num.Sign() == 0 {
		return zeroNumber
	}
	return newNumber(newNRootGenerator(num, denom, newManager))
}

// newNumber returns a new number based on gen. Unlike NewNumber, gen must
// follow the contract of Generator. Also, newNumber doesn't handle empty
// mantissas.
func newNumber(gen Generator) Number {
	digits, exp := gen.Generate()
	mantissa := mantissa{spec: newMemoizeSpec(digits)}
	f := &FiniteNumber{exponent: exp, mantissa: mantissa}
	return opaqueNumber(f)
}

func checkNumDenom(num, denom *big.Int) {
	if denom.Sign() <= 0 {
		panic("Denominator must be positive")
	}
	if num.Sign() < 0 {
		panic("Numerator must be non-negative")
	}
}

type mantissa struct {
	spec numberSpec
}

func (m mantissa) At(posit int) int {
	if m.spec == nil {
		return -1
	}
	return m.spec.At(posit)
}

func (m mantissa) IsZero() bool {
	return m.spec == nil
}

func (m mantissa) ReverseTo(start int) func() (Digit, bool) {
	digits := m.allDigits()
	index := len(digits)
	return func() (d Digit, ok bool) {
		if index <= start {
			return
		}
		index--
		return Digit{Position: index, Value: int(digits[index])}, true
	}
}

func (m mantissa) IteratorAt(index int) func() (Digit, bool) {
	if m.spec == nil {
		return func() (Digit, bool) { return Digit{}, false }
	}
	return m.spec.IteratorAt(index, math.MaxInt)
}

func (m mantissa) WithLimit(limit int) mantissa {
	return mantissa{spec: withLimit(m.spec, limit)}
}

func (m mantissa) allDigits() []int8 {
	if m.spec == nil {
		return nil
	}
	return m.spec.FirstN(math.MaxInt)
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
	switch verb {
	case 'f', 'F':
		if !precisionOk {
			precision = fPrecision
		}
		return formatSpecForF(precision, exponent), true
	case 'g', 'G', 'v':
		if !precisionOk {
			precision = gPrecision
		}
		return formatSpecForG(precision, exponent, verb == 'G'), true
	case 'e', 'E':
		if !precisionOk {
			precision = fPrecision
		}
		return formatSpecForE(precision, verb == 'E'), true
	default:
		return formatSpec{}, false
	}
}

func formatSpecForF(precision, exponent int) formatSpec {
	sigDigits := precision + exponent
	return formatSpec{sigDigits: sigDigits, exactDigitCount: true}
}

func formatSpecForG(precision, exponent int, capital bool) formatSpec {
	sigDigits := precision
	if sigDigits == 0 {
		sigDigits = 1
	}
	sci := sigDigits < exponent || bigExponent(exponent)
	return formatSpec{sigDigits: sigDigits, sci: sci, capital: capital}
}

func formatSpecForE(precision int, capital bool) formatSpec {
	return formatSpec{
		sigDigits:       precision,
		exactDigitCount: true,
		sci:             true,
		capital:         capital}
}

func (f formatSpec) PrintField(state fmt.State, n *FiniteNumber) {
	width, widthOk := state.Width()
	if !widthOk {
		f.PrintNumber(state, n)
		return
	}
	var builder strings.Builder
	f.PrintNumber(&builder, n)
	field := builder.String()
	if !state.Flag('-') && len(field) < width {
		fmt.Fprint(state, strings.Repeat(" ", width-len(field)))
	}
	fmt.Fprint(state, field)
	if state.Flag('-') && len(field) < width {
		fmt.Fprint(state, strings.Repeat(" ", width-len(field)))
	}
}

func (f formatSpec) PrintNumber(w io.Writer, n *FiniteNumber) {
	if f.sci {
		sep := "e"
		if f.capital {
			sep = "E"
		}
		f.printSci(w, n.mantissa, n.exponent, sep)
	} else {
		f.printFixed(w, n.mantissa, n.exponent)
	}
}

func (f formatSpec) printFixed(w io.Writer, m mantissa, exponent int) {
	formatter := newFormatter(w, f.sigDigits, exponent, f.exactDigitCount)
	consume2.FromGenerator[Digit](m.IteratorAt(0), formatter)
	formatter.Finish()
}

func (f formatSpec) printSci(
	w io.Writer, m mantissa, exponent int, sep string) {
	f.printFixed(w, m, 0)
	fmt.Fprint(w, sep)
	fmt.Fprintf(w, "%+03d", exponent)
}

func bigExponent(exponent int) bool {
	return exponent < -3 || exponent > 6
}

type mantissaWithStart struct {
	mantissa mantissa
	start    int
}

func (m *mantissaWithStart) Iterator() func() (Digit, bool) {
	return m.mantissa.IteratorAt(m.start)
}

func (m *mantissaWithStart) Reverse() func() (Digit, bool) {
	return m.mantissa.ReverseTo(m.start)
}

func (m *mantissaWithStart) WithStart(start int) Sequence {
	return m.FiniteWithStart(start)
}

func (m *mantissaWithStart) FiniteWithStart(start int) FiniteSequence {
	if start <= m.start {
		return m
	}
	return &mantissaWithStart{mantissa: m.mantissa, start: start}
}

func (m *mantissaWithStart) WithEnd(end int) FiniteSequence {
	return m.withMantissa(m.mantissa.WithLimit(end))
}

func (m *mantissaWithStart) withMantissa(mantissa mantissa) *mantissaWithStart {
	if mantissa == m.mantissa {
		return m
	}
	return &mantissaWithStart{mantissa: mantissa, start: m.start}
}

func (m *mantissaWithStart) private() {
}

func opaqueNumber(n Number) Number {
	if _, ok := n.(*opqNumber); ok {
		return n
	}
	return &opqNumber{Number: n}
}

type opqNumber struct {
	Number
}

func (n *opqNumber) WithStart(start int) Sequence {
	result := n.Number.WithStart(start)
	if result == n.Number {
		return n
	}
	return opaqueSequence(result)
}

func (n *opqNumber) withExponent(e int) Number {
	result := n.Number.withExponent(e)
	if result == n.Number {
		return n
	}
	return opaqueNumber(result)
}

// Package sqroot computes square roots and cube roots to arbitrary precision.
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

// Number is a reference to a square root value.
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

// FiniteNumber is a Number with a finite number of digits. FiniteNumber
// implements both Number and FiniteSequence. The zero value for FiniteNumber
// is 0.
//
// Pass FiniteNumber instances by reference not by value. Copying a
// FiniteNumber instance or overwriting a FiniteNumber instance with the
// assignment operator is not supported and may cause errors.
type FiniteNumber struct {
	spec     numberSpec
	exponent int
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

// WithStart comes from the Sequence interface.
func (n *FiniteNumber) WithStart(start int) Sequence {
	return n.FiniteWithStart(start)
}

// FiniteWithStart comes from the FiniteSequence interface.
func (n *FiniteNumber) FiniteWithStart(start int) FiniteSequence {
	if start <= 0 {
		return n
	}
	return &numberWithStart{
		number: n,
		start:  start,
	}
}

// WithEnd comes from the Sequence interface.
func (n *FiniteNumber) WithEnd(end int) FiniteSequence {
	return n.withSignificant(end)
}

// At comes from the Number interface.
func (n *FiniteNumber) At(posit int) int {
	if n.spec == nil {
		return -1
	}
	return n.spec.At(posit)
}

// WithSignificant comes from the Number interface.
func (n *FiniteNumber) WithSignificant(limit int) *FiniteNumber {
	if limit < 0 {
		panic("limit must be non-negative")
	}
	return n.withSignificant(limit)
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
	return n.spec == nil
}

// Iterator comes from the Sequence interface.
func (n *FiniteNumber) Iterator() func() (Digit, bool) {
	return n.fullIteratorAt(0)
}

// Reverse comes from the FiniteSequence interface.
func (n *FiniteNumber) Reverse() func() (Digit, bool) {
	return n.fullReverseTo(0)
}

func (n *FiniteNumber) withExponent(e int) Number {
	if e == n.exponent || n.IsZero() {
		return n
	}
	return &FiniteNumber{exponent: e, spec: n.spec}
}

func (n *FiniteNumber) fullIteratorAt(index int) func() (Digit, bool) {
	iter := n.iteratorAt(index)
	dig := iter()
	return func() (dt Digit, ok bool) {
		if dig == -1 {
			return
		}
		result := Digit{Position: index, Value: dig}
		dig = iter()
		index++
		return result, true
	}
}

func (n *FiniteNumber) fullReverseTo(start int) func() (Digit, bool) {
	digits := n.allDigits()
	index := len(digits)
	return func() (d Digit, ok bool) {
		if index <= start {
			return
		}
		index--
		return Digit{Position: index, Value: int(digits[index])}, true
	}
}

func (n *FiniteNumber) iteratorAt(index int) func() int {
	if n.spec == nil {
		return func() int { return -1 }
	}
	return n.spec.IteratorAt(index)
}

func (n *FiniteNumber) allDigits() []int8 {
	if n.spec == nil {
		return nil
	}
	return n.spec.FirstN(math.MaxInt)
}

func (n *FiniteNumber) withSpec(newSpec numberSpec) *FiniteNumber {
	if newSpec == n.spec {
		return n
	}
	if newSpec == nil {
		return zeroNumber
	}
	return &FiniteNumber{spec: newSpec, exponent: n.exponent}
}

func (n *FiniteNumber) withSignificant(limit int) *FiniteNumber {
	return n.withSpec(withLimit(n.spec, limit))
}

func (n *FiniteNumber) private() {
}

func nRootFrac(
	num, denom *big.Int, newManager func() rootManager) *FiniteNumber {
	checkNumDenom(num, denom)
	if num.Sign() == 0 {
		return zeroNumber
	}
	return newNumber(newNRootGenerator(num, denom, newManager))
}

// newNumber returns a new number based on gen. Unlike NewNumber, gen must
// follow the contract of Generator. Also, newNumber doesn't handle empty
// mantissas.
func newNumber(gen Generator) *FiniteNumber {
	digits, exp := gen.Generate()
	return &FiniteNumber{exponent: exp, spec: newMemoizeSpec(digits)}
}

func checkNumDenom(num, denom *big.Int) {
	if denom.Sign() <= 0 {
		panic("Denominator must be positive")
	}
	if num.Sign() < 0 {
		panic("Numerator must be non-negative")
	}
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
		f.printSci(w, n, n.exponent, sep)
	} else {
		f.printFixed(w, n, n.exponent)
	}
}

func (f formatSpec) printFixed(w io.Writer, n *FiniteNumber, exponent int) {
	formatter := newFormatter(w, f.sigDigits, exponent, f.exactDigitCount)
	consume2.FromIntGenerator(n.iteratorAt(0), formatter)
	formatter.Finish()
}

func (f formatSpec) printSci(
	w io.Writer, n *FiniteNumber, exponent int, sep string) {
	f.printFixed(w, n, 0)
	fmt.Fprint(w, sep)
	fmt.Fprintf(w, "%+03d", exponent)
}

func bigExponent(exponent int) bool {
	return exponent < -3 || exponent > 6
}

type numberWithStart struct {
	number *FiniteNumber
	start  int
}

func (n *numberWithStart) Iterator() func() (Digit, bool) {
	return n.number.fullIteratorAt(n.start)
}

func (n *numberWithStart) Reverse() func() (Digit, bool) {
	return n.number.fullReverseTo(n.start)
}

func (n *numberWithStart) WithStart(start int) Sequence {
	return n.FiniteWithStart(start)
}

func (n *numberWithStart) FiniteWithStart(start int) FiniteSequence {
	if start <= n.start {
		return n
	}
	return &numberWithStart{number: n.number, start: start}
}

func (n *numberWithStart) WithEnd(end int) FiniteSequence {
	return n.withNumber(n.number.withSignificant(end))
}

func (n *numberWithStart) withNumber(
	number *FiniteNumber) *numberWithStart {
	if number == n.number {
		return n
	}
	return &numberWithStart{number: number, start: n.start}
}

func (n *numberWithStart) private() {
}

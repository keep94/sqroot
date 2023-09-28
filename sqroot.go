// Package sqroot computes square roots and cube roots to arbitrary precision.
package sqroot

import (
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
	zeroNumber = &Number{}
)

// Number represents a square root value. The zero value for Number
// corresponds to 0. A Number is of the form mantissa * 10^exponent where
// mantissa is between 0.1 inclusive and 1.0 exclusive. A Number instance
// can represent an infinite number of digits. Number pointers implement
// Sequence. Number instances are safe to use with multiple goroutines.
type Number struct {
	spec     numberSpec
	exponent int
}

// Sqrt returns the square root of radican. Sqrt panics if radican is
// negative.
func Sqrt(radican int64) *Number {
	return nRootFrac(big.NewInt(radican), one, newSqrtManager)
}

// SqrtRat returns the square root of num / denom. denom must be positive,
// and num must be non-negative or else SqrtRat panics.
func SqrtRat(num, denom int64) *Number {
	return nRootFrac(big.NewInt(num), big.NewInt(denom), newSqrtManager)
}

// SqrtBigInt returns the square root of radican. SqrtBigInt panics if
// radican is negative.
func SqrtBigInt(radican *big.Int) *Number {
	return nRootFrac(radican, one, newSqrtManager)
}

// SqrtBigRat returns the square root of radican. The denominator of radican
// must be positive, and the numerator must be non-negative or else SqrtBigRat
// panics.
func SqrtBigRat(radican *big.Rat) *Number {
	return nRootFrac(radican.Num(), radican.Denom(), newSqrtManager)
}

// CubeRoot returns the cube root of radican. CubeRoot panics if radican is
// negative as Number can only hold positive results.
func CubeRoot(radican int64) *Number {
	return nRootFrac(big.NewInt(radican), one, newCubeRootManager)
}

// CubeRootRat returns the cube root of num / denom. Because Number can only
// hold positive results, denom must be positive, and num must be non-negative
// or else CubeRootRat panics.
func CubeRootRat(num, denom int64) *Number {
	return nRootFrac(big.NewInt(num), big.NewInt(denom), newCubeRootManager)
}

// CubeRootBigInt returns the cube root of radican. CubeRootBigInt panics if
// radican is negative as Number can only hold positive results.
func CubeRootBigInt(radican *big.Int) *Number {
	return nRootFrac(radican, one, newCubeRootManager)
}

// CubeRootBigRat returns the cube root of radican. Because Number can only
// hold positive results, the denominator of radican must be positive, and the
// numerator must be non-negative or else CubeRootBigRat panics.
func CubeRootBigRat(radican *big.Rat) *Number {
	return nRootFrac(radican.Num(), radican.Denom(), newCubeRootManager)
}

// NewNumberFromBigRat returns value as a Number. Because Number can only
// hold positive results, the denominator of value must be positive, and the
// numerator must be non-negative or else NewNumberFromBigRat panics.
// NewNumberFromBigRat can be used to create arbitrary Number instances for
// testing.
func NewNumberFromBigRat(value *big.Rat) *Number {
	num := value.Num()
	denom := value.Denom()
	checkNumDenom(num, denom)
	if num.Sign() == 0 {
		return zeroNumber
	}
	groups, exp := computeGroupsFromRational(num, denom, ten)
	digits := groupsToDigits(groups)
	return &Number{exponent: exp, spec: newMemoizeSpec(digits)}
}

// WithStart comes from the Sequence interface.
func (n *Number) WithStart(start int) Sequence {
	if start <= 0 {
		return n
	}
	return &numberWithStart{
		number: n,
		start:  start,
	}
}

// WithEnd comes from the Sequence interface.
func (n *Number) WithEnd(end int) Sequence {
	return n.withSignificant(end)
}

// Iterator returns the digits of the mantissa of this Number as a function.
// The first call to returned function returns the first digit; the second
// call returns the second digit and so forth. If returned function runs out
// of digits, it returns -1. If this Number is zero, the returned function
// always returns -1.
func (n *Number) Iterator() func() int {
	return n.iteratorAt(0)
}

// IteratorAt works like Iterator except that it starts at the given 0-based
// position instead of at the beginning. In fact, calling IteratorAt(0) is
// the same as calling Iterator(). If posit is negative, IteratorAt panics.
func (n *Number) IteratorAt(posit int) func() int {
	if posit < 0 {
		panic("posit must be non-negative")
	}
	return n.iteratorAt(posit)
}

// At returns the significant digit of n at the given 0 based position.
// If n has posit or fewer significant digits, At returns -1. If posit is
// negative, At returns -1.
func (n *Number) At(posit int) int {
	if n.spec == nil {
		return -1
	}
	return n.spec.At(posit)
}

// WithSignificant returns a Number like this one that has no more than
// limit significant digits. WithSignificant rounds the returned Number
// down toward zero when necessary. WithSignificant panics if limit is
// negative. WithSignificant will return n, if it can determine that n
// already has limit or fewer significant digits.
func (n *Number) WithSignificant(limit int) *Number {
	if limit < 0 {
		panic("limit must be non-negative")
	}
	return n.withSignificant(limit)
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
	formatSpec.PrintField(state, n)
}

// String returns the decimal representation of n using %g.
func (n *Number) String() string {
	var builder strings.Builder
	fs := formatSpec{sigDigits: gPrecision, sci: bigExponent(n.exponent)}
	fs.PrintNumber(&builder, n)
	return builder.String()
}

// IsZero returns true if this Number is zero.
func (n *Number) IsZero() bool {
	return n.spec == nil
}

// NumDigits returns the total number of significant digits this Number has.
// If this Number has an infinite number of digits, NumDigits runs forever.
func (n *Number) NumDigits() int {
	return len(n.allDigits())
}

// Reverse returns a function that generates the significant digits of this
// Number in reverse order. The first call to the returned function returns
// the last digit; the second call returns the second to last digit and so
// forth. When there are no more digits, the returned function returns -1.
// If this Number has an infinite number of digits, Reverse runs forever.
func (n *Number) Reverse() func() int {
	digits := n.allDigits()
	index := len(digits)
	return func() int {
		if index <= 0 {
			return -1
		}
		index--
		return int(digits[index])
	}
}

// FullIterator comes from the Sequence interface.
func (n *Number) FullIterator() func() (Digit, bool) {
	return n.fullIteratorAt(0)
}

// FullReverse comes from the Sequence interface.
func (n *Number) FullReverse() func() (Digit, bool) {
	return n.fullReverseTo(0)
}

func (n *Number) withExponent(e int) *Number {
	if e == n.exponent || n.IsZero() {
		return n
	}
	return &Number{exponent: e, spec: n.spec}
}

func (n *Number) fullIteratorAt(index int) func() (Digit, bool) {
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

func (n *Number) fullReverseTo(start int) func() (Digit, bool) {
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

func (n *Number) iteratorAt(index int) func() int {
	if n.spec == nil {
		return func() int { return -1 }
	}
	return n.spec.IteratorAt(index)
}

func (n *Number) allDigits() []int8 {
	if n.spec == nil {
		return nil
	}
	return n.spec.FirstN(math.MaxInt)
}

func (n *Number) withSpec(newSpec numberSpec) *Number {
	if newSpec == n.spec {
		return n
	}
	if newSpec == nil {
		return zeroNumber
	}
	return &Number{spec: newSpec, exponent: n.exponent}
}

func (n *Number) withSignificant(limit int) *Number {
	return n.withSpec(withLimit(n.spec, limit))
}

func (n *Number) private() {
}

func nRootFrac(num, denom *big.Int, newManager func() rootManager) *Number {
	checkNumDenom(num, denom)
	if num.Sign() == 0 {
		return zeroNumber
	}
	manager := newManager()
	groups, exp := computeGroupsFromRational(
		num, denom, manager.Base(new(big.Int)))
	digits := computeRootDigits(groups, manager)
	return &Number{exponent: exp, spec: newMemoizeSpec(digits)}
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

func (f formatSpec) PrintField(state fmt.State, n *Number) {
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

func (f formatSpec) PrintNumber(w io.Writer, n *Number) {
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

func (f formatSpec) printFixed(w io.Writer, n *Number, exponent int) {
	formatter := newFormatter(w, f.sigDigits, exponent, f.exactDigitCount)
	consume2.FromIntGenerator(n.Iterator(), formatter)
	formatter.Finish()
}

func (f formatSpec) printSci(
	w io.Writer, n *Number, exponent int, sep string) {
	f.printFixed(w, n, 0)
	fmt.Fprint(w, sep)
	fmt.Fprintf(w, "%+03d", exponent)
}

func bigExponent(exponent int) bool {
	return exponent < -3 || exponent > 6
}

type numberWithStart struct {
	number *Number
	start  int
}

func (n *numberWithStart) FullIterator() func() (Digit, bool) {
	return n.number.fullIteratorAt(n.start)
}

func (n *numberWithStart) FullReverse() func() (Digit, bool) {
	return n.number.fullReverseTo(n.start)
}

func (n *numberWithStart) WithStart(start int) Sequence {
	if start <= n.start {
		return n
	}
	return &numberWithStart{number: n.number, start: start}
}

func (n *numberWithStart) WithEnd(end int) Sequence {
	return n.withNumber(n.number.withSignificant(end))
}

func (n *numberWithStart) withNumber(number *Number) Sequence {
	if number == n.number {
		return n
	}
	return &numberWithStart{number: number, start: n.start}
}

func (n *numberWithStart) private() {
}

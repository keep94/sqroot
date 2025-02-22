package sqroot

import (
	"math/big"
)

// Interface Generator lazily generates the digits of a Number.
type Generator interface {

	// Generate returns the digits of the mantissa and the exponent for a
	// Number. Numbers are of the form mantissa*10^exp where mantissa is
	// is between 0.1 inclusive and 1.0 exclusive. Calling digits() returns
	// each digit of the mantissa in turn. The first call to digits() cannot
	// return 0 because of the range of values a mantissa can have. If a
	// call to digits() returns -1, this means that there are no more digits
	// in the mantissa. If the first call to digits() return -1, that means
	// that the resulting Number is zero. Once a call to digits() returns -1,
	// all successive calls to digits() must also return -1. digits() must
	// return values between 0 and 9 or -1.
	Generate() (digits func() int, exp int)
}

func newNRootGenerator(
	num, denom *big.Int, newManager func() rootManager) Generator {
	result := &nrootGenerator{newManager: newManager}
	result.num.Set(num)
	result.denom.Set(denom)
	return result
}

// remove for v4
func newRatGenerator(num, denom *big.Int) Generator {
	var result ratGenerator
	result.num.Set(num)
	result.denom.Set(denom)
	return &result
}

func newRepeatingGenerator(fixed, repeating []int, exp int) Generator {
	var result repeatingGenerator
	result.fixed = append([]int(nil), fixed...)
	result.repeating = append([]int(nil), repeating...)
	result.exp = exp
	return &result
}

type repeatingGenerator struct {
	fixed     []int
	repeating []int
	exp       int
}

func (g *repeatingGenerator) Generate() (func() int, int) {
	fixedIndex := 0
	repeatingIndex := 0
	gen := func() int {
		if fixedIndex < len(g.fixed) {
			temp := g.fixed[fixedIndex]
			fixedIndex++
			return temp
		}
		if len(g.repeating) == 0 {
			return -1
		}
		temp := g.repeating[repeatingIndex]
		repeatingIndex = (repeatingIndex + 1) % len(g.repeating)
		return temp
	}
	return gen, g.exp
}

// remove for v4
type ratGenerator struct {
	num   big.Int
	denom big.Int
}

func (g *ratGenerator) Generate() (func() int, int) {
	groups, exp := computeGroupsFromRational(&g.num, &g.denom, ten)
	return groupsToDigits(groups), exp
}

type nrootGenerator struct {
	num        big.Int
	denom      big.Int
	newManager func() rootManager
}

func (g *nrootGenerator) Generate() (func() int, int) {
	manager := g.newManager()
	groups, exp := computeGroupsFromRational(
		&g.num, &g.denom, manager.Base(new(big.Int)))
	return computeRootDigits(groups, manager), exp
}

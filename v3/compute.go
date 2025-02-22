package sqroot

import (
	"math/big"
)

var (
	one                  = big.NewInt(1)
	two                  = big.NewInt(2)
	six                  = big.NewInt(6)
	ten                  = big.NewInt(10)
	fortyFive            = big.NewInt(45)
	fiftyFour            = big.NewInt(54)
	oneHundred           = big.NewInt(100)
	oneHundredSeventyOne = big.NewInt(171)
	oneThousand          = big.NewInt(1000)
)

type rootManager interface {
	Next(incr *big.Int)
	NextDigit(incr *big.Int)
	Base(result *big.Int) *big.Int
}

func computeGroupsFromRational(num, denom, base *big.Int) (
	groups func(result *big.Int) *big.Int, exp int) {
	num = new(big.Int).Set(num)
	denom = new(big.Int).Set(denom)
	base = new(big.Int).Set(base)
	for num.Cmp(denom) < 0 {
		exp--
		num.Mul(num, base)
	}
	if exp < 0 {
		exp++
		num.Div(num, base)
	}
	for num.Cmp(denom) >= 0 {
		exp++
		denom.Mul(denom, base)
	}
	groups = func(result *big.Int) *big.Int {
		if num.Sign() == 0 {
			return nil
		}
		num.Mul(num, base)
		result.DivMod(num, denom, num)
		return result
	}
	return
}

// remove for v4
func groupsToDigits(groups func(result *big.Int) *big.Int) func() int {
	var nextGroupHolder big.Int
	return func() int {
		nextGroup := groups(&nextGroupHolder)
		if nextGroup == nil {
			return -1
		}
		return int(nextGroup.Int64())
	}
}

func computeRootDigits(
	radicanGroups func(result *big.Int) *big.Int,
	manager rootManager) func() int {
	base := manager.Base(new(big.Int))
	incr := big.NewInt(1)
	remainder := big.NewInt(0)
	var nextGroupHolder big.Int
	return func() int {
		nextGroup := radicanGroups(&nextGroupHolder)
		if nextGroup == nil && remainder.Sign() == 0 {
			return -1
		}
		remainder.Mul(remainder, base)
		if nextGroup != nil {
			remainder.Add(remainder, nextGroup)
		}
		digit := 0
		for remainder.Cmp(incr) >= 0 {
			remainder.Sub(remainder, incr)
			digit++
			manager.Next(incr)
		}
		manager.NextDigit(incr)
		return digit
	}
}

type sqrtManager struct {
}

func newSqrtManager() rootManager {
	return sqrtManager{}
}

func (s sqrtManager) Next(incr *big.Int) {
	incr.Add(incr, two)
}

func (s sqrtManager) NextDigit(incr *big.Int) {
	incr.Sub(incr, one).Mul(incr, ten).Add(incr, one)
}

func (s sqrtManager) Base(result *big.Int) *big.Int {
	return result.Set(oneHundred)
}

type cubeRootManager struct {
	incr2 big.Int
}

func newCubeRootManager() rootManager {
	result := &cubeRootManager{}
	result.incr2.Set(six)
	return result
}

func (c *cubeRootManager) Next(incr *big.Int) {
	incr.Add(incr, &c.incr2)
	c.incr2.Add(&c.incr2, six)
}

func (c *cubeRootManager) NextDigit(incr *big.Int) {
	var temp big.Int
	incr.Mul(incr, oneHundred)
	incr.Sub(incr, temp.Mul(&c.incr2, fortyFive))
	incr.Add(incr, oneHundredSeventyOne)

	c.incr2.Mul(&c.incr2, ten).Sub(&c.incr2, fiftyFour)
}

func (c *cubeRootManager) Base(result *big.Int) *big.Int {
	return result.Set(oneThousand)
}

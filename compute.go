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

func nRoot(
	num, denom *big.Int, newManager func() rootManager) (
	mantissa func() int, exponent int) {
	num = new(big.Int).Set(num)
	denom = new(big.Int).Set(denom)
	base := newManager().Base(new(big.Int))
	exp := 0
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
	g := &nRootDigitGenerator{newManager: newManager}
	g.num.Set(num)
	g.denom.Set(denom)
	return g.iterator(), exp
}

type nRootDigitGenerator struct {
	num        big.Int
	denom      big.Int
	newManager func() rootManager
}

func (n *nRootDigitGenerator) iterator() func() int {
	manager := n.newManager()
	base := manager.Base(new(big.Int))
	incr := big.NewInt(1)
	remainder := big.NewInt(0)
	radicanGroups := n.generateRadicanGroups()
	return func() int {
		nextGroup := radicanGroups()
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

func (n *nRootDigitGenerator) generateRadicanGroups() func() *big.Int {
	num := new(big.Int).Set(&n.num)
	denom := new(big.Int).Set(&n.denom)
	base := n.newManager().Base(new(big.Int))
	return func() *big.Int {
		if num.Sign() == 0 {
			return nil
		}
		num.Mul(num, base)
		group, _ := new(big.Int).DivMod(num, denom, num)
		return group
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

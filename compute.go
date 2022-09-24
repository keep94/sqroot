package sqroot

import (
	"math/big"
)

var (
	oneHundred = big.NewInt(100)
	two        = big.NewInt(2)
	one        = big.NewInt(1)
	ten        = big.NewInt(10)
)

type mantissaSpec interface {
	Iterator() func() int
}

type sqrtSpec struct {
	num   big.Int
	denom big.Int
}

func (s *sqrtSpec) Iterator() func() int {
	incr := big.NewInt(1)
	remainder := big.NewInt(0)
	radicanGroups := generateQuotientBase100(&s.num, &s.denom)
	return func() int {
		nextGroup := radicanGroups()
		if nextGroup == nil && remainder.Sign() == 0 {
			return -1
		}
		remainder.Mul(remainder, oneHundred)
		if nextGroup != nil {
			remainder.Add(remainder, nextGroup)
		}
		digit := 0
		for remainder.Cmp(incr) >= 0 {
			remainder.Sub(remainder, incr)
			digit++
			incr.Add(incr, two)
		}
		incr.Sub(incr, one).Mul(incr, ten).Add(incr, one)
		return digit
	}
}

func generateQuotientBase100(num, denom *big.Int) func() *big.Int {
	num = new(big.Int).Set(num)
	denom = new(big.Int).Set(denom)
	return func() *big.Int {
		if num.Sign() == 0 {
			return nil
		}
		num.Mul(num, oneHundred)
		group, _ := new(big.Int).DivMod(num, denom, num)
		return group
	}
}

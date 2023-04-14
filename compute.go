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
	IteratorFrom(index int) func() int
	At(index int) int
	IsMemoize() bool
}

type sqrtSpec struct {
	num   big.Int
	denom big.Int
}

func (s *sqrtSpec) At(index int) int {
	return simpleAt(s.iterator(), index)
}

func (s *sqrtSpec) IteratorFrom(index int) func() int {
	return fastForward(s.iterator(), index)
}

func (s *sqrtSpec) IsMemoize() bool { return false }

func (s *sqrtSpec) iterator() func() int {
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

type limitSpec struct {
	delegate mantissaSpec
	limit    int
}

func withLimit(spec mantissaSpec, limit int) mantissaSpec {
	if limit < 0 {
		panic("limit must be non-negative")
	}
	if limit == 0 || spec == nil {
		return nil
	}
	ls, ok := spec.(*limitSpec)
	if ok {
		if limit >= ls.limit {
			return spec
		}
		return &limitSpec{delegate: ls.delegate, limit: limit}
	}
	return &limitSpec{delegate: spec, limit: limit}
}

func (l *limitSpec) At(index int) int {
	if index >= l.limit {
		return -1
	}
	return l.delegate.At(index)
}

func (l *limitSpec) IsMemoize() bool {
	return l.delegate.IsMemoize()
}

func (l *limitSpec) IteratorFrom(index int) func() int {
	if index > l.limit {
		index = l.limit
	}
	iter := l.delegate.IteratorFrom(index)
	return func() int {
		if index == l.limit {
			return -1
		}
		index++
		return iter()
	}
}

func withMemoize(spec mantissaSpec) mantissaSpec {
	if spec == nil {
		return nil
	}
	if spec.IsMemoize() {
		return spec
	}
	return newMemoizer(spec.IteratorFrom(0))
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

func simpleAt(iter func() int, index int) int {
	if index < 0 {
		return -1
	}
	return fastForward(iter, index)()
}

func fastForward(iter func() int, index int) func() int {
	if index < 0 {
		panic("index must be non-negative")
	}
	for i := 0; i < index && iter() != -1; i++ {
	}
	return iter
}

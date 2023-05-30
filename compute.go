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

type mantissaSpec interface {
	IteratorAt(index int) func() int
	At(index int) int
	IsMemoize() bool
	FirstN(n int) []int
}

type nRootSpec struct {
	num        big.Int
	denom      big.Int
	newManager func() rootManager
}

func (n *nRootSpec) At(index int) int {
	return simpleAt(n.iterator(), index)
}

func (n *nRootSpec) IteratorAt(index int) func() int {
	return fastForward(n.iterator(), index)
}

func (n *nRootSpec) IsMemoize() bool { return false }

func (n *nRootSpec) FirstN(size int) []int {
	panic("FirstN not supported")
}

func (n *nRootSpec) iterator() func() int {
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

func (n *nRootSpec) generateRadicanGroups() func() *big.Int {
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

func (l *limitSpec) IteratorAt(index int) func() int {
	if index > l.limit {
		index = l.limit
	}
	iter := l.delegate.IteratorAt(index)
	return func() int {
		if index == l.limit {
			return -1
		}
		index++
		return iter()
	}
}

func (l *limitSpec) FirstN(n int) []int {
	if n > l.limit {
		n = l.limit
	}
	return l.delegate.FirstN(n)
}

func withMemoize(spec mantissaSpec) mantissaSpec {
	if spec == nil {
		return nil
	}
	if spec.IsMemoize() {
		return spec
	}
	return newMemoizer(spec.IteratorAt(0))
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

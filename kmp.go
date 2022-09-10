package sqroot

import (
	"github.com/keep94/consume2"
)

// pattern must be non-empty
func ttable(pattern []int) []int {
	result := make([]int, len(pattern)+1)
	result[0] = -1
	posit := -1
	for i := 1; i < len(pattern); i++ {
		posit++
		result[i] = posit
		for posit != -1 && pattern[i] != pattern[posit] {
			posit = result[posit]
		}
	}
	result[len(pattern)] = posit + 1
	return result
}

type zeroPattern struct {
	consume2.Consumer[int]
	textIndex int
}

func (z *zeroPattern) Consume(x int) {
	z.Consumer.Consume(z.textIndex)
	z.textIndex++
}

type kmp struct {
	consume2.Consumer[int]
	pattern      []int
	table        []int
	textIndex    int
	patternIndex int
}

func (k *kmp) Consume(x int) {
	if x == k.pattern[k.patternIndex] {
		k.textIndex++
		k.patternIndex++
		if k.patternIndex == len(k.pattern) {
			k.Consumer.Consume(k.textIndex - k.patternIndex)
			k.patternIndex = k.table[k.patternIndex]
		}
		return
	}
	for k.patternIndex != -1 && k.pattern[k.patternIndex] != x {
		k.patternIndex = k.table[k.patternIndex]
	}
	k.patternIndex++
	k.textIndex++
}

package sqroot

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

func zeroPattern(f func() (Digit, bool)) func() int {
	return func() int {
		d, ok := f()
		if !ok {
			return -1
		}
		return d.Position
	}
}

func kmp(f func() (Digit, bool), pattern []int, reverse bool) func() int {
	kernel := newKmpKernel(pattern)
	return func() int {
		for {
			d, ok := f()
			if !ok {
				return -1
			}
			if kernel.Visit(d.Value) {
				if reverse {
					return d.Position
				}
				return d.Position + 1 - len(pattern)
			}
		}
	}
}

type kmpKernel struct {
	table        []int
	pattern      []int
	patternIndex int
}

func newKmpKernel(pattern []int) *kmpKernel {
	return &kmpKernel{
		table:   ttable(pattern),
		pattern: pattern,
	}
}

func (k *kmpKernel) Visit(digit int) bool {
	if digit == k.pattern[k.patternIndex] {
		k.patternIndex++
		if k.patternIndex == len(k.pattern) {
			k.patternIndex = k.table[k.patternIndex]
			return true
		}
		return false
	}
	for k.patternIndex != -1 && k.pattern[k.patternIndex] != digit {
		k.patternIndex = k.table[k.patternIndex]
	}
	k.patternIndex++
	return false
}

func patternReverse(pattern []int) []int {
	length := len(pattern)
	result := make([]int, length)
	for i := range pattern {
		result[length-i-1] = pattern[i]
	}
	return result
}

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

func zeroPattern(f func() positDigit) func() int {
	return func() int {
		pd := f()
		if !pd.Valid() {
			return -1
		}
		return pd.Posit
	}
}

func kmp(f func() positDigit, pattern []int, reverse bool) func() int {
	kernel := newKmpKernel(pattern)
	direction := 1
	if reverse {
		direction = -1
	}
	expectedIndex := -1
	return func() int {
		for {
			pd := f()
			if !pd.Valid() {
				return -1
			}
			if pd.Posit != expectedIndex {
				kernel.Reset()
			}
			expectedIndex = pd.Posit + direction
			if kernel.Visit(pd.Digit) {
				if reverse {
					return pd.Posit
				}
				return pd.Posit + 1 - len(pattern)
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

func (k *kmpKernel) Reset() {
	k.patternIndex = 0
}

func patternReverse(pattern []int) []int {
	length := len(pattern)
	result := make([]int, length)
	for i := range pattern {
		result[length-i-1] = pattern[i]
	}
	return result
}

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

func kmp(f func() positDigit, pattern []int) func() int {
	table := ttable(pattern)
	textIndex := 0
	patternIndex := 0
	return func() int {
		for {
			pd := f()
			if !pd.Valid() {
				return -1
			}
			if pd.Posit > textIndex {
				patternIndex = 0
				textIndex = pd.Posit
			}
			if pd.Digit == pattern[patternIndex] {
				textIndex++
				patternIndex++
				if patternIndex == len(pattern) {
					result := textIndex - patternIndex
					patternIndex = table[patternIndex]
					return result
				}
				continue
			}
			for patternIndex != -1 && pattern[patternIndex] != pd.Digit {
				patternIndex = table[patternIndex]
			}
			patternIndex++
			textIndex++
		}
	}
}

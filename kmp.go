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

func zeroPattern(f func() int) func() int {
	textIndex := 0
	return func() int {
		digit := f()
		if digit == -1 {
			return -1
		}
		textIndex++
		return textIndex - 1
	}
}

func kmp(f func() int, pattern []int) func() int {
	table := ttable(pattern)
	textIndex := 0
	patternIndex := 0
	return func() int {
		for {
			digit := f()
			if digit == -1 {
				return -1
			}
			if digit == pattern[patternIndex] {
				textIndex++
				patternIndex++
				if patternIndex == len(pattern) {
					result := textIndex - patternIndex
					patternIndex = table[patternIndex]
					return result
				}
				continue
			}
			for patternIndex != -1 && pattern[patternIndex] != digit {
				patternIndex = table[patternIndex]
			}
			patternIndex++
			textIndex++
		}
	}
}

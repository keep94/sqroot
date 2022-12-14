package sqroot

// Find returns a function that returns the next zero based index of the
// match for pattern in s. If s has a finite number of digits and there
// are no more matches for pattern, the returned function returns -1.
// Pattern is a sequence of digits between 0 and 9.
func Find(s Sequence, pattern []int) func() int {
	if len(pattern) == 0 {
		return zeroPattern(s.positDigitIter())
	}
	patternCopy := make([]int, len(pattern))
	copy(patternCopy, pattern)
	return kmp(s.positDigitIter(), patternCopy, false)
}

// FindFirst finds the zero based index of the first match of pattern in s.
// FindFirst returns -1 if pattern is not found only if s has a finite number
// of digits. If s has an infinite number of digits and pattern is not found,
// FindFirst will run forever. pattern is a sequence of digits between 0 and 9.
func FindFirst(s Sequence, pattern []int) int {
	iter := find(s, pattern)
	return iter()
}

// FindFirstN works like FindFirst but it finds the first n matches and
// returns the zero based index of each match. If s has a finite
// number of digits, FindFirstN may return fewer than n matches.
// Like FindFirst, FindFirstN may run forever if s has an infinite
// number of digits, and there are not n matches available.
// pattern is a sequence of digits between 0 and 9.
func FindFirstN(s Sequence, pattern []int, n int) []int {
	var result []int
	iter := find(s, pattern)
	for index := iter(); index != -1 && len(result) < n; index = iter() {
		result = append(result, index)
	}
	return result
}

// FindAll finds all the matches of pattern in s and returns the zero based
// index of each match. If s has an infinite number of digits, FindAll will
// run forever. pattern is a sequence of digits between 0 and 9.
func FindAll(s Sequence, pattern []int) []int {
	var result []int
	iter := find(s, pattern)
	for index := iter(); index != -1; index = iter() {
		result = append(result, index)
	}
	return result
}

// FindLast finds the zero based index of the last match of pattern in s.
// FindLast returns -1 if pattern is not found in s. If s has an infinite
// number of digits, FindLast will run forever. pattern is a sequence of
// digits between 0 and 9.
func FindLast(s Sequence, pattern []int) int {
	result := FindLastN(s, pattern, 1)
	if len(result) == 0 {
		return -1
	}
	return result[0]
}

// FindLastN works like FindLast but it finds the last n matches and
// returns the zero based index of each match. Last matches come first
// in returned array. If s has an infinite number of digits, FindLastN
// will run forever. pattern is a sequence of digits between 0 and 9.
func FindLastN(s Sequence, pattern []int, n int) []int {
	d, ok := s.(Digits)
	if ok {
		return d.findLastN(pattern, n)
	}
	return findLastN(s, pattern, n)
}

func find(s Sequence, pattern []int) func() int {
	if len(pattern) == 0 {
		return zeroPattern(s.positDigitIter())
	}
	return kmp(s.positDigitIter(), pattern, false)
}

func findLastN(s Sequence, pattern []int, n int) []int {
	if n == 0 {
		return make([]int, 0)
	}
	var buffer []int
	iter := find(s, pattern)
	var index int
	for index = iter(); index != -1 && len(buffer) < n; index = iter() {
		buffer = append(buffer, index)
	}
	posit := 0
	for ; index != -1; index = iter() {
		buffer[posit] = index
		posit = (posit + 1) % n
	}
	result := make([]int, len(buffer))
	nposit := 0
	for i := posit - 1; i >= 0; i-- {
		result[nposit] = buffer[i]
		nposit++
	}
	for i := len(buffer) - 1; i >= posit; i-- {
		result[nposit] = buffer[i]
		nposit++
	}
	return result
}

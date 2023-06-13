package sqroot

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func assertStartsAt(t *testing.T, s Sequence, start int) {
	t.Helper()
	iter := s.digitIter()
	d, ok := iter()
	assert.True(t, ok)
	assert.Equal(t, start, d.Position)
	assert.Equal(t, (start+1)%10, d.Value)
}

func assertRange(t *testing.T, s Sequence, start, end int) {
	t.Helper()
	if assertForwardRange(t, s, start, end) && hasReverse(s) {
		assertReverseRange(t, s, start, end)
	}
}

func assertForwardRange(t *testing.T, s Sequence, start, end int) bool {
	t.Helper()
	iter := s.digitIter()
	for i := start; i < end; i++ {
		d, ok := iter()
		if !assert.True(t, ok) {
			return false
		}
		if !assert.Equal(t, i, d.Position) {
			return false
		}
		if !assert.Equal(t, (i+1)%10, d.Value) {
			return false
		}
	}
	_, ok := iter()
	return assert.False(t, ok)
}

func assertReverseRange(t *testing.T, s Sequence, start, end int) bool {
	t.Helper()
	r := s.(reverseSequence)
	iter := r.reverseDigitIter()
	for i := end - 1; i >= start; i-- {
		d, ok := iter()
		if !assert.True(t, ok) {
			return false
		}
		if !assert.Equal(t, i, d.Position) {
			return false
		}
		if !assert.Equal(t, (i+1)%10, d.Value) {
			return false
		}
	}
	_, ok := iter()
	return assert.False(t, ok)
}

func assertEmpty(t *testing.T, s Sequence) {
	t.Helper()
	assertRange(t, s, 0, 0)
}

func hasReverse(s Sequence) bool {
	r, ok := s.(reverseSequence)
	return ok && r.canReverse()
}

func hasSubRange(s Sequence) bool {
	sr, ok := s.(subRangeSequence)
	return ok && sr.canSubRange()
}

func subRange(s Sequence, start, end int) Sequence {
	sr := s.(subRangeSequence)
	return sr.subRange(start, end)
}

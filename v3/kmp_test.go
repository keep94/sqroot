package sqroot

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTTable(t *testing.T) {
	pattern := []int{0, 1, 2, 3, 4, 5, 4, 0, 1, 3, 6, 7, 4, 8, 7, 0, 1, 2, 0, 5, 9}
	expect := []int{-1, 0, 0, 0, 0, 0, 0, 0, 1, 2, 0, 0, 0, 0, 0, 0, 1, 2, 3, 1, 0, 0}
	assert.Equal(t, expect, ttable(pattern))
}

func TestTTableAgain(t *testing.T) {
	pattern := []int{1, 2, 2, 1, 2, 1, 2, 2, 1, 2, 2, 1}
	expect := []int{-1, 0, 0, 0, 1, 2, 1, 2, 3, 4, 5, 3, 4}
	assert.Equal(t, expect, ttable(pattern))
}

func TestTTableSingle(t *testing.T) {
	assert.Equal(t, []int{-1, 0}, ttable([]int{3}))
}

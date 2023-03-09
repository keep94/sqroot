package sqroot

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPositionsBuilder(t *testing.T) {
	var pb PositionsBuilder
	pb.AddRange(0, 2).Add(4).Add(10).AddRange(-1, 3)
	pb.AddRange(15, 17)
	pb.AddRange(-3, -1)
	pb.AddRange(13, 15)
	pb.AddRange(17, 19)
	pb.Add(1)
	pb.AddRange(20, 25)
	pb.AddRange(21, 26)
	pb.AddRange(22, 23)
	assert.True(t, pb.unsorted)
	p := pb.Build()
	assert.False(t, pb.unsorted)
	assert.Len(t, pb.ranges, 0)
	expected := []positionRange{
		{Start: 0, End: 3},
		{Start: 4, End: 5},
		{Start: 10, End: 11},
		{Start: 13, End: 19},
		{Start: 20, End: 26},
	}
	assert.Equal(t, expected, p.ranges)
	assert.Equal(t, 26, p.limit())
}

func TestPositionsBuilderSorted(t *testing.T) {
	var pb PositionsBuilder
	pb.AddRange(0, 3).AddRange(1, 4).Add(2)
	pb.AddRange(4, 6).AddRange(10, 15).AddRange(6, 6).AddRange(7, 5)
	pb.AddRange(12, 17)
	for i := 100; i < 200; i++ {
		pb.Add(i)
	}
	assert.False(t, pb.unsorted)
	assert.Len(t, pb.ranges, 3)
	p := pb.Build()
	assert.False(t, pb.unsorted)
	assert.Len(t, pb.ranges, 0)
	expected := []positionRange{
		{Start: 0, End: 6},
		{Start: 10, End: 17},
		{Start: 100, End: 200},
	}
	assert.Equal(t, expected, p.ranges)
	assert.Equal(t, 200, p.limit())
}

func TestPositionsBuilderNegative(t *testing.T) {
	var pb PositionsBuilder
	pb.Add(-1)
	assert.Zero(t, pb.Build())
}

func TestPositionsBuilderZero(t *testing.T) {
	var pb PositionsBuilder
	assert.Zero(t, pb.Build())
}

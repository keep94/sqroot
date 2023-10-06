package sqroot

import (
	"sort"
)

// PositionsBuilder builds Positions objects. The zero value has no
// positions in it and is ready to use.
type PositionsBuilder struct {
	ranges   []PositionRange
	unsorted bool
}

// Add adds posit to this instance and returns this instance for chaining.
// If posit is negative, Add is a no-op.
func (p *PositionsBuilder) Add(posit int) *PositionsBuilder {
	return p.AddRange(posit, posit+1)
}

// AddRange adds a range of positions to this instance and returns this
// instance for chaining. The range is between start inclusive and end
// exclusive. AddRange ignores any negative positions within the specified
// range. If end <= start, AddRange is a no-op.
func (p *PositionsBuilder) AddRange(start, end int) *PositionsBuilder {
	if start < 0 {
		start = 0
	}
	if end <= start {
		return p
	}
	newRange := PositionRange{Start: start, End: end}
	length := len(p.ranges)
	if length == 0 {
		p.ranges = append(p.ranges, newRange)
		return p
	}
	if start < p.ranges[length-1].Start {
		p.ranges = append(p.ranges, newRange)
		p.unsorted = true
		return p
	}
	appendNotBefore(newRange, &p.ranges)
	return p
}

// Build builds a Positions instance from this builder and resets this builder
// so that it has no positions in it.
func (p *PositionsBuilder) Build() Positions {
	if !p.unsorted {
		result := p.ranges
		*p = PositionsBuilder{}
		return Positions{ranges: result}
	}
	sort.Slice(
		p.ranges,
		func(i, j int) bool {
			return p.ranges[i].Start < p.ranges[j].Start
		},
	)
	var result []PositionRange
	result = append(result, p.ranges[0])
	for _, prange := range p.ranges[1:] {
		appendNotBefore(prange, &result)
	}
	*p = PositionsBuilder{}
	return Positions{ranges: result}
}

// Positions represents a set of zero based positions for which to fetch
// digits. The zero value contains no positions.
type Positions struct {
	ranges []PositionRange
}

// UpTo returns the positions from 0 up to but not including end.
func UpTo(end int) Positions {
	var pb PositionsBuilder
	return pb.AddRange(0, end).Build()
}

// Between returns the positions from start up to but not including end.
func Between(start, end int) Positions {
	var pb PositionsBuilder
	return pb.AddRange(start, end).Build()
}

// Ranges returns a function that generates all the non overlapping ranges
// of positions in p. The returned function generates all the ranges in
// increasing order and returns false when there are no more.
func (p Positions) Ranges() func() (pr PositionRange, ok bool) {
	index := 0
	return func() (pr PositionRange, ok bool) {
		if index == len(p.ranges) {
			return
		}
		pr = p.ranges[index]
		ok = true
		index++
		return
	}
}

// End returns the last zero based position in p plus 1. If p is the zero
// value, End returns 0.
func (p Positions) End() int {
	length := len(p.ranges)
	if length == 0 {
		return 0
	}
	return p.ranges[length-1].End
}

// PositionRange is a single range of positions within a Positions instance.
type PositionRange struct {

	// The zero based starting position inclusive.
	Start int

	// The zero based ending position exclusive.
	End int
}

func appendNotBefore(item PositionRange, ranges *[]PositionRange) {
	length := len(*ranges)
	lastItem := &(*ranges)[length-1]
	if item.Start <= lastItem.End {
		if item.End > lastItem.End {
			lastItem.End = item.End
		}
	} else {
		*ranges = append(*ranges, item)
	}
}

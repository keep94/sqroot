package sqroot

import (
	"sort"
)

// PositionsBuilder builds Positions objects. The zero value has no
// positions in it and is ready to use.
type PositionsBuilder struct {
	ranges   []positionRange
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
	newRange := positionRange{Start: start, End: end}
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
	var result []positionRange
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
	ranges []positionRange
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

func (p Positions) limit() int {
	length := len(p.ranges)
	if length == 0 {
		return 0
	}
	return p.ranges[length-1].End
}

func (p Positions) filter() *positionsFilter {
	return &positionsFilter{ranges: p.ranges}
}

type positionRange struct {
	Start int
	End   int
}

type positionsFilter struct {
	ranges     []positionRange
	startIndex int
	limit      int
}

func (p *positionsFilter) Includes(posit int) bool {
	p.update(posit)
	return posit < p.limit
}

func (p *positionsFilter) update(posit int) {
	for p.startIndex < len(p.ranges) && p.ranges[p.startIndex].Start <= posit {
		p.limit = p.ranges[p.startIndex].End
		p.startIndex++
	}
}

func appendNotBefore(item positionRange, ranges *[]positionRange) {
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

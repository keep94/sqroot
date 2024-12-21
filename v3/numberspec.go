package sqroot

import (
	"math"
	"sync"
)

const (
	kMemoizerChunkSize = 100
	kMaxChunks         = math.MaxInt / kMemoizerChunkSize
)

// Digit represents a digit and its zero based position in a mantissa.
type Digit struct {

	// The 0 based position of the digit.
	Position int

	// The value of the digit. Always between 0 and 9.
	Value int
}

type numberSpec interface {
	IteratorAt(index, limit int) func() (Digit, bool)
	Scan(index, limit int, yield func(index, value int) bool)
	ScanValues(index, limit int, yield func(value int) bool)
	At(index int) int
	FirstN(n int) []int8
}

type memoizer struct {
	iter            func() int
	mu              sync.Mutex
	mustGrow        *sync.Cond
	updateAvailable *sync.Cond
	data            []int8
	maxLength       int
	done            bool
}

func newMemoizeSpec(iter func() int) numberSpec {
	result := &memoizer{iter: iter}
	result.mustGrow = sync.NewCond(&result.mu)
	result.updateAvailable = sync.NewCond(&result.mu)
	go result.run()
	return result
}

func (m *memoizer) At(index int) int {
	if index < 0 {
		return -1
	}
	data, ok := m.wait(index)
	if !ok {
		return -1
	}
	return int(data[index])
}

func (m *memoizer) FirstN(n int) []int8 {
	if n <= 0 {
		return nil
	}
	data, _ := m.wait(n - 1)
	if len(data) > n {
		return data[:n]
	}
	return data
}

func (m *memoizer) IteratorAt(index, limit int) func() (Digit, bool) {
	if index < 0 {
		panic("index must be non-negative")
	}
	var data []int8
	var ok, initialized bool
	return func() (Digit, bool) {
		if !initialized {
			data, ok = m.wait(index)
			initialized = true
		}
		if !ok || index >= limit {
			return Digit{}, false
		}
		result := Digit{Position: index, Value: int(data[index])}
		index++
		if index == len(data) {
			data, ok = m.wait(index)
		}
		return result, true
	}
}

func (m *memoizer) Scan(index, limit int, yield func(index, value int) bool) {
	if index < 0 {
		panic("index must be non-negative")
	}
	data, ok := m.wait(index)
	for ok && index < limit {
		if !yield(index, int(data[index])) {
			return
		}
		index++
		if index == len(data) {
			data, ok = m.wait(index)
		}
	}
}

func (m *memoizer) ScanValues(index, limit int, yield func(value int) bool) {
	if index < 0 {
		panic("index must be non-negative")
	}
	data, ok := m.wait(index)
	for ok && index < limit {
		if !yield(int(data[index])) {
			return
		}
		index++
		if index == len(data) {
			data, ok = m.wait(index)
		}
	}
}

func (m *memoizer) wait(index int) ([]int8, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if !m.done && m.maxLength <= index {
		chunkCount := index/kMemoizerChunkSize + 1

		// Have to prevent integer overflow in case index = math.MaxInt - 1
		if chunkCount > kMaxChunks {
			chunkCount = kMaxChunks
		}
		m.maxLength = kMemoizerChunkSize * chunkCount
		m.mustGrow.Signal()
	}
	for !m.done && len(m.data) <= index {
		m.updateAvailable.Wait()
	}
	return m.data, len(m.data) > index
}

func (m *memoizer) waitToGrow() {
	m.mu.Lock()
	defer m.mu.Unlock()
	for len(m.data) >= m.maxLength {
		m.mustGrow.Wait()
	}
}

func (m *memoizer) setData(data []int8, done bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data = data
	m.done = done
	m.updateAvailable.Broadcast()
}

func (m *memoizer) run() {
	var data []int8
	for i := 0; i < kMaxChunks; i++ {
		m.waitToGrow()
		for j := 0; j < kMemoizerChunkSize; j++ {
			x := m.iter()
			if digitOutOfRange(x) {
				m.setData(data, true)
				return
			}
			data = append(data, int8(x))
		}
		m.setData(data, false)
	}
	m.setData(data, true)
}

type limitSpec struct {
	delegate numberSpec
	limit    int
}

func withLimit(spec numberSpec, limit int) numberSpec {
	if limit <= 0 || spec == nil {
		return nil
	}
	ls, ok := spec.(*limitSpec)
	if ok {
		if limit >= ls.limit {
			return spec
		}
		return &limitSpec{delegate: ls.delegate, limit: limit}
	}
	return &limitSpec{delegate: spec, limit: limit}
}

func (l *limitSpec) At(index int) int {
	if index >= l.limit {
		l.delegate.At(l.limit)
		return -1
	}
	return l.delegate.At(index)
}

func (l *limitSpec) IteratorAt(index, limit int) func() (Digit, bool) {
	index = min(index, l.limit)
	limit = min(limit, l.limit)
	return l.delegate.IteratorAt(index, limit)
}

func (l *limitSpec) Scan(index, limit int, yield func(index, value int) bool) {
	index = min(index, l.limit)
	limit = min(limit, l.limit)
	l.delegate.Scan(index, limit, yield)
}

func (l *limitSpec) ScanValues(index, limit int, yield func(value int) bool) {
	index = min(index, l.limit)
	limit = min(limit, l.limit)
	l.delegate.ScanValues(index, limit, yield)
}

func (l *limitSpec) FirstN(n int) []int8 {
	if n > l.limit {
		n = l.limit
	}
	return l.delegate.FirstN(n)
}

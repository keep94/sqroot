package sqroot

import (
	"math"
	"sync"
)

const (
	kMemoizerChunkSize = 100
	kMaxChunks         = math.MaxInt / kMemoizerChunkSize
)

type memoizer struct {
	iter            func() int
	mu              sync.Mutex
	mustGrow        *sync.Cond
	updateAvailable *sync.Cond
	data            []int
	maxLength       int
	done            bool
}

func newMemoizer(iter func() int) *memoizer {
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
	return data[index]
}

func (m *memoizer) FirstN(n int) []int {
	if n <= 0 {
		return nil
	}
	data, _ := m.wait(n - 1)
	if len(data) > n {
		return data[:n]
	}
	return data
}

func (m *memoizer) IsMemoize() bool { return true }

func (m *memoizer) IteratorAt(index int) func() int {
	if index < 0 {
		panic("index must be non-negative")
	}
	data, ok := m.wait(index)
	return func() int {
		if !ok {
			return -1
		}
		result := data[index]
		index++
		if index == len(data) {
			data, ok = m.wait(index)
		}
		return result
	}
}

func (m *memoizer) wait(index int) ([]int, bool) {
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
	for !m.done && len(m.data) < m.maxLength {
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

func (m *memoizer) setData(data []int, done bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data = data
	m.done = done
	m.updateAvailable.Broadcast()
}

func (m *memoizer) run() {
	var data []int
	for {
		m.waitToGrow()
		for i := 0; i < kMemoizerChunkSize; i++ {
			x := m.iter()
			if x == -1 {
				m.setData(data, true)
				return
			}
			data = append(data, x)
		}
		m.setData(data, false)
	}
}

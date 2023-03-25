package sqroot

import (
	"sync"
)

const (
	kMemoizerChunkSize = 100
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

func (m *memoizer) Memoize() bool { return true }

func (m *memoizer) Iterator() func() int {
	posit := 0
	data, ok := m.wait(posit)
	return func() int {
		if !ok {
			return -1
		}
		result := data[posit]
		posit++
		if posit == len(data) {
			data, ok = m.wait(posit)
		}
		return result
	}
}

func (m *memoizer) wait(index int) ([]int, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if !m.done && m.maxLength <= index {
		m.maxLength = kMemoizerChunkSize * ((index / kMemoizerChunkSize) + 1)
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

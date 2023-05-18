package sqroot

import (
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMemoize(t *testing.T) {
	n := Sqrt(5)
	nm := n.WithMemoize()
	assert.True(t, nm.Mantissa().IsMemoize())
	assert.False(t, n.Mantissa().IsMemoize())
	expected := fmt.Sprintf("%.10000g", n)
	var actual [10]string
	var wg sync.WaitGroup
	for i := 0; i < len(actual); i++ {
		wg.Add(1)
		go func(index int) {
			actual[index] = fmt.Sprintf("%.10000g", nm)
			wg.Done()
		}(i)
	}
	wg.Wait()
	for i := 0; i < len(actual); i++ {
		assert.Equal(t, expected, actual[i])
	}
}

func TestMemoizeAt(t *testing.T) {
	n := Sqrt(7)
	d := AllDigits(n.Mantissa().WithSignificant(10000))
	var expected, actual1, actual2 [10000]int
	iter := d.Items()
	i := 0
	for digit, ok := iter(); ok; digit, ok = iter() {
		expected[i] = digit.Value
		i++
	}
	nm := n.WithMemoize()
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		m := nm.Mantissa()
		for i := 9999; i >= 0; i-- {
			actual1[i] = m.At(i)
		}
		wg.Done()
	}()
	go func() {
		m := nm.Mantissa()
		for i := 0; i < 10000; i++ {
			actual2[i] = m.At(i)
		}
		wg.Done()
	}()
	wg.Wait()
	assert.Equal(t, expected, actual1)
	assert.Equal(t, expected, actual2)
}

func TestMemoizeOutOfBounds(t *testing.T) {
	n := Sqrt(111).WithSignificant(1000)
	expectedDigits := AllDigits(n.Mantissa())
	nm := n.WithMemoize()
	assert.Equal(t, 2, nm.Exponent())
	m := nm.Mantissa()
	assert.Equal(t, -1, m.At(1000))
	assert.Equal(t, -1, m.At(-1))
	assert.Equal(t, expectedDigits.At(999), m.At(999))
	assert.Equal(t, expectedDigits.At(0), m.At(0))
	assert.Equal(t, expectedDigits.Sprint(), AllDigits(m).Sprint())
}

func TestMemoizeOutOfBounds2(t *testing.T) {
	n := Sqrt(111)
	expectedDigits := AllDigits(n.Mantissa().WithSignificant(1000))
	nm := n.WithMemoize().WithSignificant(1000)
	assert.Equal(t, 2, nm.Exponent())
	m := nm.Mantissa()
	assert.Equal(t, -1, m.At(1000))
	assert.Equal(t, -1, m.At(-1))
	assert.Equal(t, expectedDigits.At(999), m.At(999))
	assert.Equal(t, expectedDigits.At(0), m.At(0))
	assert.Equal(t, expectedDigits.Sprint(), AllDigits(m).Sprint())
}

func TestMemoizeOddBoundary(t *testing.T) {
	n := Sqrt(97)
	var pb PositionsBuilder
	exdigits := GetDigits(n.Mantissa(), pb.AddRange(153, 158).Build())
	n = n.WithSignificant(158).WithMemoize()
	m := n.Mantissa()
	assert.Equal(t, exdigits.At(153), m.At(153))
	assert.Equal(t, exdigits.At(154), m.At(154))
	assert.Equal(t, -1, m.At(158))
	assert.Equal(t, exdigits.At(155), m.At(155))
	assert.Equal(t, exdigits.At(156), m.At(156))
	assert.Equal(t, exdigits.At(157), m.At(157))
	start153 := m.WithStart(153)
	assert.Equal(t, exdigits.Sprint(), AllDigits(start153).Sprint())
	assert.Zero(t, AllDigits(m.WithStart(158)))
	pattern := []int{m.At(153), m.At(154), m.At(155), m.At(156), m.At(157)}
	assert.Equal(t, 153, FindFirst(start153, pattern))
	assert.Equal(t, 153, FindLast(start153, pattern))
	start154 := m.WithStart(154)
	assert.Equal(t, -1, FindFirst(start154, pattern))
	assert.Equal(t, -1, FindLast(start154, pattern))
}

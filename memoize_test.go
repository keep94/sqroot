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
	assert.True(t, nm.Mantissa().Memoize())
	assert.False(t, n.Mantissa().Memoize())
	expected := fmt.Sprintf("%.10000g", n)
	var actual1 string
	var actual2 string
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		actual1 = fmt.Sprintf("%.10000g", nm)
		wg.Done()
	}()
	go func() {
		actual2 = fmt.Sprintf("%.10000g", nm)
		wg.Done()
	}()
	wg.Wait()
	assert.Equal(t, expected, actual1)
	assert.Equal(t, expected, actual2)
}

func TestMemoizeAt(t *testing.T) {
	n := Sqrt(7)
	d := n.Mantissa().WithSignificant(10000).Digits()
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
	expectedDigits := n.Mantissa().Digits()
	nm := n.WithMemoize().WithMemoize()
	assert.Equal(t, 2, nm.Exponent())
	m := nm.Mantissa()
	assert.Equal(t, -1, m.At(1000))
	assert.Equal(t, -1, m.At(-1))
	assert.Equal(t, expectedDigits.At(999), m.At(999))
	assert.Equal(t, expectedDigits.At(0), m.At(0))
	assert.Equal(t, expectedDigits.Sprint(), m.Digits().Sprint())
}

func TestMemoizeOutOfBounds2(t *testing.T) {
	n := Sqrt(111)
	expectedDigits := n.Mantissa().WithSignificant(1000).Digits()
	nm := n.WithMemoize().WithSignificant(1000)
	assert.Equal(t, 2, nm.Exponent())
	m := nm.Mantissa()
	assert.Equal(t, -1, m.At(1000))
	assert.Equal(t, -1, m.At(-1))
	assert.Equal(t, expectedDigits.At(999), m.At(999))
	assert.Equal(t, expectedDigits.At(0), m.At(0))
	assert.Equal(t, expectedDigits.Sprint(), m.Digits().Sprint())
}

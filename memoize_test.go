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
	assert.True(t, nm.IsMemoize())
	assert.False(t, n.IsMemoize())
	expected := fmt.Sprintf("%.10000g", n)
	var actual [10]string
	var wg sync.WaitGroup
	for i := range actual {
		wg.Add(1)
		go func(index int) {
			actual[index] = fmt.Sprintf("%.10000g", nm)
			wg.Done()
		}(i)
	}
	wg.Wait()
	for i := range actual {
		assert.Equal(t, expected, actual[i])
	}
}

func TestMemoizeAt(t *testing.T) {
	n := Sqrt(7)
	d := AllDigits(n.WithSignificant(10000))
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
		for i := 9999; i >= 0; i-- {
			actual1[i] = nm.At(i)
		}
		wg.Done()
	}()
	go func() {
		for i := 0; i < 10000; i++ {
			actual2[i] = nm.At(i)
		}
		wg.Done()
	}()
	wg.Wait()
	assert.Equal(t, expected, actual1)
	assert.Equal(t, expected, actual2)
}

func TestMemoizeOutOfBounds(t *testing.T) {
	n := Sqrt(111).WithSignificant(1000)
	expectedDigits := AllDigits(n)
	nm := n.WithMemoize()
	assert.Equal(t, 2, nm.Exponent())
	assert.Equal(t, -1, nm.At(1000))
	assert.Equal(t, -1, nm.At(-1))
	assert.Equal(t, expectedDigits.At(999), nm.At(999))
	assert.Equal(t, expectedDigits.At(0), nm.At(0))
	assert.Equal(t, expectedDigits.Sprint(), AllDigits(nm).Sprint())
}

func TestMemoizeOutOfBounds2(t *testing.T) {
	n := Sqrt(111)
	expectedDigits := AllDigits(n.WithSignificant(1000))
	nm := n.WithMemoize().WithSignificant(1000)
	assert.Equal(t, 2, nm.Exponent())
	assert.Equal(t, -1, nm.At(1000))
	assert.Equal(t, -1, nm.At(-1))
	assert.Equal(t, expectedDigits.At(999), nm.At(999))
	assert.Equal(t, expectedDigits.At(0), nm.At(0))
	assert.Equal(t, expectedDigits.Sprint(), AllDigits(nm).Sprint())
}

func TestMemoizeOddBoundary(t *testing.T) {
	n := Sqrt(97)
	var pb PositionsBuilder
	exdigits := GetDigits(n, pb.AddRange(153, 158).Build())
	n = n.WithSignificant(158).WithMemoize()
	assert.Equal(t, exdigits.At(153), n.At(153))
	assert.Equal(t, exdigits.At(154), n.At(154))
	assert.Equal(t, -1, n.At(158))
	assert.Equal(t, exdigits.At(155), n.At(155))
	assert.Equal(t, exdigits.At(156), n.At(156))
	assert.Equal(t, exdigits.At(157), n.At(157))
	start153 := n.WithStart(153)
	assert.Equal(t, exdigits.Sprint(), AllDigits(start153).Sprint())
	assert.Zero(t, AllDigits(n.WithStart(158)))
	pattern := []int{n.At(153), n.At(154), n.At(155), n.At(156), n.At(157)}
	assert.Equal(t, 153, FindFirst(start153, pattern))
	assert.Equal(t, 153, FindLast(start153, pattern))
	start154 := n.WithStart(154)
	assert.Equal(t, -1, FindFirst(start154, pattern))
	assert.Equal(t, -1, FindLast(start154, pattern))
}

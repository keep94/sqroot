package sqroot

import (
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMemoize(t *testing.T) {
	expected := fmt.Sprintf("%.10000g", Sqrt(5))
	n := Sqrt(5)
	var actual [10]string
	var wg sync.WaitGroup
	for i := range actual {
		wg.Add(1)
		go func(index int) {
			actual[index] = fmt.Sprintf("%.10000g", n)
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
	var expected [10000]int
	var actual [10][10000]int
	for i := range expected {
		expected[i] = n.At(i)
	}
	n1 := Sqrt(7)
	var wg sync.WaitGroup
	for i := 0; i < len(actual)/2; i++ {
		wg.Add(2)
		go func(idx int) {
			for i := 9999; i >= 0; i-- {
				actual[idx][i] = n1.At(i)
			}
			wg.Done()
		}(2 * i)
		go func(idx int) {
			for i := 0; i < 10000; i++ {
				actual[idx][i] = n1.At(i)
			}
			wg.Done()
		}(2*i + 1)
	}
	wg.Wait()
	for i := range actual {
		assert.Equal(t, expected, actual[i])
	}
}

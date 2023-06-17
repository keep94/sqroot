package sqroot

import (
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMemoize(t *testing.T) {
	n := Sqrt(5)
	expected := fmt.Sprintf("%.10000g", n)
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
	var expected, actual1, actual2 [10000]int
	for i := range expected {
		expected[i] = n.At(i)
	}
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		for i := 9999; i >= 0; i-- {
			actual1[i] = n.At(i)
		}
		wg.Done()
	}()
	go func() {
		for i := 0; i < 10000; i++ {
			actual2[i] = n.At(i)
		}
		wg.Done()
	}()
	wg.Wait()
	assert.Equal(t, expected, actual1)
	assert.Equal(t, expected, actual2)
}

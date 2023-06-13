package sqroot

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFPrintP(t *testing.T) {
	p := new(PositionsBuilder).AddRange(5, 8).AddRange(45, 50).Build()
	d := GetDigits(fakeNumber, p)
	var builder strings.Builder
	num, err := d.Fprint(&builder, DigitsPerRow(11), DigitsPerColumn(10))
	assert.Equal(t, 27, num)
	assert.NoError(t, err)
	expected := `  0......678.. .
44  .67890`
	assert.Equal(t, expected, builder.String())
}

func TestPrintP(t *testing.T) {
	p := new(PositionsBuilder).AddRange(5, 8).AddRange(45, 50).Build()
	actual := GetDigits(fakeNumber, p).Sprint(
		DigitsPerRow(11), DigitsPerColumn(10))
	expected := `  0......678.. .
44  .67890`
	assert.Equal(t, expected, actual)
}

func TestPrintDigits(t *testing.T) {
	var pb PositionsBuilder
	p := pb.AddRange(100, 200).AddRange(300, 400).AddRange(500, 600).Build()
	d := GetDigits(fakeNumber, p)
	pb.AddRange(0, 101).AddRange(200, 300).AddRange(399, 500)
	q := pb.AddRange(700, 800).Build()
	actual := Sprint(d, q, DigitsPerRow(10))
	expected := `100  1.... .....
390  ..... ....0`
	assert.Equal(t, expected, actual)
}

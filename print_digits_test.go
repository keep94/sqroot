package sqroot

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrintP(t *testing.T) {
	p := new(PositionsBuilder).AddRange(5, 8).AddRange(45, 50).Build()
	actual := GetDigits(fakeMantissa, p).Sprint(
		DigitsPerRow(11), DigitsPerColumn(10))
	expected := `  0......678.. .
44  .67890`
	assert.Equal(t, expected, actual)
}

func TestPrintP2(t *testing.T) {
	var pb PositionsBuilder
	p := pb.AddRange(0, 2).AddRange(3, 5).AddRange(8, 11).Build()
	digits := GetDigits(fakeMantissa, p)
	actual := digits.Sprint(DigitsPerRow(11), DigitsPerColumn(10))
	expected := `0.12.45...90 1`
	assert.Equal(t, expected, actual)
}

func TestPrintPGaps(t *testing.T) {
	var pb PositionsBuilder
	p := pb.AddRange(22, 44).AddRange(66, 77).Build()
	digits := GetDigits(fakeMantissa, p)
	actual := digits.Sprint(DigitsPerRow(11), DigitsPerColumn(10))
	expected := `22  3456789012 3
33  4567890123 4
66  7890123456 7`
	assert.Equal(t, expected, actual)
}

func TestPrintPGaps2(t *testing.T) {
	var pb PositionsBuilder
	pb.AddRange(0, 10).AddRange(11, 21).AddRange(33, 43).AddRange(66, 76)
	p := pb.Build()
	digits := GetDigits(fakeMantissa, p)
	actual := digits.Sprint(DigitsPerRow(11), DigitsPerColumn(10))
	expected := `  0.1234567890 .
11  2345678901 .
33  4567890123 .
66  7890123456`
	assert.Equal(t, expected, actual)
}

func TestPrintPGaps3(t *testing.T) {
	var pb PositionsBuilder
	p := pb.AddRange(21, 33).AddRange(65, 77).Build()
	digits := GetDigits(fakeMantissa, p)
	actual := digits.Sprint(
		DigitsPerRow(11), DigitsPerColumn(10), MissingDigit('-'))
	expected := `11  ---------- 2
22  3456789012 3
55  ---------- 6
66  7890123456 7`
	assert.Equal(t, expected, actual)
}

func TestPrintPNoShowCount(t *testing.T) {
	var pb PositionsBuilder
	p := pb.AddRange(21, 33).AddRange(65, 77).Build()
	digits := GetDigits(fakeMantissa, p)
	actual := digits.Sprint(
		DigitsPerRow(11), DigitsPerColumn(10), ShowCount(false))
	expected := `0........... .
  .......... 2
  3456789012 3
  .......... .
  .......... .
  .......... 6
  7890123456 7`
	assert.Equal(t, expected, actual)
}

func TestPrintPNarrow(t *testing.T) {
	var pb PositionsBuilder
	p := pb.Add(3).Add(5).Add(8).Build()
	digits := GetDigits(fakeMantissa, p)
	actual := digits.Sprint(DigitsPerRow(1))
	expected := `3  4
5  6
8  9`
	assert.Equal(t, expected, actual)
}

func TestPrintPDefaults(t *testing.T) {
	var pb PositionsBuilder
	p := pb.AddRange(0, 75).Build()
	digits := GetDigits(fakeMantissa, p)
	actual := digits.Sprint()
	expected := `  0.12345 67890 12345 67890 12345 67890 12345 67890 12345 67890
50  12345 67890 12345 67890 12345`
	assert.Equal(t, expected, actual)
}

func TestPrintPTooShort(t *testing.T) {
	n := Sqrt(100489)
	p := new(PositionsBuilder).AddRange(3, 5).Build()
	digits := GetDigits(n.Mantissa(), p)
	assert.Zero(t, digits)
	assert.Empty(t, digits.Sprint())
}

func TestPrintPZero(t *testing.T) {
	var m Mantissa
	p := new(PositionsBuilder).AddRange(3, 5).Build()
	digits := GetDigits(m, p)
	assert.Zero(t, digits)
	assert.Empty(t, digits.Sprint())
}

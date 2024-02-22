package sqroot

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWriteZeroDigits(t *testing.T) {
	n := fakeNumber()
	assert.Equal(t, "\n", Swrite(n.WithEnd(0)))
	assert.Equal(t, "\n", Swrite(n.WithStart(5).WithEnd(5)))
}

func TestWriteNoOptions(t *testing.T) {
	n := fakeNumber()
	actual := Swrite(n.WithEnd(12))
	expected := "0  12345 67890 12\n"
	assert.Equal(t, expected, actual)
}

func TestWriteLessThanOneRow(t *testing.T) {
	n := fakeNumber()
	actual := Swrite(
		n.WithEnd(12), DigitsPerRow(12), DigitsPerColumn(0))
	expected := "0  123456789012\n"
	assert.Equal(t, expected, actual)
}

func TestWriteColumns(t *testing.T) {
	n := fakeNumber()
	actual := Swrite(
		n.WithEnd(12),
		DigitsPerColumn(4),
		DigitsPerRow(0),
		ShowCount(false))
	expected := "1234 5678 9012\n"
	assert.Equal(t, expected, actual)
}

func TestWriteColumnsShow(t *testing.T) {
	n := fakeNumber()
	actual := Swrite(n.WithEnd(12), DigitsPerColumn(5), DigitsPerRow(0))
	expected := "0  12345 67890 12\n"
	assert.Equal(t, expected, actual)
}

func TestWriteRows10(t *testing.T) {
	n := fakeNumber()
	actual := Swrite(
		n.WithEnd(110),
		DigitsPerRow(10),
		DigitsPerColumn(0),
		ShowCount(false))
	expected := `1234567890
1234567890
1234567890
1234567890
1234567890
1234567890
1234567890
1234567890
1234567890
1234567890
1234567890
`
	assert.Equal(t, expected, actual)
}

func TestWriteRows10LeadingDecimalNoExtraLF(t *testing.T) {
	n := fakeNumber()
	actual := Swrite(
		n.WithEnd(110),
		DigitsPerRow(10),
		DigitsPerColumn(0),
		ShowCount(false),
		TrailingLF(false),
		LeadingDecimal(true))
	expected := `0.1234567890
  1234567890
  1234567890
  1234567890
  1234567890
  1234567890
  1234567890
  1234567890
  1234567890
  1234567890
  1234567890`
	assert.Equal(t, expected, actual)
}

func TestWriteRows10Between(t *testing.T) {
	n := fakeNumber()
	actual := Swrite(
		n.WithStart(55).WithEnd(110),
		DigitsPerRow(10),
		DigitsPerColumn(0),
		ShowCount(false))
	expected := `..........
..........
..........
..........
..........
.....67890
1234567890
1234567890
1234567890
1234567890
1234567890
`
	assert.Equal(t, expected, actual)
}

func TestWriteRows10Columns(t *testing.T) {
	n := fakeNumber()
	actual := Swrite(
		n.WithEnd(110),
		DigitsPerRow(10),
		DigitsPerColumn(10),
		ShowCount(false))
	expected := `1234567890
1234567890
1234567890
1234567890
1234567890
1234567890
1234567890
1234567890
1234567890
1234567890
1234567890
`
	assert.Equal(t, expected, actual)
}

func TestWriteRows11Columns(t *testing.T) {
	n := fakeNumber()
	actual := Swrite(
		n.WithEnd(110),
		DigitsPerRow(11),
		DigitsPerColumn(10),
		ShowCount(false))
	expected := `1234567890 1
2345678901 2
3456789012 3
4567890123 4
5678901234 5
6789012345 6
7890123456 7
8901234567 8
9012345678 9
0123456789 0
`
	assert.Equal(t, expected, actual)
}

func TestWriteRows10Show(t *testing.T) {
	n := fakeNumber()
	actual := Swrite(
		n.WithEnd(110), DigitsPerRow(10), DigitsPerColumn(0))
	expected := `  0  1234567890
 10  1234567890
 20  1234567890
 30  1234567890
 40  1234567890
 50  1234567890
 60  1234567890
 70  1234567890
 80  1234567890
 90  1234567890
100  1234567890
`
	assert.Equal(t, expected, actual)
}

func TestWriteRows10ColumnsShow(t *testing.T) {
	n := fakeNumber()
	actual := Swrite(
		n.WithEnd(110), DigitsPerRow(10), DigitsPerColumn(10))
	expected := `  0  1234567890
 10  1234567890
 20  1234567890
 30  1234567890
 40  1234567890
 50  1234567890
 60  1234567890
 70  1234567890
 80  1234567890
 90  1234567890
100  1234567890
`
	assert.Equal(t, expected, actual)
}

func TestWriteRows11ColumnsShow(t *testing.T) {
	n := fakeNumber()
	actual := Swrite(
		n.WithEnd(110), DigitsPerRow(11), DigitsPerColumn(10))
	expected := ` 0  1234567890 1
11  2345678901 2
22  3456789012 3
33  4567890123 4
44  5678901234 5
55  6789012345 6
66  7890123456 7
77  8901234567 8
88  9012345678 9
99  0123456789 0
`
	assert.Equal(t, expected, actual)
}

func TestWriteRows11ColumnsShow109(t *testing.T) {
	n := fakeNumber()
	actual := Swrite(
		n.WithEnd(109), DigitsPerRow(11), DigitsPerColumn(10))
	expected := ` 0  1234567890 1
11  2345678901 2
22  3456789012 3
33  4567890123 4
44  5678901234 5
55  6789012345 6
66  7890123456 7
77  8901234567 8
88  9012345678 9
99  0123456789
`
	assert.Equal(t, expected, actual)
}

func TestWriteRows11ColumnsShow111(t *testing.T) {
	n := fakeNumber()
	actual := Swrite(
		n.WithEnd(111),
		DigitsPerRow(11),
		DigitsPerColumn(10),
		ShowCount(true))
	expected := `  0  1234567890 1
 11  2345678901 2
 22  3456789012 3
 33  4567890123 4
 44  5678901234 5
 55  6789012345 6
 66  7890123456 7
 77  8901234567 8
 88  9012345678 9
 99  0123456789 0
110  1
`
	assert.Equal(t, expected, actual)
}

func TestWriteWithBetween(t *testing.T) {
	n := fakeNumber()
	actual := Swrite(
		n.WithStart(50).WithEnd(70),
		DigitsPerRow(11),
		DigitsPerColumn(10),
		MissingDigit('-'))
	expected := `44  ------1234 5
55  6789012345 6
66  7890
`
	assert.Equal(t, expected, actual)
}

func TestWriteWithStart(t *testing.T) {
	number := fakeNumber()
	actual := Swrite(number.WithStart(502).WithEnd(505))
	expected := "500  ..345\n"
	assert.Equal(t, expected, actual)
}

func TestWriteCountBytes(t *testing.T) {
	w := &maxBytesWriter{maxBytes: 100000}

	// Prints 200 rows. Each row 10 columns 6 chars per column + (4+2) chars
	// for the margin. 66*200=13200 bytes
	n, err := Fwrite(w, fakeNumber().WithEnd(10000))
	assert.Equal(t, 13200, n)
	assert.NoError(t, err)
}

func TestWRiteCountBytes2(t *testing.T) {
	w := &maxBytesWriter{maxBytes: 10000}

	// Prints 20 rows. Each row 10 columns 6 chars per column + (3+2) chars
	// for the margin. 65*20=1300 bytes
	n, err := Fwrite(w, fakeNumber().WithEnd(1000))
	assert.Equal(t, 1300, n)
	assert.NoError(t, err)
}

func TestWriteErrorAtAllStages(t *testing.T) {
	number := fakeNumber()

	// Internally Fprint uses a bufio.Writer which defaults to a buffer size
	// of 4096 bytes (could change in future go versions). This means writing
	// can fail only when buffer fills up or is flushed. We pick 601 which
	// is prime and small enough that we test an error happening when manually
	// flushing the buffer. There exists a k such that 3*4096 < 601*k < 13200.
	for i := 0; i < 13200; i += 601 {
		w := &maxBytesWriter{maxBytes: i}
		n, err := Fwrite(w, number.WithEnd(10000))
		assert.Equal(t, i, n)
		assert.Error(t, err)
	}
}

func TestWriteErrorAtAllStages2(t *testing.T) {
	number := fakeNumber()

	// Set buffer size to 1 so that we get extra coverage.
	for i := 0; i < 1300; i++ {
		w := &maxBytesWriter{maxBytes: i}
		n, err := Fwrite(w, number.WithEnd(1000), bufferSize(1))
		assert.Equal(t, i, n)
		assert.Error(t, err)
	}
}

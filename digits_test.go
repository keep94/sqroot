package sqroot

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDigits(t *testing.T) {
	positions := new(PositionsBuilder).Add(25).Add(15).Add(50).Build()
	digits := GetDigits(Sqrt(2), positions)
	assert.Equal(t, 5, digits.At(15))
	assert.Equal(t, 7, digits.At(25))
	assert.Equal(t, 4, digits.At(50))
	assert.Equal(t, -1, digits.At(26))
	assert.Equal(t, -1, digits.At(51))
	assert.Equal(t, -1, digits.At(14))
	checkFullIter(t, digits.Items(), 15, 5, 25, 7, 50, 4)
	checkFullIter(t, digits.WithStart(25).Items(), 25, 7, 50, 4)
	assert.Equal(t, 15, digits.Min())
	assert.Equal(t, 50, digits.Max())
	checkFullIter(t, digits.ReverseItems(), 50, 4, 25, 7, 15, 5)
	checkFullIter(t, digits.WithEnd(50).ReverseItems(), 25, 7, 15, 5)
}

func TestGetDigitsFromDigits(t *testing.T) {
	var pb PositionsBuilder
	pb.AddRange(0, 100).AddRange(200, 300).AddRange(400, 500)
	digits := GetDigits(Sqrt(2), pb.Build())
	pb.AddRange(100, 200).AddRange(300, 400).AddRange(500, 600)
	assert.Zero(t, GetDigits(digits, pb.Build()))
}

func TestGetDigitsFromDigits2(t *testing.T) {
	n := Sqrt(2)
	var pb PositionsBuilder
	pb.AddRange(100, 200).AddRange(300, 400).AddRange(500, 600)
	digits := GetDigits(n, pb.Build())

	pb.AddRange(0, 101).AddRange(200, 301).Add(500).AddRange(1000, 2000000000)
	digits = GetDigits(digits, pb.Build())
	expected := GetDigits(n, pb.Add(100).Add(300).Add(500).Build())
	assert.Equal(t, expected.Sprint(), digits.Sprint())
}

func TestGetDigitsFromDigitsPick(t *testing.T) {
	n := Sqrt(2)
	var pb PositionsBuilder
	pb.AddRange(100, 200).AddRange(300, 400).AddRange(500, 600)
	digits := GetDigits(n, pb.Build())
	pb.Add(153).Add(200)
	digits = GetDigits(digits, pb.Build())
	expected := GetDigits(n, pb.Add(153).Build())
	assert.Equal(t, expected.Sprint(), digits.Sprint())
}

func TestDigitsFinite(t *testing.T) {
	n := Sqrt(2).WithSignificant(50)
	positions := new(PositionsBuilder).Add(25).Add(15).Add(50).Build()
	digits := GetDigits(n, positions)
	assert.Equal(t, 5, digits.At(15))
	assert.Equal(t, 7, digits.At(25))
	assert.Equal(t, -1, digits.At(50))
	checkFullIter(t, digits.Items(), 15, 5, 25, 7)
	checkFullIter(t, digits.ReverseItems(), 25, 7, 15, 5)
}

func TestDigitsNone(t *testing.T) {
	var positions Positions
	assert.Zero(t, GetDigits(Sqrt(2), positions))
}

func TestDigitsNoneZeroNumber(t *testing.T) {
	var n Number
	var p Positions
	var pb PositionsBuilder
	assert.Zero(t, GetDigits(&n, p))
	assert.Zero(t, GetDigits(&n, pb.AddRange(0, 100).Build()))
}

func TestDigitsNoneFromDigits(t *testing.T) {
	digits := AllDigits(Sqrt(2).WithSignificant(1000))
	assert.Equal(t, 1000, digits.Len())
	var zeroPosits Positions
	zeroDigits := GetDigits(digits, zeroPosits)
	assert.Zero(t, zeroDigits)
	assert.Zero(t, GetDigits(zeroDigits, zeroPosits))
	var pb PositionsBuilder
	assert.Zero(t, GetDigits(zeroDigits, pb.AddRange(0, 100).Build()))
}

func TestDigitsZero(t *testing.T) {
	var digits Digits
	iter := digits.Items()
	_, ok := iter()
	assert.False(t, ok)
	iter = digits.ReverseItems()
	_, ok = iter()
	assert.False(t, ok)
	assert.Equal(t, -1, digits.At(0))
	assert.Empty(t, digits.Sprint())
	assert.Equal(t, -1, digits.Min())
	assert.Equal(t, -1, digits.Max())
	assert.Equal(t, 0, digits.Len())
	assert.Empty(t, FindAll(digits, nil))
	assert.Empty(t, FindAll(digits, []int{5}))
	assert.Equal(t, 0, digits.WithStart(5).Len())
	assert.Equal(t, 0, digits.WithEnd(5).Len())
}

func TestDigitsBinary(t *testing.T) {
	n := Sqrt(2)
	var pb PositionsBuilder
	pb.AddRange(1000, 2000).AddRange(5000, 5999).AddRange(10000, 10999).Add(11000)
	digits := GetDigits(n, pb.Build())
	arr, err := digits.MarshalBinary()
	assert.Equal(t, 1509, len(arr))
	assert.NoError(t, err)
	var copy Digits
	assert.NoError(t, copy.UnmarshalBinary(arr))
	assert.Equal(t, digits.Sprint(), copy.Sprint())
}

func TestDigitsBinary2(t *testing.T) {
	n := Sqrt(2)
	var pb PositionsBuilder
	pb.Add(50).Add(25).Add(15).Add(0).AddRange(100, 102)
	digits := GetDigits(n, pb.Build())
	arr, err := digits.MarshalBinary()
	assert.NoError(t, err)
	var copy Digits
	assert.NoError(t, copy.UnmarshalBinary(arr))
	assert.Equal(t, digits.Sprint(), copy.Sprint())
}

func TestDigitsBinary3(t *testing.T) {
	n := Sqrt(2)
	var pb PositionsBuilder
	pb.AddRange(0, 4).Add(26)
	digits := GetDigits(n, pb.Build())
	arr, err := digits.MarshalBinary()
	assert.NoError(t, err)
	assert.Equal(t, 6, len(arr))
	var copy Digits
	assert.NoError(t, copy.UnmarshalBinary(arr))
	assert.Equal(t, digits.Sprint(), copy.Sprint())
}

func TestDigitsBinaryZero(t *testing.T) {
	var digits Digits
	arr, err := digits.MarshalBinary()
	assert.NoError(t, err)
	assert.Equal(t, 1, len(arr))
	var copy Digits
	assert.NoError(t, copy.UnmarshalBinary(arr))
	assert.Zero(t, copy)
	assert.Equal(t, digits.Sprint(), copy.Sprint())
}

func TestDigitsBinaryEmpty(t *testing.T) {
	var digits Digits
	assert.Error(t, digits.UnmarshalBinary(nil))
}

func TestDigitsBinaryBadVersion(t *testing.T) {
	var digits Digits
	assert.Error(t, digits.UnmarshalBinary([]byte{51}))
}

func TestDigitsBinaryUnmarshalErrors(t *testing.T) {
	// A very large position gap results in a negative position delta.
	data := []byte{0xbb, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x01, 0x44}
	var digits Digits
	assert.Error(t, digits.UnmarshalBinary(data))
	data = []byte{0xbb, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x01, 0x6d}
	assert.Error(t, digits.UnmarshalBinary(data))
	data = []byte{0xbb, 0xff}
	assert.Error(t, digits.UnmarshalBinary(data))
}

func TestDigitsText(t *testing.T) {
	n := Sqrt(2)
	var pb PositionsBuilder
	pb.AddRange(0, 4).Add(26)
	digits := GetDigits(n, pb.Build())
	arr, err := digits.MarshalText()
	assert.NoError(t, err)
	assert.Equal(t, "v1:1414[26]2", string(arr))
	var copy Digits
	assert.NoError(t, copy.UnmarshalText(arr))
	assert.Equal(t, digits.Sprint(), copy.Sprint())
}

func TestDigitsText2(t *testing.T) {
	n := Sqrt(2)
	var pb PositionsBuilder
	pb.AddRange(3, 6).AddRange(10, 12)
	digits := GetDigits(n, pb.Build())
	arr, err := digits.MarshalText()
	assert.NoError(t, err)
	assert.Equal(t, "v1:[3]421[10]37", string(arr))
	var copy Digits
	assert.NoError(t, copy.UnmarshalText(arr))
	assert.Equal(t, digits.Sprint(), copy.Sprint())
}

func TestDigitsText3(t *testing.T) {
	digits := AllDigits(Sqrt(2).WithSignificant(6))
	arr, err := digits.MarshalText()
	assert.NoError(t, err)
	assert.Equal(t, "v1:141421", string(arr))
	var copy Digits
	assert.NoError(t, copy.UnmarshalText(arr))
	assert.Equal(t, digits.Sprint(), copy.Sprint())
}

func TestDigitsTextZero(t *testing.T) {
	var digits Digits
	arr, err := digits.MarshalText()
	assert.NoError(t, err)
	assert.Equal(t, "v1:", string(arr))
	var copy Digits
	assert.NoError(t, copy.UnmarshalText(arr))
	assert.Zero(t, copy)
	assert.Equal(t, digits.Sprint(), copy.Sprint())
}

func TestDigitsTextUnmarshalErrors(t *testing.T) {
	text := []byte("v1:12345[")
	var digits Digits
	assert.Error(t, digits.UnmarshalText(text))
	text = []byte("v1:12345[67]")
	assert.Error(t, digits.UnmarshalText(text))
	text = []byte("v1:12abc")
	assert.Error(t, digits.UnmarshalText(text))
	text = []byte("v1:12345[6a]")
	assert.Error(t, digits.UnmarshalText(text))
	text = []byte("v2:")
	assert.Error(t, digits.UnmarshalText(text))
	text = []byte("")
	assert.Error(t, digits.UnmarshalText(text))
}

func TestDigitsWithStartAndEnd(t *testing.T) {
	digits := AllDigits(Sqrt(2).WithSignificant(1000))
	assert.NotEqual(t, -1, digits.At(700))
	assert.NotEqual(t, -1, digits.At(399))
	digits = digits.WithStart(200).WithEnd(900)
	digits = digits.WithStart(400).WithEnd(700).WithStart(300).WithEnd(800)
	assert.Equal(t, 300, digits.Len())
	assert.Equal(t, 400, digits.Min())
	assert.Equal(t, 699, digits.Max())
	assert.Equal(t, -1, digits.At(700))
	assert.Equal(t, -1, digits.At(399))
	assert.NotEqual(t, -1, digits.At(400))
	assert.NotEqual(t, -1, digits.At(699))
	assert.Len(t, digits.Sprint(), 389)
}

func TestDigitsWithStartZero(t *testing.T) {
	digits := AllDigits(Sqrt(2).WithSignificant(1000))
	digits = digits.WithStart(1000)
	assert.Equal(t, 0, digits.Len())
	assert.Equal(t, -1, digits.Min())
	assert.Equal(t, -1, digits.Max())
	assert.Equal(t, -1, digits.At(500))
}

func TestDigitsWithEndZero(t *testing.T) {
	digits := AllDigits(Sqrt(2).WithSignificant(1000))
	digits = digits.WithEnd(0)
	assert.Equal(t, 0, digits.Len())
	assert.Equal(t, -1, digits.Min())
	assert.Equal(t, -1, digits.Max())
	assert.Equal(t, -1, digits.At(500))
}

func TestDigitsBuilder(t *testing.T) {
	var builder digitsBuilder
	assert.NoError(t, builder.AddDigit(0, 3))
	assert.NoError(t, builder.AddDigit(1, 1))
	assert.NoError(t, builder.AddDigit(2, 4))
	assert.NoError(t, builder.AddDigit(3, 1))
	assert.NoError(t, builder.AddDigit(6, 4))
	assert.NoError(t, builder.AddDigit(7, 1))
	assert.NoError(t, builder.AddDigit(8, 4))
	digits := builder.Build()
	assert.NoError(t, builder.AddDigit(0, 4))
	assert.NoError(t, builder.AddDigit(1, 6))
	assert.NoError(t, builder.AddDigit(2, 6))
	assert.NoError(t, builder.AddDigit(3, 9))
	assert.NoError(t, builder.AddDigit(4, 2))
	assert.NoError(t, builder.AddDigit(5, 0))
	feigenbaum := builder.Build()
	actual := digits.Sprint(
		DigitsPerRow(0), DigitsPerColumn(0), ShowCount(false))
	assert.Equal(t, "0.3141..414", actual)
	assert.Equal(t, 7, digits.Len())
	matches := FindAll(digits, []int{1, 4})
	assert.Equal(t, []int{1, 7}, matches)
	assert.Equal(t, 6, feigenbaum.Len())
}

func TestDigitBuilderErrors(t *testing.T) {
	var builder digitsBuilder
	assert.Error(t, builder.AddDigit(-1, 3))
	assert.Error(t, builder.AddDigit(0, -1))
	assert.Error(t, builder.AddDigit(0, 10))
	assert.NoError(t, builder.AddDigit(0, 2))
	assert.Error(t, builder.AddDigit(0, 1))
	digits := builder.Build()
	assert.Equal(t, 1, digits.Len())
}

func TestDigitLookup(t *testing.T) {
	n := Sqrt(11).WithSignificant(10000)
	pattern := []int{4, 5, 7}
	finds := FindAll(n, pattern)
	assert.Len(t, finds, 7)
	var pb PositionsBuilder
	for _, find := range finds {
		pb.AddRange(find, find+3)
	}
	digits := GetDigits(n, pb.Build())
	iter := digits.Items()
	count := 0
	for digit, ok := iter(); ok; digit, ok = iter() {
		assert.Equal(t, pattern[count%3], digit.Value)
		assert.Equal(t, finds[count/3]+count%3, digit.Position)
		count++
	}
	assert.Equal(t, 7*3, count)
}

func TestFindWithDigits(t *testing.T) {
	n := Sqrt(5)
	// '000' in Sqrt(5) at 424 569 3663 4506 4601 6113 7173 9110 9114
	var pb PositionsBuilder
	pb.AddRange(7173, 7176).AddRange(9110, 9117).AddRange(4500, 4600)
	digits := GetDigits(n, pb.Build())
	finds := FindAll(digits, []int{0, 0, 0})
	assert.Equal(t, []int{4506, 7173, 9110, 9114}, finds)
}

func checkFullIter(t *testing.T, iter func() (Digit, bool), positionDigits ...int) {
	t.Helper()
	var actual []int
	for digit, ok := iter(); ok; digit, ok = iter() {
		actual = append(actual, digit.Position, digit.Value)
	}
	assert.Equal(t, positionDigits, actual)
	_, ok := iter()
	assert.False(t, ok)
}

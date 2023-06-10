package sqroot

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/keep94/consume2"
)

const (
	digitsBinaryVersion        = 187
	unmarshalTextUnexpectedEnd = "sqroot: Digits.UnmarshalText hit unexpected end of text"
)

// Digits holds the digits found at selected positions of a Number so
// that they can be quickly retrieved. With Memoization available on
// Number instances, Digits instances are needed only in rare cases.
// Retrieving any digit by position takes O(log N) time where N is the
// total number of digits. Using Items() or ReverseItems() to retrieve
// all the digits in order takes O(N) time. The zero value is no digits.
// Digits implements Sequence.
type Digits struct {
	digits []Digit
}

// GetDigits gets the digits from s found at the zero based positions
// in p.
func GetDigits(s Sequence, p Positions) Digits {
	var builder digitsBuilder
	fromSequenceWithPositions(s, p, &builder)
	return builder.Build()
}

// AllDigits gets all digits from s. If s has an infinite number of digits,
// AllDigits will run forever.
func AllDigits(s Sequence) Digits {
	var builder digitsBuilder
	consume2.FromGenerator[Digit](s.digitIter(), &builder)
	return builder.Build()
}

// WithStart returns a Digits like this one that only has positions greater
// than or equal to start. Returned instance shares memory with this instance.
// Therefore, to change only the starting position it is more efficient to use
// WithStart than GetDigits. WithStart runs in O(log N) time where N is the
// number of digits in this instance.
func (d Digits) WithStart(start int) Digits {
	index := sort.Search(
		len(d.digits),
		func(x int) bool { return d.digits[x].Position >= start },
	)
	return Digits{digits: d.digits[index:]}
}

// WithEnd returns a Digits like this one that only has positions less than
// end. Returned instance shares memory with this instance. Therefore, to
// change only the ending position it is more efficient to use WithEnd than
// GetDigits. WithEnd runs in O(log N) time where N is the number of digits
// in this instance.
func (d Digits) WithEnd(end int) Digits {
	index := sort.Search(
		len(d.digits),
		func(x int) bool { return d.digits[x].Position >= end },
	)
	return Digits{digits: d.digits[:index]}
}

// MarshalBinary implements the encoding.BinaryMarshaler interface.
func (d Digits) MarshalBinary() ([]byte, error) {
	iter := d.Items()
	nextPosit := 0
	result := []byte{digitsBinaryVersion}
	state := 0
	pair := uint64(0)
	for digit, ok := iter(); ok; digit, ok = iter() {
		delta := digit.Position - nextPosit
		if delta > 0 {
			if state == 1 {
				result = binary.AppendUvarint(result, 100+pair)
				state = 0
				pair = 0
			}
			result = binary.AppendUvarint(result, uint64(delta)+109)
		}
		nextPosit = digit.Position + 1
		pair = 10*pair + uint64(digit.Value)
		if state == 1 {
			result = binary.AppendUvarint(result, pair)
			pair = 0
		}
		state = 1 - state
	}
	if state == 1 {
		result = binary.AppendUvarint(result, 100+pair)
	}
	return result, nil
}

// MarshalText implements the encoding.TextMarshaler interface.
func (d Digits) MarshalText() ([]byte, error) {
	iter := d.Items()
	nextPosit := 0
	result := []byte("v1:")
	for digit, ok := iter(); ok; digit, ok = iter() {
		if digit.Position > nextPosit {
			result = append(result, byte('['))
			result = strconv.AppendInt(result, int64(digit.Position), 10)
			result = append(result, byte(']'))
			nextPosit = digit.Position
		}
		result = append(result, byte('0'+digit.Value))
		nextPosit++
	}
	return result, nil
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface.
func (d *Digits) UnmarshalBinary(b []byte) error {
	if len(b) == 0 || b[0] != digitsBinaryVersion {
		return errors.New("sqroot: Bad Digits Binary Version")
	}
	var builder digitsBuilder
	posit := 0
	reader := bytes.NewReader(b[1:])
	for reader.Len() > 0 {
		val, err := binary.ReadUvarint(reader)
		if err != nil {
			return err
		}
		if val >= 110 {
			posit += int(val - 109)
		} else if val >= 100 {
			if err := builder.AddDigit(posit, int(val-100)); err != nil {
				return err
			}
			posit++
		} else {
			if err := builder.AddDigit(posit, int(val/10)); err != nil {
				return err
			}
			posit++
			if err := builder.AddDigit(posit, int(val%10)); err != nil {
				return err
			}
			posit++
		}
	}
	*d = builder.Build()
	return nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (d *Digits) UnmarshalText(text []byte) error {
	var builder digitsBuilder
	posit := 0
	version, i, err := readVersion(text)
	if err != nil || version != "v1" {
		return errors.New("sqroot: Bad Digits Text Version")
	}
	for i < len(text) {
		if text[i] == '[' {
			posit, i, err = readPositiveInt(text, i+1)
			if err != nil {
				return err
			}
		}
		digit := int(text[i] - '0')
		if err := builder.AddDigit(posit, digit); err != nil {
			return err
		}
		posit++
		i++
	}
	*d = builder.Build()
	return nil
}

// At returns the digit between 0 and 9 at the given zero based position.
// If the digit at posit is unknown or if posit is negative, At returns -1.
// At runs in O(log N) time where N is the number of digits in this instance.
func (d Digits) At(posit int) int {
	index := sort.Search(
		len(d.digits),
		func(x int) bool { return d.digits[x].Position >= posit },
	)
	if index == len(d.digits) || d.digits[index].Position != posit {
		return -1
	}
	return d.digits[index].Value
}

// Items returns a function that generates all the digits in this instance
// from lowest to highest position. When there are no more digits,
// the returned function returns false.
func (d Digits) Items() func() (digit Digit, ok bool) {
	return d.digitIter()
}

// ReverseItems returns a function that generates all the digits in this
// instance from highest to lowest position. When there are no more digits,
// the returned function returns false.
func (d Digits) ReverseItems() func() (digit Digit, ok bool) {
	return d.reverseDigitIter()
}

// Min returns the minimum position in this instance. If this instance
// is empty, Min returns -1.
func (d Digits) Min() int {
	if len(d.digits) == 0 {
		return -1
	}
	return d.digits[0].Position
}

// Max returns the maximum position in this instance. If this instance
// is empty, Max returns -1.
func (d Digits) Max() int {
	if len(d.digits) == 0 {
		return -1
	}
	return d.digits[len(d.digits)-1].Position
}

// Len returns the number of digits in this instance.
func (d Digits) Len() int {
	return len(d.digits)
}

// Print works like Fprint printing this instance to stdout.
func (d Digits) Print(options ...Option) (n int, err error) {
	return d.Fprint(os.Stdout, options...)
}

// Sprint works like Fprint printing this instance to the returned string.
func (d Digits) Sprint(options ...Option) string {
	var builder strings.Builder
	d.Fprint(&builder, options...)
	return builder.String()
}

// Fprint prints this instance to w and returns number of bytes written
// and any error encountered. For options, the default is 50 digits per
// row, 5 digits per column, show digit count, and period (.) for missing
// digits.
func (d Digits) Fprint(w io.Writer, options ...Option) (n int, err error) {
	settings := &printerSettings{
		digitsPerRow:    50,
		digitsPerColumn: 5,
		showCount:       true,
		missingDigit:    '.',
	}
	p := newPrinter(w, d.limit(), mutateSettings(options, settings))
	consume2.FromGenerator[Digit](d.digitIter(), p)
	return p.byteCount, p.err
}

func (d Digits) limit() int {
	return d.Max() + 1
}

func (d Digits) digitIter() func() (Digit, bool) {
	index := 0
	return func() (digit Digit, ok bool) {
		if index == len(d.digits) {
			return
		}
		result := d.digits[index]
		index++
		return result, true
	}
}

func (d Digits) canReverse() bool {
	return true
}

func (d Digits) reverseDigitIter() func() (Digit, bool) {
	index := len(d.digits)
	return func() (digit Digit, ok bool) {
		if index == 0 {
			return
		}
		index--
		return d.digits[index], true
	}
}

func (d Digits) enabled() bool {
	return true
}

func (d Digits) get(start, end int) Sequence {
	return d.WithStart(start).WithEnd(end)
}

func readVersion(text []byte) (string, int, error) {
	idx := bytes.Index(text, []byte(":"))
	if idx == -1 {
		return "", 0, errors.New("sqroot: Digits.UnmarhalText: Can't read version")
	}
	return string(text[:idx]), idx + 1, nil
}

func readPositiveInt(text []byte, i int) (int, int, error) {
	result := 0
	for i < len(text) {
		if text[i] == ']' {
			if i+1 == len(text) {
				return 0, 0, errors.New(unmarshalTextUnexpectedEnd)
			}
			return result, i + 1, nil
		} else if text[i] >= '0' && text[i] <= '9' {
			result = result*10 + int(text[i]-'0')
		} else {
			return 0, 0, fmt.Errorf("sqroot: Digits.UnmarshalText unexpected character in text: %c", text[i])
		}
		i++
	}
	return 0, 0, errors.New(unmarshalTextUnexpectedEnd)
}

type digitsBuilder struct {
	digits []Digit
}

func (d *digitsBuilder) CanConsume() bool {
	return true
}

func (d *digitsBuilder) Consume(digit Digit) {
	d.digits = append(d.digits, digit)
}

func (d *digitsBuilder) AddDigit(posit int, digit int) error {
	if posit < 0 {
		return fmt.Errorf(
			"sqroot: digitsBuilder.AddDigit: posit must be non-negative but was %d", posit)
	}
	if digit < 0 || digit > 9 {
		return fmt.Errorf(
			"sqroot: digitsBuilder.AddDigit: digit must be between 0 and 9 but was %d", digit)
	}
	if len(d.digits) > 0 && d.digits[len(d.digits)-1].Position >= posit {
		return fmt.Errorf(
			"sqroot: digitsBuilder.AddDigit: posit must be ever increasing was %d now %d",
			d.digits[len(d.digits)-1].Position,
			posit,
		)
	}
	d.digits = append(d.digits, Digit{Position: posit, Value: digit})
	return nil
}

func (d *digitsBuilder) Build() Digits {
	result := Digits{digits: d.digits}
	*d = digitsBuilder{}
	return result
}

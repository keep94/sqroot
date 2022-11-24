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
)

const (
	digitsBinaryVersion        = 187
	unmarshalTextUnexpectedEnd = "sqroot: Digits.UnmarshalText hit unexpected end of text"
)

// Digits holds the digits found at selected positions of a Mantissa so
// that they can be quickly retrieved. The zero value is no digits.
// Digits implements Sequence.
type Digits struct {
	digits map[int]int
	posits []int
}

// GetDigits gets the digits from s found at the zero based positions
// in p.
func GetDigits(s Sequence, p Positions) Digits {
	d, ok := s.(Digits)
	if ok && p.count < d.Len()+len(p.ranges) {

		// Optimization: Just choose what we want instead of iterating
		// over all of s.
		return d.pick(p)
	}
	var builder digitsBuilder
	sendPositDigits(newPart(s, p), &builder)
	return builder.Build()
}

// MarshalBinary implements the encoding.BinaryMarshaler interface.
func (d Digits) MarshalBinary() ([]byte, error) {
	iter := d.Iterator()
	nextPosit := 0
	result := []byte{digitsBinaryVersion}
	state := 0
	pair := uint64(0)
	for posit := iter(); posit != -1; posit = iter() {
		delta := posit - nextPosit
		if delta > 0 {
			if state == 1 {
				result = binary.AppendUvarint(result, 100+pair)
				state = 0
				pair = 0
			}
			result = binary.AppendUvarint(result, uint64(delta)+109)
		}
		nextPosit = posit + 1
		pair = 10*pair + uint64(d.At(posit))
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
	iter := d.Iterator()
	nextPosit := 0
	result := []byte("v1:")
	for posit := iter(); posit != -1; posit = iter() {
		if posit > nextPosit {
			result = append(result, byte('['))
			result = strconv.AppendInt(result, int64(posit), 10)
			result = append(result, byte(']'))
			nextPosit = posit
		}
		result = append(result, byte('0'+d.At(posit)))
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
func (d Digits) At(posit int) int {
	digit, ok := d.digits[posit]
	if !ok {
		return -1
	}
	return digit
}

// Iterator returns a function that generates all the zero based positions
// in this instance from lowest to highest. When there are no more positions,
// the returned function returns -1.
func (d Digits) Iterator() func() int {
	return d.iteratorAtIndex(0)
}

// Reverse returns a function that generates all the zero based positions
// in this instance from highest to lowest. When there are no more positions,
// the returned function returns -1.
func (d Digits) Reverse() func() int {
	return d.reverseAtIndex(len(d.posits))
}

// IteratorAt returns a function that generates all the zero based
// positions in this instance from lowest to highest starting at posit.
// When there are no more positions, the returned function returns -1.
func (d Digits) IteratorAt(posit int) func() int {
	index := sort.Search(
		len(d.posits), func(x int) bool { return d.posits[x] >= posit })
	return d.iteratorAtIndex(index)
}

// ReverseAt returns a function that generates all the zero based positions
// in this instance from highest to lowest that come before posit. When
// there are no more positions, the returned function returns -1.
func (d Digits) ReverseAt(posit int) func() int {
	index := sort.Search(
		len(d.posits), func(x int) bool { return d.posits[x] >= posit })
	return d.reverseAtIndex(index)
}

// Min returns the minimum position in this instance. If this instance
// is empty, Min returns -1.
func (d Digits) Min() int {
	if len(d.posits) == 0 {
		return -1
	}
	return d.posits[0]
}

// Max returns the maximum position in this instance. If this instance
// is empty, Max returns -1.
func (d Digits) Max() int {
	if len(d.posits) == 0 {
		return -1
	}
	return d.posits[len(d.posits)-1]
}

// Len returns the number of digits in this instance.
func (d Digits) Len() int {
	return len(d.posits)
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
	return fprint(w, d, mutateSettings(options, settings))
}

func (d Digits) limit() int {
	if len(d.posits) == 0 {
		return 0
	}
	return d.posits[len(d.posits)-1] + 1
}

func (d Digits) iteratorAtIndex(index int) func() int {
	return func() int {
		if index == len(d.posits) {
			return -1
		}
		result := d.posits[index]
		index++
		return result
	}
}

func (d Digits) reverseAtIndex(index int) func() int {
	index--
	return func() int {
		if index == -1 {
			return -1
		}
		result := d.posits[index]
		index--
		return result
	}
}

func (d Digits) positDigitIter() func() positDigit {
	index := 0
	return func() positDigit {
		if index == len(d.posits) {
			return invalidPositDigit
		}
		posit := d.posits[index]
		result := positDigit{Posit: posit, Digit: d.digits[posit]}
		index++
		return result
	}
}

func (d Digits) reversePositDigitIter() func() positDigit {
	index := len(d.posits)
	return func() positDigit {
		if index == 0 {
			return invalidPositDigit
		}
		index--
		posit := d.posits[index]
		return positDigit{Posit: posit, Digit: d.digits[posit]}
	}
}

func (d Digits) rfind(pattern []int) func() int {
	if len(pattern) == 0 {
		return zeroPattern(d.reversePositDigitIter())
	}
	return kmp(d.reversePositDigitIter(), patternReverse(pattern), true)
}

func (d Digits) findLastN(pattern []int, n int) []int {
	var result []int
	iter := d.rfind(pattern)
	for index := iter(); index != -1 && len(result) < n; index = iter() {
		result = append(result, index)
	}
	return result
}

func (d Digits) pick(p Positions) Digits {
	var builder digitsBuilder
	for _, pr := range p.ranges {
		for posit := pr.Start; posit < pr.End; posit++ {
			digit := d.At(posit)
			if digit != -1 {
				builder.unsafeAddDigit(posit, digit)
			}
		}
	}
	return builder.Build()
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
	digits map[int]int
	posits []int
}

func (d *digitsBuilder) CanConsume() bool {
	return true
}

func (d *digitsBuilder) Consume(pd positDigit) {
	d.unsafeAddDigit(pd.Posit, pd.Digit)
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
	if len(d.posits) > 0 && d.posits[len(d.posits)-1] >= posit {
		return fmt.Errorf(
			"sqroot: digitsBuilder.AddDigit: posit must be ever increasing was %d now %d",
			d.posits[len(d.posits)-1],
			posit,
		)
	}
	d.unsafeAddDigit(posit, digit)
	return nil
}

func (d *digitsBuilder) Build() Digits {
	result := Digits{digits: d.digits, posits: d.posits}
	*d = digitsBuilder{}
	return result
}

func (d *digitsBuilder) unsafeAddDigit(posit int, digit int) {
	if d.digits == nil {
		d.digits = make(map[int]int)
	}
	d.digits[posit] = digit
	d.posits = append(d.posits, posit)
}

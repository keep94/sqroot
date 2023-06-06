package sqroot

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNumberZeroValueString(t *testing.T) {
	var number Number
	assert.Equal(t, "0", number.String())
}

func TestNumberFPositiveExponent(t *testing.T) {
	number := fakeNumber.withExponent(5)
	actual := fmt.Sprintf("%f", number)
	assert.Equal(t, "12345.678901", actual)
	actual = fmt.Sprintf("%.1f", number)
	assert.Equal(t, "12345.6", actual)
	actual = fmt.Sprintf("%.0f", number)
	assert.Equal(t, "12345", actual)
}

func TestNumberFPositiveExponentFiniteDigits(t *testing.T) {
	number := fakeNumber.WithSignificant(9).withExponent(5)
	actual := fmt.Sprintf("%F", number)
	assert.Equal(t, "12345.678900", actual)
}

func TestNumberFNegExponent(t *testing.T) {
	number := fakeNumber.withExponent(-5)
	actual := fmt.Sprintf("%f", number)
	assert.Equal(t, "0.000001", actual)
	actual = fmt.Sprintf("%.10f", number)
	assert.Equal(t, "0.0000012345", actual)
	actual = fmt.Sprintf("%.5f", number)
	assert.Equal(t, "0.00000", actual)
	actual = fmt.Sprintf("%.1f", number)
	assert.Equal(t, "0.0", actual)
	actual = fmt.Sprintf("%.0f", number)
	assert.Equal(t, "0", actual)
}

func TestNumberFZeroExponent(t *testing.T) {
	number := fakeNumber.withExponent(0)
	actual := fmt.Sprintf("%f", number)
	assert.Equal(t, "0.123456", actual)
	actual = fmt.Sprintf("%.10f", number)
	assert.Equal(t, "0.1234567890", actual)
	actual = fmt.Sprintf("%.5f", number)
	assert.Equal(t, "0.12345", actual)
	actual = fmt.Sprintf("%.1f", number)
	assert.Equal(t, "0.1", actual)
	actual = fmt.Sprintf("%.0f", number)
	assert.Equal(t, "0", actual)
}

func TestNumberFZero(t *testing.T) {
	var number Number
	actual := fmt.Sprintf("%f", &number)
	assert.Equal(t, "0.000000", actual)
	actual = fmt.Sprintf("%.3f", &number)
	assert.Equal(t, "0.000", actual)
	actual = fmt.Sprintf("%.1f", &number)
	assert.Equal(t, "0.0", actual)
	actual = fmt.Sprintf("%.0f", &number)
	assert.Equal(t, "0", actual)
}

func TestNumberGPositiveExponent(t *testing.T) {
	number := fakeNumber.withExponent(5)
	actual := fmt.Sprintf("%g", number)
	assert.Equal(t, "12345.67890123456", actual)
	actual = fmt.Sprintf("%.8g", number)
	assert.Equal(t, "12345.678", actual)
	actual = fmt.Sprintf("%.5g", number)
	assert.Equal(t, "12345", actual)
	actual = fmt.Sprintf("%.4g", number)
	assert.Equal(t, "0.1234e+05", actual)
	actual = fmt.Sprintf("%.0G", number)
	assert.Equal(t, "0.1E+05", actual)
}

func TestNumberGPositiveExponentShort(t *testing.T) {
	number := fakeNumber.WithSignificant(3).withExponent(5)
	actual := fmt.Sprintf("%g", number)
	assert.Equal(t, "12300", actual)
	actual = fmt.Sprintf("%.5g", number)
	assert.Equal(t, "12300", actual)
	actual = fmt.Sprintf("%.4g", number)
	assert.Equal(t, "0.123e+05", actual)
}

func TestNumberGPositiveExponentFiniteDigits(t *testing.T) {
	number := fakeNumber.WithSignificant(9).withExponent(5)
	actual := fmt.Sprintf("%G", number)
	assert.Equal(t, "12345.6789", actual)
}

func TestNumberGZeroExponent(t *testing.T) {
	number := fakeNumber.withExponent(0)
	actual := fmt.Sprintf("%g", number)
	assert.Equal(t, "0.1234567890123456", actual)
	actual = fmt.Sprintf("%.8g", number)
	assert.Equal(t, "0.12345678", actual)
	actual = fmt.Sprintf("%.0g", number)
	assert.Equal(t, "0.1", actual)
}

func TestNumberGNegExponent(t *testing.T) {
	number := fakeNumber.withExponent(-3)
	actual := fmt.Sprintf("%g", number)
	assert.Equal(t, "0.0001234567890123456", actual)
	actual = fmt.Sprintf("%.8g", number)
	assert.Equal(t, "0.00012345678", actual)
	actual = fmt.Sprintf("%.0g", number)
	assert.Equal(t, "0.0001", actual)
}

func TestNumberGZero(t *testing.T) {
	var number Number
	actual := fmt.Sprintf("%G", &number)
	assert.Equal(t, "0", actual)
	actual = fmt.Sprintf("%.0g", &number)
	assert.Equal(t, "0", actual)
}

func TestNumberGLargePosExponent(t *testing.T) {
	number := fakeNumber.withExponent(7)
	actual := fmt.Sprintf("%G", number)
	assert.Equal(t, "0.1234567890123456E+07", actual)
	actual = fmt.Sprintf("%.8g", number)
	assert.Equal(t, "0.12345678e+07", actual)
	actual = fmt.Sprintf("%.0g", number)
	assert.Equal(t, "0.1e+07", actual)
	number = fakeNumber.withExponent(6)
	actual = fmt.Sprintf("%.6g", number)
	assert.Equal(t, "123456", actual)
	number = fakeNumber.withExponent(10)
	actual = fmt.Sprintf("%.10g", number)
	assert.Equal(t, "0.1234567890e+10", actual)
}

func TestNumberGLargePosExponentFiniteDigits(t *testing.T) {
	number := fakeNumber.WithSignificant(9).withExponent(7)
	actual := fmt.Sprintf("%g", number)
	assert.Equal(t, "0.123456789e+07", actual)
}

func TestNumberGLargeNegExponent(t *testing.T) {
	number := fakeNumber.withExponent(-4)
	actual := fmt.Sprintf("%G", number)
	assert.Equal(t, "0.1234567890123456E-04", actual)
}

func TestNumberEPositiveExponent(t *testing.T) {
	number := fakeNumber.withExponent(5)
	actual := fmt.Sprintf("%e", number)
	assert.Equal(t, "0.123456e+05", actual)
	actual = fmt.Sprintf("%.1E", number)
	assert.Equal(t, "0.1E+05", actual)
	actual = fmt.Sprintf("%.0e", number)
	assert.Equal(t, "0e+05", actual)
}

func TestNumberEPositiveExponentFiniteDigits(t *testing.T) {
	number := fakeNumber.WithSignificant(9).withExponent(5)
	actual := fmt.Sprintf("%.14e", number)
	assert.Equal(t, "0.12345678900000e+05", actual)
}

func TestNumberEZeroExponent(t *testing.T) {
	number := fakeNumber.withExponent(0)
	actual := fmt.Sprintf("%e", number)
	assert.Equal(t, "0.123456e+00", actual)
	actual = fmt.Sprintf("%.1E", number)
	assert.Equal(t, "0.1E+00", actual)
	actual = fmt.Sprintf("%.0e", number)
	assert.Equal(t, "0e+00", actual)
}

func TestNumberENegExponent(t *testing.T) {
	number := fakeNumber.withExponent(-5)
	actual := fmt.Sprintf("%e", number)
	assert.Equal(t, "0.123456e-05", actual)
	actual = fmt.Sprintf("%.1E", number)
	assert.Equal(t, "0.1E-05", actual)
	actual = fmt.Sprintf("%.0e", number)
	assert.Equal(t, "0e-05", actual)
}

func TestNumberEZero(t *testing.T) {
	var number Number
	actual := fmt.Sprintf("%E", &number)
	assert.Equal(t, "0.000000E+00", actual)
	actual = fmt.Sprintf("%.1e", &number)
	assert.Equal(t, "0.0e+00", actual)
	actual = fmt.Sprintf("%.0e", &number)
	assert.Equal(t, "0e+00", actual)
}

func TestNumberWidth(t *testing.T) {
	number := fakeNumber.withExponent(5)
	actual := fmt.Sprintf("%20v", number)
	assert.Equal(t, "   12345.67890123456", actual)
	actual = fmt.Sprintf("%16v", number)
	assert.Equal(t, "12345.67890123456", actual)
	actual = fmt.Sprintf("%-20v", number)
	assert.Equal(t, "12345.67890123456   ", actual)
	actual = fmt.Sprintf("%-16v", number)
	assert.Equal(t, "12345.67890123456", actual)
	actual = fmt.Sprintf("%6.5v", number)
	assert.Equal(t, " 12345", actual)
}

func TestNumberString(t *testing.T) {
	number := fakeNumber.WithSignificant(9).withExponent(6)
	assert.Equal(t, "123456.789", number.String())
	number = fakeNumber.withExponent(6)
	assert.Equal(t, "123456.7890123456", number.String())
	number = fakeNumber.withExponent(7)
	assert.Equal(t, "0.1234567890123456e+07", number.String())
	number = fakeNumber.withExponent(11)
	assert.Equal(t, "0.1234567890123456e+11", number.String())
	number = fakeNumber.withExponent(-3)
	assert.Equal(t, "0.0001234567890123456", number.String())
	number = fakeNumber.withExponent(-4)
	assert.Equal(t, "0.1234567890123456e-04", number.String())
	number = &Number{}
	assert.Equal(t, "0", number.String())
}

func TestNumberBadVerb(t *testing.T) {
	number := fakeNumber.WithSignificant(9).withExponent(5)
	actual := fmt.Sprintf("%h", number)
	assert.Equal(t, "%!h(number=12345.6789)", actual)
}

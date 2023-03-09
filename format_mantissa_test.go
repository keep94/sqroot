package sqroot

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormat(t *testing.T) {
	actual := fmt.Sprintf("%.14f", fakeMantissa)
	assert.Equal(t, "0.12345678901234", actual)
}

func TestFormatNoPrecision(t *testing.T) {
	actual := fmt.Sprintf("%f", fakeMantissa)
	assert.Equal(t, "0.123456", actual)
}

func TestFormatNoPrecisionCapital(t *testing.T) {
	actual := fmt.Sprintf("%F", fakeMantissa)
	assert.Equal(t, "0.123456", actual)
}

func TestFormatNotInfinite(t *testing.T) {
	actual := fmt.Sprintf("%.14f", fakeMantissa.WithSignificant(9))
	assert.Equal(t, "0.12345678900000", actual)
}

func TestFormatNotInfiniteNoPrecision(t *testing.T) {
	actual := fmt.Sprintf("%f", fakeMantissa.WithSignificant(3))
	assert.Equal(t, "0.123000", actual)
}

func TestFormatWidth(t *testing.T) {
	actual := fmt.Sprintf("%10f", fakeMantissa)
	assert.Equal(t, "  0.123456", actual)
}

func TestFormatShortWidth(t *testing.T) {
	actual := fmt.Sprintf("%5f", fakeMantissa)
	assert.Equal(t, "0.123456", actual)
}

func TestFormatWidthLeftJustify(t *testing.T) {
	actual := fmt.Sprintf("%-10f", fakeMantissa)
	assert.Equal(t, "0.123456  ", actual)
}

func TestFormatWidthAndPrecision(t *testing.T) {
	actual := fmt.Sprintf("%-20.13f", fakeMantissa)
	assert.Equal(t, "0.1234567890123     ", actual)
}

func TestFormatPrecisionSetToZero(t *testing.T) {
	actual := fmt.Sprintf("%.0f", fakeMantissa)
	assert.Equal(t, "0", actual)
}

func TestFormatWidthAndPrecisionNotInfinite(t *testing.T) {
	var builder strings.Builder
	n, err := fmt.Fprintf(&builder, "%-20.13f", fakeMantissa.WithSignificant(9))
	assert.Equal(t, "0.1234567890000     ", builder.String())
	assert.Equal(t, 20, n)
	assert.NoError(t, err)
}

func TestFormatWidthAndPrecisionNotInfiniteError(t *testing.T) {
	for i := 0; i < 20; i++ {
		w := &maxBytesWriter{maxBytes: i}
		n, err := fmt.Fprintf(w, "%-20.13f", fakeMantissa.WithSignificant(9))
		assert.Equal(t, i, n)
		assert.Error(t, err)
	}
}

func TestFormatNoPrecisionG(t *testing.T) {
	actual := fmt.Sprintf("%g", fakeMantissa)
	assert.Equal(t, "0.1234567890123456", actual)
}

func TestFormatNoPrecisionCapitalG(t *testing.T) {
	actual := fmt.Sprintf("%G", fakeMantissa)
	assert.Equal(t, "0.1234567890123456", actual)
}

func TestFormatNotInfiniteG14(t *testing.T) {
	actual := fmt.Sprintf("%.14g", fakeMantissa.WithSignificant(9))
	assert.Equal(t, "0.123456789", actual)
}

func TestFormatNotInfiniteG(t *testing.T) {
	actual := fmt.Sprintf("%g", fakeMantissa.WithSignificant(9))
	assert.Equal(t, "0.123456789", actual)
}

func TestFormatG7(t *testing.T) {
	actual := fmt.Sprintf("%.7g", fakeMantissa)
	assert.Equal(t, "0.1234567", actual)
}

func TestFormatG0(t *testing.T) {
	actual := fmt.Sprintf("%.0g", fakeMantissa)
	assert.Equal(t, "0.1", actual)
}

func TestFormatV(t *testing.T) {
	actual := fmt.Sprintf("%v", fakeMantissa)
	assert.Equal(t, "0.1234567890123456", actual)
}

func TestFormatE(t *testing.T) {
	actual := fmt.Sprintf("%e", fakeMantissa)
	assert.Equal(t, "0.123456e+00", actual)
}

func TestPrintZero(t *testing.T) {
	var mantissa Mantissa
	actual := mantissa.Sprint(45)
	assert.Equal(t, "0", actual)
}

func TestFormatZero(t *testing.T) {
	var mantissa Mantissa
	actual := fmt.Sprintf("%f", mantissa)
	assert.Equal(t, "0.000000", actual)
}

func TestFormatZeroPrecision(t *testing.T) {
	var mantissa Mantissa
	actual := fmt.Sprintf("%.5f", mantissa)
	assert.Equal(t, "0.00000", actual)
}

func TestFormatZeroWidth(t *testing.T) {
	var mantissa Mantissa
	actual := fmt.Sprintf("%4.0f", mantissa)
	assert.Equal(t, "   0", actual)
}

func TestFormatZeroG(t *testing.T) {
	var mantissa Mantissa
	actual := fmt.Sprintf("%G", mantissa)
	assert.Equal(t, "0", actual)
}

func TestFormatZeroPrecisionG(t *testing.T) {
	var mantissa Mantissa
	actual := fmt.Sprintf("%.5G", mantissa)
	assert.Equal(t, "0", actual)
}

func TestFormatZeroZeroPrecisionG(t *testing.T) {
	var mantissa Mantissa
	actual := fmt.Sprintf("%.0G", mantissa)
	assert.Equal(t, "0", actual)
}

func TestFormatZeroV(t *testing.T) {
	var mantissa Mantissa
	actual := fmt.Sprintf("%5v", mantissa)
	assert.Equal(t, "    0", actual)
}

func TestFormatZeroE(t *testing.T) {
	var mantissa Mantissa
	actual := fmt.Sprintf("%.5E", mantissa)
	assert.Equal(t, "0.00000E+00", actual)
}

func TestFormatBadVerb(t *testing.T) {
	actual := fmt.Sprintf("%h", fakeMantissa)
	assert.Equal(t, "%!h(mantissa=0.1234567890123456)", actual)
}

func TestPrint(t *testing.T) {
	actual := fmt.Sprint(fakeMantissa)
	assert.Equal(t, "0.1234567890123456", actual)
}

func TestPrintNil(t *testing.T) {
	var mantissa Mantissa
	actual := fmt.Sprint(mantissa)
	assert.Equal(t, "0", actual)
}

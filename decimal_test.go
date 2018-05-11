package decimal

import (
	"strings"
	"testing"

	"github.com/gotoxu/assert"
)

func TestFromInt(t *testing.T) {
	tests := []struct {
		input  int64
		output string
	}{
		{-12345, "-12345"},
		{-1, "-1"},
		{1, "1"},
		{-9223372036854775807, "-9223372036854775807"},
		{-9223372036854775808, "-9223372036854775808"},
	}

	for _, tt := range tests {
		dec := NewDecFromInt(tt.input)
		str := dec.ToString()
		assert.DeepEqual(t, str, tt.output)
	}
}

func TestFromUint(t *testing.T) {
	tests := []struct {
		input  uint64
		output string
	}{
		{12345, "12345"},
		{0, "0"},
		{18446744073709551615, "18446744073709551615"},
	}

	for _, tt := range tests {
		var dec MyDecimal
		dec.FromUint(tt.input)
		str := dec.ToString()
		assert.DeepEqual(t, str, tt.output)
	}
}

func TestToInt(t *testing.T) {
	tests := []struct {
		input  string
		output int64
		err    error
	}{
		{"18446744073709551615", 9223372036854775807, ErrOverflow},
		{"-1", -1, nil},
		{"1", 1, nil},
		{"-1.23", -1, ErrTruncated},
		{"-9223372036854775807", -9223372036854775807, nil},
		{"-9223372036854775808", -9223372036854775808, nil},
		{"9223372036854775808", 9223372036854775807, ErrOverflow},
		{"-9223372036854775809", -9223372036854775808, ErrOverflow},
	}

	for _, tt := range tests {
		var dec MyDecimal
		dec.FromString([]byte(tt.input))
		result, ec := dec.ToInt()

		if ec != nil {
			assert.DeepEqual(t, ec, tt.err)
		} else {
			assert.DeepEqual(t, result, tt.output)
		}
	}
}

func TestToUint(t *testing.T) {
	tests := []struct {
		input  string
		output uint64
		err    error
	}{
		{"12345", 12345, nil},
		{"0", 0, nil},
		/* ULLONG_MAX = 18446744073709551615ULL */
		{"18446744073709551615", 18446744073709551615, nil},
		{"18446744073709551616", 18446744073709551615, ErrOverflow},
		{"-1", 0, ErrOverflow},
		{"1.23", 1, ErrTruncated},
		{"9999999999999999999999999.000", 18446744073709551615, ErrOverflow},
	}

	for _, tt := range tests {
		var dec MyDecimal
		dec.FromString([]byte(tt.input))
		result, ec := dec.ToUint()

		if ec != nil {
			assert.DeepEqual(t, ec, tt.err)
		} else {
			assert.DeepEqual(t, result, tt.output)
		}
	}
}

func TestFromFloat(t *testing.T) {
	tests := []struct {
		s string
		f float64
	}{
		{"12345", 12345},
		{"123.45", 123.45},
		{"-123.45", -123.45},
		{"0.00012345000098765", 0.00012345000098765},
		{"1234500009876.5", 1234500009876.5},
	}

	for _, tt := range tests {
		dec, err := NewDecFromFloatForTest(tt.f)
		assert.Nil(t, err)
		str := dec.ToString()
		assert.DeepEqual(t, str, tt.s)
	}
}

func TestToFloat(t *testing.T) {
	tests := []struct {
		s string
		f float64
	}{
		{"12345", 12345},
		{"123.45", 123.45},
		{"-123.45", -123.45},
		{"0.00012345000098765", 0.00012345000098765},
		{"1234500009876.5", 1234500009876.5},
	}

	for _, ca := range tests {
		var dec MyDecimal
		dec.FromString([]byte(ca.s))
		f, err := dec.ToFloat64()
		assert.Nil(t, err)
		assert.DeepEqual(t, f, ca.f)
	}
}

func TestRoundWithHalfEven(t *testing.T) {
	tests := []struct {
		input  string
		scale  int
		output string
		err    error
	}{
		{"123456789.987654321", 1, "123456790.0", nil},
		{"15.1", 0, "15", nil},
		{"15.5", 0, "16", nil},
		{"15.9", 0, "16", nil},
		{"-15.1", 0, "-15", nil},
		{"-15.5", 0, "-16", nil},
		{"-15.9", 0, "-16", nil},
		{"15.1", 1, "15.1", nil},
		{"-15.1", 1, "-15.1", nil},
		{"15.17", 1, "15.2", nil},
		{"15.4", -1, "20", nil},
		{"-15.4", -1, "-20", nil},
		{"5.4", -1, "10", nil},
		{".999", 0, "1", nil},
		{"999999999", -9, "1000000000", nil},
	}

	for _, ca := range tests {
		var dec MyDecimal
		dec.FromString([]byte(ca.input))
		var rounded MyDecimal
		err := dec.Round(&rounded, ca.scale, ModeHalfEven)
		if err != nil {
			assert.DeepEqual(t, err, ca.err)
		} else {
			result := rounded.ToString()
			assert.DeepEqual(t, result, ca.output)
		}
	}
}

func TestRoundWithTruncate(t *testing.T) {
	tests := []struct {
		input  string
		scale  int
		output string
		err    error
	}{
		{"123456789.987654321", 1, "123456789.9", nil},
		{"15.1", 0, "15", nil},
		{"15.5", 0, "15", nil},
		{"15.9", 0, "15", nil},
		{"-15.1", 0, "-15", nil},
		{"-15.5", 0, "-15", nil},
		{"-15.9", 0, "-15", nil},
		{"15.1", 1, "15.1", nil},
		{"-15.1", 1, "-15.1", nil},
		{"15.17", 1, "15.1", nil},
		{"15.4", -1, "10", nil},
		{"-15.4", -1, "-10", nil},
		{"5.4", -1, "0", nil},
		{".999", 0, "0", nil},
		{"999999999", -9, "0", nil},
	}

	for _, ca := range tests {
		var dec MyDecimal
		dec.FromString([]byte(ca.input))
		var rounded MyDecimal
		err := dec.Round(&rounded, ca.scale, ModeTruncate)
		if err != nil {
			assert.DeepEqual(t, err, ca.err)
		} else {
			result := rounded.ToString()
			assert.DeepEqual(t, result, ca.output)
		}
	}
}

func TestRoundWithCeil(t *testing.T) {
	tests := []struct {
		input  string
		scale  int
		output string
		err    error
	}{
		{"123456789.987654321", 1, "123456790.0", nil},
		{"15.1", 0, "16", nil},
		{"15.5", 0, "16", nil},
		{"15.9", 0, "16", nil},
		//TODO:fix me
		{"-15.1", 0, "-16", nil},
		{"-15.5", 0, "-16", nil},
		{"-15.9", 0, "-16", nil},
		{"15.1", 1, "15.1", nil},
		{"-15.1", 1, "-15.1", nil},
		{"15.17", 1, "15.2", nil},
		{"15.4", -1, "20", nil},
		{"-15.4", -1, "-20", nil},
		{"5.4", -1, "10", nil},
		{".999", 0, "1", nil},
		{"999999999", -9, "1000000000", nil},
	}

	for _, ca := range tests {
		var dec MyDecimal
		dec.FromString([]byte(ca.input))
		var rounded MyDecimal
		err := dec.Round(&rounded, ca.scale, modeCeiling)
		if err != nil {
			assert.DeepEqual(t, err, ca.err)
		} else {
			result := rounded.ToString()
			assert.DeepEqual(t, result, ca.output)
		}
	}
}

func TestFromString(t *testing.T) {
	type tcase struct {
		input  string
		output string
		err    error
	}
	tests := []tcase{
		{"12345", "12345", nil},
		{"12345.", "12345", nil},
		{"123.45.", "123.45", nil},
		{"-123.45.", "-123.45", nil},
		{".00012345000098765", "0.00012345000098765", nil},
		{".12345000098765", "0.12345000098765", nil},
		{"-.000000012345000098765", "-0.000000012345000098765", nil},
		{"1234500009876.5", "1234500009876.5", nil},
		{"123E5", "12300000", nil},
		{"123E-2", "1.23", nil},
	}

	for _, ca := range tests {
		var dec MyDecimal
		err := dec.FromString([]byte(ca.input))
		result := dec.ToString()
		assert.Nil(t, err)
		assert.DeepEqual(t, result, ca.output)
	}
}

func TestToString(t *testing.T) {
	type tcase struct {
		input  string
		output string
	}
	tests := []tcase{
		{"123.123", "123.123"},
		{"123.1230", "123.1230"},
		{"00123.123", "123.123"},
	}

	for _, ca := range tests {
		var dec MyDecimal
		dec.FromString([]byte(ca.input))
		result := dec.ToString()
		assert.DeepEqual(t, result, ca.output)
	}
}

func TestCompare(t *testing.T) {
	type tcase struct {
		a   string
		b   string
		cmp int
	}
	tests := []tcase{
		{"12", "13", -1},
		{"13", "12", 1},
		{"-10", "10", -1},
		{"10", "-10", 1},
		{"-12", "-13", 1},
		{"0", "12", -1},
		{"-10", "0", -1},
		{"4", "4", 0},
		{"-1.1", "-1.2", 1},
		{"1.2", "1.1", 1},
		{"1.1", "1.2", -1},
	}

	for _, tt := range tests {
		var a, b MyDecimal
		a.FromString([]byte(tt.a))
		b.FromString([]byte(tt.b))

		cmp, err := a.Compare(&b)
		assert.Nil(t, err)
		assert.DeepEqual(t, cmp, tt.cmp)
	}
}

func TestMaxDecimal(t *testing.T) {
	type tcase struct {
		prec   int
		frac   int
		result string
	}
	tests := []tcase{
		{1, 1, "0.9"},
		{1, 0, "9"},
		{2, 1, "9.9"},
		{4, 2, "99.99"},
		{6, 3, "999.999"},
		{8, 4, "9999.9999"},
		{10, 5, "99999.99999"},
		{12, 6, "999999.999999"},
		{14, 7, "9999999.9999999"},
		{16, 8, "99999999.99999999"},
		{18, 9, "999999999.999999999"},
		{20, 10, "9999999999.9999999999"},
		{20, 20, "0.99999999999999999999"},
		{20, 0, "99999999999999999999"},
		{40, 20, "99999999999999999999.99999999999999999999"},
	}

	for _, tt := range tests {
		var dec MyDecimal
		maxDecimal(tt.prec, tt.frac, &dec)
		str := dec.ToString()
		assert.DeepEqual(t, str, tt.result)
	}
}

func TestAdd(t *testing.T) {
	type testCase struct {
		a      string
		b      string
		result string
		err    error
	}
	tests := []testCase{
		{".00012345000098765", "123.45", "123.45012345000098765", nil},
		{".1", ".45", "0.55", nil},
		{"1234500009876.5", ".00012345000098765", "1234500009876.50012345000098765", nil},
		{"9999909999999.5", ".555", "9999910000000.055", nil},
		{"99999999", "1", "100000000", nil},
		{"989999999", "1", "990000000", nil},
		{"999999999", "1", "1000000000", nil},
		{"12345", "123.45", "12468.45", nil},
		{"-12345", "-123.45", "-12468.45", nil},
		{"-12345", "123.45", "-12221.55", nil},
		{"12345", "-123.45", "12221.55", nil},
		{"123.45", "-12345", "-12221.55", nil},
		{"-123.45", "12345", "12221.55", nil},
		{"5", "-6.0", "-1.0", nil},
		{"2" + strings.Repeat("1", 71), strings.Repeat("8", 81), "8888888890" + strings.Repeat("9", 71), nil},
	}

	for _, tt := range tests {
		a, _ := NewDecFromStringForTest(tt.a)
		b, _ := NewDecFromStringForTest(tt.b)
		var sum MyDecimal
		err := DecimalAdd(a, b, &sum)
		if err != nil {
			assert.DeepEqual(t, err, tt.err)
		} else {
			result := sum.ToString()
			assert.DeepEqual(t, result, tt.result)
		}
	}
}

func TestSub(t *testing.T) {
	type tcase struct {
		a      string
		b      string
		result string
		err    error
	}
	tests := []tcase{
		{".00012345000098765", "123.45", "-123.44987654999901235", nil},
		{"1234500009876.5", ".00012345000098765", "1234500009876.49987654999901235", nil},
		{"9999900000000.5", ".555", "9999899999999.945", nil},
		{"1111.5551", "1111.555", "0.0001", nil},
		{".555", ".555", "0", nil},
		{"10000000", "1", "9999999", nil},
		{"1000001000", ".1", "1000000999.9", nil},
		{"1000000000", ".1", "999999999.9", nil},
		{"12345", "123.45", "12221.55", nil},
		{"-12345", "-123.45", "-12221.55", nil},
		{"123.45", "12345", "-12221.55", nil},
		{"-123.45", "-12345", "12221.55", nil},
		{"-12345", "123.45", "-12468.45", nil},
		{"12345", "-123.45", "12468.45", nil},
	}

	for _, tt := range tests {
		var a, b, sum MyDecimal
		a.FromString([]byte(tt.a))
		b.FromString([]byte(tt.b))
		err := DecimalSub(&a, &b, &sum)
		if err != nil {
			assert.DeepEqual(t, err, tt.err)
		} else {
			result := sum.ToString()
			assert.DeepEqual(t, result, tt.result)
		}
	}
}

func TestMul(t *testing.T) {
	type tcase struct {
		a      string
		b      string
		result string
		err    error
	}
	tests := []tcase{
		{"12", "10", "120", nil},
		{"-123.456", "98765.4321", "-12193185.1853376", nil},
		{"-123456000000", "98765432100000", "-12193185185337600000000000", nil},
		{"123456", "987654321", "121931851853376", nil},
		{"123456", "9876543210", "1219318518533760", nil},
		{"123", "0.01", "1.23", nil},
		{"123", "0", "0", nil},
		{"1" + strings.Repeat("0", 60), "1" + strings.Repeat("0", 60), "0", ErrOverflow},
	}

	for _, tt := range tests {
		var a, b, product MyDecimal
		a.FromString([]byte(tt.a))
		b.FromString([]byte(tt.b))
		err := DecimalMul(&a, &b, &product)
		if err != nil {
			assert.DeepEqual(t, err, tt.err)
		} else {
			result := product.ToString()
			assert.DeepEqual(t, result, tt.result)
		}
	}
}

func TestDivMod(t *testing.T) {
	type tcase struct {
		a      string
		b      string
		result string
		err    error
	}
	tests := []tcase{
		{"120", "10", "12.000000000", nil},
		{"123", "0.01", "12300.000000000", nil},
		{"120", "100000000000.00000", "0.000000001200000000", nil},
		{"123", "0", "", ErrDivByZero},
		{"0", "0", "", ErrDivByZero},
		{"-12193185.1853376", "98765.4321", "-123.456000000000000000", nil},
		{"121931851853376", "987654321", "123456.000000000", nil},
		{"0", "987", "0", nil},
		{"1", "3", "0.333333333", nil},
		{"1.000000000000", "3", "0.333333333333333333", nil},
		{"1", "1", "1.000000000", nil},
		{"0.0123456789012345678912345", "9999999999", "0.000000000001234567890246913578148141", nil},
		{"10.333000000", "12.34500", "0.837019036046982584042122316", nil},
		{"10.000000000060", "2", "5.000000000030000000", nil},
		{"51", "0.003430", "14868.804664723032069970", nil},
	}

	for _, tt := range tests {
		var a, b, to MyDecimal
		a.FromString([]byte(tt.a))
		b.FromString([]byte(tt.b))
		err := doDivMod(&a, &b, &to, nil, 5)

		if tt.err == ErrDivByZero {
			continue
		}
		if err != nil {
			assert.DeepEqual(t, err, tt.err)
		} else {
			result := to.ToString()
			assert.DeepEqual(t, result, tt.result)
		}
	}

	tests = []tcase{
		{"234", "10", "4", nil},
		{"234.567", "10.555", "2.357", nil},
		{"-234.567", "10.555", "-2.357", nil},
		{"234.567", "-10.555", "2.357", nil},
		{"99999999999999999999999999999999999999", "3", "0", nil},
		{"51", "0.003430", "0.002760", nil},
	}
	for _, tt := range tests {
		var a, b, to MyDecimal
		a.FromString([]byte(tt.a))
		b.FromString([]byte(tt.b))
		ec := doDivMod(&a, &b, nil, &to, 0)
		if tt.err == ErrDivByZero {
			continue
		}

		if ec != nil {
			assert.DeepEqual(t, ec, tt.err)
		} else {
			result := to.ToString()
			assert.DeepEqual(t, result, tt.result)
		}
	}

	tests = []tcase{
		{"1", "1", "1.0000", nil},
		{"1.00", "1", "1.000000", nil},
		{"1", "1.000", "1.0000", nil},
		{"2", "3", "0.6667", nil},
		{"51", "0.003430", "14868.8047", nil},
	}
	for _, tt := range tests {
		var a, b, to MyDecimal
		a.FromString([]byte(tt.a))
		b.FromString([]byte(tt.b))
		ec := DecimalDiv(&a, &b, &to, DivFracIncr)
		if tt.err == ErrDivByZero {
			continue
		}

		if ec != nil {
			assert.DeepEqual(t, ec, tt.err)
		} else {
			s, _ := to.String()
			assert.DeepEqual(t, s, tt.result)
		}
	}

	tests = []tcase{
		{"1", "2.0", "1.0", nil},
		{"1.0", "2", "1.0", nil},
		{"2.23", "3", "2.23", nil},
		{"51", "0.003430", "0.002760", nil},
	}
	for _, tt := range tests {
		var a, b, to MyDecimal
		a.FromString([]byte(tt.a))
		b.FromString([]byte(tt.b))
		ec := DecimalMod(&a, &b, &to)

		if tt.err == ErrDivByZero {
			continue
		}

		if ec != nil {
			assert.DeepEqual(t, ec, tt.err)
		} else {
			s, _ := to.String()
			assert.DeepEqual(t, s, tt.result)
		}
	}
}

func TestMaxOrMin(t *testing.T) {
	type tcase struct {
		neg    bool
		prec   int
		frac   int
		result string
	}
	tests := []tcase{
		{true, 2, 1, "-9.9"},
		{false, 1, 1, "0.9"},
		{true, 1, 0, "-9"},
		{false, 0, 0, "0"},
		{false, 4, 2, "99.99"},
	}
	for _, tt := range tests {
		dec, _ := NewMaxOrMinDec(tt.neg, tt.prec, tt.frac)
		s, _ := dec.String()
		assert.DeepEqual(t, s, tt.result)
	}
}

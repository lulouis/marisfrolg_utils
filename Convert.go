package marisfrolg_utils

import (
	"math"
	"math/big"
	"reflect"
	"time"
	"unsafe"
)

func ToFloat64(ori []byte) (re float64) {
	var bi big.Int
	var neg bool
	var i int

	neg, i = decodeDecimal(ori, &bi)
	re = bigIntToFloat(neg, &bi, i)
	return re
}

func Substr(str string, start int, length int) string {
	rs := []rune(str)
	rl := len(rs)
	end := 0

	if start < 0 {
		start = rl - 1 + start
	}
	end = start + length

	if start > end {
		start, end = end, start
	}

	if start < 0 {
		start = 0
	}
	if start > rl {
		start = rl
	}
	if end < 0 {
		end = 0
	}
	if end > rl {
		end = rl
	}

	return string(rs[start:end])
}

// 获取source的子串,如果start小于0或者end大于source长度则返回""
// start:开始index，从0开始，包括0
// end:结束index，以end结束，但不包括end
func SubString(source string, start int, end int) string {
	var r = []rune(source)
	length := len(r)

	if start < 0 || end > length || start > end {
		return ""
	}

	if start == 0 && end == length {
		return source
	}

	return string(r[start:end])
}

func decodeDecimal(b []byte, m *big.Int) (bool, int) {

	//bigint word size (*--> src/pkg/math/big/arith.go)
	const (
		dec128Bias = 6176
		// Compute the size _S of a Word in bytes.
		_m    = ^big.Word(0)
		_logS = _m>>8&1 + _m>>16&1 + _m>>32&1
		_S    = 1 << _logS
	)

	neg := (b[15] & 0x80) != 0
	exp := int((((uint16(b[15])<<8)|uint16(b[14]))<<1)>>2) - dec128Bias

	b14 := b[14]  // save b[14]
	b[14] &= 0x01 // keep the mantissa bit (rest: sign and exp)

	//most significand byte
	msb := 14
	for msb > 0 {
		if b[msb] != 0 {
			break
		}
		msb--
	}

	//calc number of words
	numWords := (msb / _S) + 1
	w := make([]big.Word, numWords)

	k := numWords - 1
	d := big.Word(0)
	for i := msb; i >= 0; i-- {
		d |= big.Word(b[i])
		if k*_S == i {
			w[k] = d
			k--
			d = 0
		}
		d <<= 8
	}
	b[14] = b14 // restore b[14]
	m.SetBits(w)
	return neg, exp
}

func bigIntToFloat(sign bool, m *big.Int, exp int) float64 {
	var neg int64
	if sign {
		neg = -1
	} else {
		neg = 1
	}

	return float64(neg*m.Int64()) * math.Pow10(exp)
}

func Round(f float64, n int) float64 {
	n10 := math.Pow10(n)
	return math.Trunc((f+0.5/n10)*n10) / n10
}

func MapToStruct(data map[string]interface{}, result interface{}) error {
	rv := indirect(reflect.ValueOf(result), false)
	for key, value := range data {
		fieldValue := rv.FieldByName(key)
		if fieldValue.IsValid() {
			valueValue := reflect.ValueOf(value)
			if value != nil {
				if fieldValue.Type().Name() == "Time" && valueValue.Type().Name() == "string" {
					setUnexportedField(fieldValue, reflect.ValueOf(parseOracleTime(value)))
				} else {
					setUnexportedField(fieldValue, valueValue)
				}
			}
		}
	}
	return nil
}

func setUnexportedField(field reflect.Value, value reflect.Value) {
	reflect.NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr())).Elem().Set(value)
}

func indirect(v reflect.Value, decodingNull bool) reflect.Value {
	v0 := v
	haveAddr := false

	if v.Kind() != reflect.Pointer && v.Type().Name() != "" && v.CanAddr() {
		haveAddr = true
		v = v.Addr()
	}
	for {
		if v.Kind() == reflect.Interface && !v.IsNil() {
			e := v.Elem()
			if e.Kind() == reflect.Pointer && !e.IsNil() && (!decodingNull || e.Elem().Kind() == reflect.Pointer) {
				haveAddr = false
				v = e
				continue
			}
		}

		if v.Kind() != reflect.Pointer {
			break
		}

		if decodingNull && v.CanSet() {
			break
		}

		if v.Elem().Kind() == reflect.Interface && v.Elem().Elem() == v {
			v = v.Elem()
			break
		}
		if v.IsNil() {
			v.Set(reflect.New(v.Type().Elem()))
		}

		if haveAddr {
			v = v0
			haveAddr = false
		} else {
			v = v.Elem()
		}
	}
	return v
}

func parseOracleTime(input interface{}) (result time.Time) {
	if _date, ok := input.(string); ok {
		result, _ = time.ParseInLocation("2006-01-02T15:04:05", _date, CstZone)
	}
	return result
}

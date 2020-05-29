package Marisfrolg_utils

import (
	"errors"
	"math"
	"math/big"
	"math/rand"
	"reflect"
	"strconv"
	"time"
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

//获取source的子串,如果start小于0或者end大于source长度则返回""
//start:开始index，从0开始，包括0
//end:结束index，以end结束，但不包括end
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

func ValueToString(data interface{}) string {
	value := reflect.ValueOf(data)
	switch value.Kind() {
	case reflect.Bool:
		return strconv.FormatBool(value.Bool())

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(value.Int(), 10)

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(value.Uint(), 10)

	case reflect.Float32, reflect.Float64:
		return strconv.FormatFloat(value.Float(), 'g', -1, 64)

	case reflect.String:
		return value.String()
	}

	return ""
}

// 随机生成大写字母
func GetRandomString(l int) string {
	str := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

const (
	LongDateFormat  = "2006-01-02 15:04:05"
	ShortDateFormat = "2006-01-02"
)

//获取日期格式
func GetLongDateString(date string, Hours int64) (dateString string, err error) {
	if len(date) <= 0 {
		return "", errors.New("时间不能为空")
	}
	inputDate, err := time.Parse(LongDateFormat, date)
	if err == nil {
		h, _ := time.ParseDuration("1h")
		d := inputDate.Add(time.Duration(Hours) * h)
		return d.Format(LongDateFormat), err
	} else {
		return "", errors.New("时间格式错误")
	}
}

//获取日期格式
func GetShortDateString(date string, Hours int64) (dateString string, err error) {
	if len(date) <= 0 {
		return "", errors.New("时间为空")
	}
	inputDate, err := time.Parse(ShortDateFormat, date)
	if err == nil {
		h, _ := time.ParseDuration("1h")
		d := inputDate.Add(time.Duration(Hours) * h)
		return d.Format(LongDateFormat), err
	} else {
		return "", errors.New("时间格式错误")
	}
}

//获取相差时间
func GetMinuteDiffer(start_time, end_time string) int64 {
	var hour int64
	t1, err := time.ParseInLocation("2006-01-02 15:04:05", start_time, time.Local)
	t2, err := time.ParseInLocation("2006-01-02 15:04:05", end_time, time.Local)
	if err == nil && t1.Before(t2) {
		diff := t2.Unix() - t1.Unix() //
		hour = diff / 60
		return hour
	} else {
		return hour
	}
}

//获取某一天的0点时间
func GetZeroTime(d time.Time) time.Time {
	return time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, d.Location())
}

//获取某一天的23:59:59点时间
func GetZeroTimeEnd(d time.Time) time.Time {
	return time.Date(d.Year(), d.Month(), d.Day(), 23, 59, 59, 59, d.Location())
}

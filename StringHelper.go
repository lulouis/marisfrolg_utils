package marisfrolg_utils

import (
	"math/rand"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"time"
)

//位数不够自动左补全
//str 传入的字符串 ，totalWidth 补的长度 ，paddingChar 补全的字符
func PadLeft(str string, totalWidth int, paddingChar string) (r string) {
	var result string
	if str == "" {
		return ""
	} else {
		l := len(str)
		result = str
		for i := 0; i < totalWidth-l; i++ {
			result = paddingChar + result
		}
	}
	return result
}

//位数不够自动右补全
// str 传入的字符串 ，totalWidth 补的长度 ，paddingChar 补全的字符
func PadRight(str string, totalWidth int, paddingChar string) (r string) {
	var result string
	if str == "" {
		return ""
	} else {
		l := len(str)
		result = str
		for i := 0; i < totalWidth-l; i++ {
			result = result + paddingChar
		}
	}
	return result
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

//判断是否存在
func Contains(array interface{}, val interface{}) (index int) {
	index = -1
	switch reflect.TypeOf(array).Kind() {
	case reflect.Slice:
		{
			s := reflect.ValueOf(array)
			for i := 0; i < s.Len(); i++ {
				if reflect.DeepEqual(val, s.Index(i).Interface()) {
					index = i
					return
				}
			}
		}
	}
	return
}
//验证手机号
func ValidateMobile(mobileNum string) bool {
	regMobile := `^1([38][0-9]|4[579]|5[0-3,5-9]|6[6]|7[0135678]|9[89])\d{8}$`
	reg := regexp.MustCompile(regMobile)
	return reg.MatchString(mobileNum)
}

//移除重复数据
func RemoveRepeatedElement(arr []string) (newArr []string) {
	newArr = make([]string, 0)
	sort.Strings(arr)
	for i := 0; i < len(arr); i++ {
		repeat := false
		for j := i + 1; j < len(arr); j++ {
			if arr[i] == arr[j] {
				repeat = true
				break
			}
		}
		if !repeat {
			newArr = append(newArr, arr[i])
		}
	}
	return
}

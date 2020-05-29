package marisfrolg_utils

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

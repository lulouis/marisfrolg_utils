package marisfrolg_utils

import (
	"errors"
	"math"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

//计算公式结果
func Calculate(Str string) interface{} {
	if len(Str) <= 0 {
		return nil
	}
	postfix := Infix2ToPostfix(Str)
	var tmp string
	sk := InitStack()
	var c string
	var Operand string
	var x, y float64
	for i := 0; i < len(postfix); i++ {
		c = string(postfix[i])
		//added c==',' for germany culture
		if unicode.IsDigit(rune(postfix[i])) || c == "." || c == "," {
			//数据值收集.
			Operand += string(c)
		} else if c == " " && len(Operand) > 0 {
			//运算数转换
			tmp = Operand
			//负数的转换
			if strings.HasPrefix(tmp, "-") {
				v2, _ := strconv.ParseFloat(tmp[1:len(tmp)-1], 64)
				sk.Push(-v2)
			} else {
				v2, _ := strconv.ParseFloat(tmp, 64)
				sk.Push(v2)
			}
			Operand = ""
		} else if c == "+" || c == "-" || c == "*" || c == "/" || c == "%" || c == "^" { //运算符处理.双目运算处理.
			// 双目运算
			if sk.Count() > 0 { /*如果输入的表达式根本没有包含运算符.或是根本就是空串.这里的逻辑就有意义了.*/
				y = sk.Pop().(float64)
			} else {
				sk.Push(0)
				break
			}
			if sk.Count() > 0 {
				x = sk.Pop().(float64)
			} else {
				sk.Push(y)
				break
			}
			switch c {
			case "+":
				sk.Push(x + y)
				break
			case "-":
				sk.Push(x - y)
				break
			case "*":
				if y == 0 {
					sk.Push(x * 1)
				} else {
					sk.Push(x * y)
				}
				break
			case "/":
				if y == 0 {
					sk.Push(x / 1)
				} else {
					sk.Push(x / y)
				}
				break
			case "%":
				//sk.Push(x % y)
				break
			case "^":
				if x > 0 {
					sk.Push(math.Pow(x, y))
				} else {
					//t := y
					//t = 1 / (2 * t)
					//ts := strconv.FormatFloat(t, 'E', -1, 62)

				}
				break
			}
		} else if c == "!" { //单目取反. )
			sk.Push(-(sk.Pop().(float64)))
		}
	}
	if sk.Count() > 1 {
		return nil
	}
	if sk.Count() == 0 {
		return nil
	}
	return sk.Pop()
}

func Infix2ToPostfix(exp string) string {
	sb := exp
	sk := InitStack()
	var re string
	var c string
	for i := 0; i < len(sb); i++ {
		c = string(sb[i])
		//数字要
		if unicode.IsDigit(rune(sb[i])) || c == "," {
			re += c
		}
		unicode.IsLetter(rune(sb[i])) //如果是空白,那么不要.现在字母也不要.
		switch c {                    //如果是其它字符...列出的要,没有列出的不要.
		case "+", "-", "*", "/", "%", "^", "!", "(", ")", ".":
			re += c
			break
		default:
			continue
		}
	}
	sb = re
	rb := []rune(sb)
	sb = ""
	//对负号进行预转义处理.负号变单目运算符求反.
	for i := 0; i < len(rb); i++ {
		if string(rb[i]) == "-" && (i == 0 || string(rb[i-1]) == "(") {
			rb[i] = 33
		}
		sb += string(rb[i])
	}
	//字符转义.
	//将中缀表达式变为后缀表达式.
	re = ""
	for i := 0; i < len(sb); i++ {
		r := string(sb[i])
		//如果是数值.
		if unicode.IsDigit(rune(sb[i])) || r == "." {
			re += r
			//加入后缀式
		} else if r == "+" || r == "-" || r == "*" || r == "/" || r == "%" || r == "^" || r == "!" {
			//运算符处理
			//栈不为空时
			for {
				if sk.Count() > 0 {
					c = sk.Pop().(string)
					//将栈中的操作符弹出.
					//如果发现左括号.停.
					if c == "(" {
						sk.Push(c)
						//将弹出的左括号压回.因为还有右括号要和它匹配.
						break
						//中断.
					} else {
						//如果优先级比上次的高,则压栈.
						if Power(c) < Power(r) {
							sk.Push(c)
							break
						} else {
							re += " "
							re += c
						}
						//如果不是左括号,那么将操作符加入后缀式中.
					}
				} else {
					break
				}
			}
			sk.Push(r)
			//把新操作符入栈.
			re += " "

		} else if r == "(" { //基本优先级提升
			sk.Push("(")
			re += " "
		} else if r == ")" { //基本优先级下调
			for {
				if sk.Count() > 0 { //栈不为空时
					c = sk.Pop().(string)
					//pop Operator
					if c != "(" {
						re += " "
						re += c
						//加入空格主要是为了防止不相干的数据相临产生解析错误.
						re += " "
					} else {
						break
					}
				}
			}
		} else {
			re += string(sb[i])
		}

	}
	for {
		if sk.Count() > 0 { //这是最后一个弹栈啦.
			re += " "
			re += sk.Pop().(string)
		} else {
			break
		}
	}
	re += " "
	return FormatSpace(re)
	//在这里进行一次表达式格式化.这里就是后缀式了.
}

//运算优先级别
func Power(opr string) (p int) {
	switch opr {
	case "+", "-":
		p = 1
	case "*", "/":
		p = 2
	case "%", "^", "!":
		p = 3
	default:
		p = 0
	}
	return
}

// 规范化逆波兰表达式.
func FormatSpace(s string) (exp string) {
	postfix := ""
	for i := 0; i < len(s); i++ {
		if !(len(s) > i+1 && s[i] == ' ' && s[i+1] == ' ') {
			postfix += string(s[i])
		} else {
			postfix += string(s[i])
		}
	}
	return postfix
}

//公式检测
func CalculateCheck(str string) (flag bool, err error) {
	if len(str) <= 0 {
		return false, errors.New("空值")
	}

	//如果有连续的符号的话返回false
	if f, _ := regexp.MatchString("[+-/\\*]{2,}", str); f {
		return false, errors.New("存在连续符号")
	}

	//如果出现连续的括号的话返回false
	if f, _ := regexp.MatchString("[\\(\\)]{2,}", str); f {
		return false, errors.New("存在连续符号")
	}

	//如果左括号后面出现+*/符号的时候返回false
	if f, _ := regexp.MatchString("\\([+/\\*]+", str); f {

		return false, errors.New("左括号后面出现+*/符号")
	}

	//如果右括号后面出现数字或者没有出现+-*/的时候返回false
	if f, _ := regexp.MatchString("\\)[^+-/\\*]", str); f {
		return false, errors.New("右括号后面出现数字或者没有出现+-*/")
	}
	//+-*/ 后面没有数字
	// if f, _ := regexp.MatchString("[+-/\\*]^[1-9]\\d*$", str); f {
	// 	return false, errors.New("右括号后面出现数字或者没有出现+-*/")
	// }

	//如果左括号前面没有出现+-*/的时候返回false
	if f, _ := regexp.MatchString("[^+-/\\*]+\\(", str); f {
		return false, errors.New("左括号前面没有出现+-*/")
	}
	//如果右括号前面出现+-*/的时候返回fasle
	if f, _ := regexp.MatchString("[+-/\\*]+\\)", str); f {
		return false, errors.New("右括号前面出现+-*/")
	}

	//递归检查括号是否成对出现
	var item string
	sk := InitStack()
	for i := 0; i < len(str); i++ {
		item = string(str[i])
		// if f, _ := regexp.MatchString("[+-/\\*]", item); f {
		// 	if i+1 < len(str) {
		// 		item1 := string(str[i+1])
		// 		if g, _ := regexp.MatchString("^[0-9]*$", item1); !g {
		// 			return false, errors.New("运算符后面接数字")
		// 		}
		// 	}
		// }
		if item == "(" {
			sk.Push('(')
		} else if item == ")" {
			if sk.Count() > 0 {
				sk.Pop()
			} else {
				return false, errors.New("括号不成对")
			}
		}
	}
	if sk.Count() != 0 {
		return false, errors.New("括号不成对")
	}
	return true, err
}

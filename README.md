# marisfrolg_utils
玛丝菲尔golang项目帮助类库

## 项目引用
go get github.com/lulouis/marisfrolg_utils

## 使用案例
帮助类分为9个部分

1、字符串类

2、数据库类


### 案例1：字符串函数

(```)
package main
import (
	"fmt"

	convert "github.com/lulouis/marisfrolg_utils"
)

func main() {
	substring := convert.Substr("marisfrolg_utils", 0, 10)
	fmt.Println(substring)

}

(```)

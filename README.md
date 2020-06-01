# marisfrolg_utils
玛丝菲尔golang项目帮助类库

## 项目引用
go get github.com/lulouis/marisfrolg_utils

## 使用案例
帮助类分为10个部分

>1、Convert

>2、数据库类(DBHelper)

```
1、ExecuteNonQueryByTran func ExecuteNonQueryByTran(db *sql.DB, SqlList []string) error  批量执行 含事务 nil 成功 err 失败 MongoDB 禁止使用，其他数据库自行斟酌（目前支持Oracle）
   测试指令: go test -v -run TestExecuteNonQueryByTran DB_test.go
   参数说明：db: 数据库链接 Sql: 语句集合
2、AssemblyParameters(IdList, Field string) (condition string) SQL IN()的查询里不能超过1000列，将大于1000列的以900为间隔分开组装
   测试指令： go test -v -run TestAssemblyParameters DB_test.go
   参数说明：IdList:数据列 例如："'987654321','123456789'"  Field:字段名
3、GetSqlList(Parameters string, Field string) string  对数据库查询请求参数的优化
   测试指令：go test -v -run TestGetSqlList DB_test.go
   参数说明：Parameters:数据查询参数; Field: 字段名
4、StringToRuneArr(Parameters string) []string 对HANA查询结果的特殊符号优化
   测试指令:go test -v -run TestStringToRuneArr DB_test.go
   参数说明:Parameters  数据结果列请求参数
 ```
>3、身份证加解密相关操作(DesHelper)

>4、发送邮件相关操作(EmailHelper)
 ```
1、(emailServer EmailServer) SendMail(emailRequest EmailRequest) error 发送邮件
   测试指令：go test -v -run TestSendMail Email_test.go
   参数说明：emailServer:配置发件邮箱账号和密码端口号,emailRequest 配置收件人、标题、主题、正文、附加路径
 ```



>5、文件相关操作(FileHelper)
 ```
1、PathlogExistsFile(path string) (bool, error)检查制定路径下是否存在文件如果不存在直接创建文件夹
  测试指令:go test -v -run TestPathlogExistsFile File_test.go
  参数说明:path:文件夹路径
 ```

>6、Http相关操作(HttpHelper)
```
1、HttpPostOnlyBody(apiURL string, parmbody string, token string) (rs []byte, err error) POST请求带token验证的URL
  测试指令:go test -v -run TestHttpPostOnlyBody Http_test.go
  参数说明：apiURL:请求路径;parmbody:Body参数;token:需要验证的token
```


>7、文本日志相关操作类(LogUtil)
```
1、AddOperationLog(Type string, Title string, Message string,Filepath string) 添加操作日志 
   测试指令: go test -v -run TestLog Log_test.go
   参数说明：Type:文件名称(可以为空) Title:函数名称 Message:要记录在日志的内容 Filepath 日志文件要存放的路径
   日志格式: 日期：[2020-05-29 17:29:02]    IP：169.254.126.100
            标题：测试
            内容：消息内容
2、PrintMyName() 获取此函数被哪一个函数调用的上一级函数名称
   测试指令:go test -v -run TestPrintMyName Log_test.go
   参数说明:参数无
3、GetFileName() 获取当前文件名称
   测试指令:go test -v -run TestGetFileName Log_test.go
   参数说明:参数无
```

>8、加解密相关操作类(RSA)

>9、Redis缓存相关操作类(RedisHelper)

>10、字符串相关操作类(StringHelper)
```
1、PadLeft(str string, totalWidth int, paddingChar string) (r string)位数不够自动左补全
   测试指令:go test -v -run TestPadLeft StringHelper_test.go
   参数说明:str 传入的字符串 ，totalWidth 补的长度 ，paddingChar 补全的字符
2、PadRight(str string, totalWidth int, paddingChar string) (r string)位数不够自动右补全
   测试指令:go test -v -run TestPadRight StringHelper_test.go
   参数说明:str 传入的字符串 ，totalWidth 补的长度 ，paddingChar 补全的字符
3、GetRandomString(Length int) string 随机生成指定位数的大写字母
   测试指令:go test -v -run TestGetRandomString StringHelper_test.go
   参数说明:Length:位数
4、ValidateMobile(mobileNum string) bool 验证手机号是否符合要求
   测试指令:go test -v -run TestValidateMobile StringHelper_test.go
   参数说明:mobileNum:手机号码
5、RemoveRepeatedElement(arr []string) (newArr []string)  移除string数组中重复数据,返回已经处理好的无重复的string数组
   测试指令:go test -v -run TestRemoveRepeatedElement StringHelper_test.go
   参数说明:arr 带有重复的string数组
```
## 测试案例 Usage 参考网站：http://c.biancheng.net/view/124.html
新建xxx_test.go文件
1.测试整个文件:
```
go test -v Convert_test.go
```
2.指定某个库某个方法进行测试:
```
go test -v -run TestPadRight StringHelper_test.go
```
### 案例1：字符串函数
```
package main

import (
	"fmt"
	"testing"

	convert "github.com/lulouis/marisfrolg_utils"
)

func TestConvert(t *testing.T) {
	substring := convert.Substr("marisfrolg_utils", 0, 10)
	fmt.Println(substring)

}
go test -v Convert_test.go
```
### 案例2:请查看example里面的 StringHelper_test.go
运行测试案例
```
cd example 
go test -v StringHelper_test.go
```

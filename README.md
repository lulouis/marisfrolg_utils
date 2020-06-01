# marisfrolg_utils
玛丝菲尔golang项目帮助类库

## 项目引用
go get github.com/lulouis/marisfrolg_utils

## 使用案例
帮助类分为10个部分

>1、Convert

>2、数据库类(DBHelper)

```
1、 ExecuteNonQueryByTran func ExecuteNonQueryByTran(db *sql.DB, SqlList []string) error  批量执行 含事务 nil 成功 err 失败 MongoDB 禁止使用，其他数据库自行斟酌（目前支持Oracle）
  测试指令: go test -v -run TestExecuteNonQueryByTran DB_test.go
  参数说明：db: 数据库链接 Sql: 语句集合
2、AssemblyParameters(IdList, Field string) (condition string) SQL IN()的查询里不能超过1000列，将大于1000列的以900为间隔分开组装
  测试指令： go test -v -run TestAssemblyParameters DB_test.go
  参数说明：IdList:数据列 例如："'987654321','123456789'"  Field:字段名
 ```
>3、 身份证加解密相关操作(DesHelper)

>4、 发送邮件相关操作(EmailHelper)

>5、 文件相关操作(FileHelper)


>6、 Http相关操作(HttpHelper)


>7、 文本日志相关操作类(LogUtil)
```
1、 AddOperationLog(Type string, Title string, Message string,Filepath string) 添加操作日志 
   测试指令: go test -v -run TestLog Log_test.go
   参数说明：Type:文件名称(可以为空) Title:函数名称 Message:要记录在日志的内容 Filepath 日志文件要存放的路径
   日志格式: 日期：[2020-05-29 17:29:02]    IP：169.254.126.100
            标题：测试
            内容：消息内容
2、 PrintMyName() 获取此函数被哪一个函数调用的上一级函数名称
3、 GetFileName() 获取当前文件名称
```

>8、 加解密相关操作类(RSA)

>9、 Redis缓存相关操作类(RedisHelper)

>10、 字符串相关操作类(StringHelper)
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

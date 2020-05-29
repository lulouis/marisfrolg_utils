package main

import (
	convert "github.com/lulouis/marisfrolg_utils"
	"testing"
)

func TestLog(t *testing.T) {
	convert.AddOperationLog("文件名", "测试", "消息内容", "LOG")
}

package main

import (
	 "github.com/lulouis/marisfrolg_utils"
	"testing"
)

func TestLog(t *testing.T) {
	marisfrolg_utils.AddOperationLog("文件名", "测试", "消息内容", "LOG")
}

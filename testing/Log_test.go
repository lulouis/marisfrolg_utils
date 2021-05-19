package testing

import (
	"fmt"
	"testing"

	"github.com/lulouis/marisfrolg_utils"
)

func TestLog(t *testing.T) {
	marisfrolg_utils.AddOperationLog("文件名", "测试", "消息内容", "LOG")
}

func TestPrintMyName(t *testing.T) {
	name := marisfrolg_utils.PrintMyName()
	fmt.Println(name)
}

func TestGetFileName(t *testing.T) {
	name := marisfrolg_utils.GetFileName()
	fmt.Println(name)
}

package testing

import (
	"fmt"
	"testing"

	"github.com/lulouis/marisfrolg_utils"
)

func TestLog(t *testing.T) {
	var err error
	logInfo := marisfrolg_utils.NewAuditLogInfo("marisfrolg_utils", marisfrolg_utils.GetMethodName(), "00000")
	logInfo.SetRequest("测试，你好")
	defer func() {
		if rec := recover(); rec != nil {
			err = fmt.Errorf("%v", rec)
		}
		if err != nil {
			marisfrolg_utils.AuditLog(logInfo, err)
		}
	}()
}

func TestPrintMyName(t *testing.T) {
	name := marisfrolg_utils.PrintMyName()
	fmt.Println(name)
}

func TestGetFileName(t *testing.T) {
	name := marisfrolg_utils.GetFileName()
	fmt.Println(name)
}

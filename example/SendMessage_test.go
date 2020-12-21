package main

import (
	"fmt"
	"github.com/lulouis/marisfrolg_utils"
	"testing"
)

func TestSendMessage(t *testing.T) {
	r := marisfrolg_utils.SendShortMessage("hangzhou","LTAI4G6joAA8gwV3pd1FUoG4","accessKeySecret","13534023573","测试","qwerty","")
	fmt.Println(r)
}

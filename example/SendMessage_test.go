package main

import (
	"fmt"
	"github.com/lulouis/marisfrolg_utils"
	"testing"
)

func TestSendMessage(t *testing.T) {
	r := marisfrolg_utils.SendShortMessage("regionId","accessKeyId","accessKeySecret","12345678910","测试","qwerty","")
	fmt.Println(r)
}

package main

import (
	"fmt"
	"github.com/lulouis/marisfrolg_utils"
	"testing"
)

func TestHttpPostOnlyBody(t *testing.T) {
	r, err := marisfrolg_utils.HttpPostOnlyBody("www.baidu.com",fmt.Sprintf(`{"data":"123"}`),"123456")
	fmt.Println(string(r),err)
}

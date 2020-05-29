package main

import (
	"fmt"
	"github.com/lulouis/marisfrolg_utils"
	"testing"
)

func TestPadLeft(t *testing.T){
	r:=marisfrolg_utils.PadLeft("marisfrolg",15,"p")
	fmt.Println(r)
}
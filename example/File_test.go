package main

import (
	"fmt"
	"github.com/lulouis/marisfrolg_utils"
	"testing"
)

func TestPathlogExistsFile(t *testing.T) {
	file, err := marisfrolg_utils.PathlogExistsFile("./Log")
	fmt.Println(file,err)
}

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

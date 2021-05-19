package testing

import (
	"fmt"
	"testing"

	"github.com/lulouis/marisfrolg_utils"
)

// go test -v -run Convert .\testing\Convert_test.go
func TestConvert(t *testing.T) {
	substring := marisfrolg_utils.Substr("marisfrolg_utils", 0, 10)
	fmt.Println(substring)

}

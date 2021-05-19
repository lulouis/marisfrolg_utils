package testing

import (
	"fmt"
	"testing"

	"github.com/lulouis/marisfrolg_utils"
)

// go test -v -run Calculate .\testing\Calculate_test.go
func TestCalculate(t *testing.T) {
	result := marisfrolg_utils.Calculate("(1+10)*100*0.2")
	fmt.Println(result)
}

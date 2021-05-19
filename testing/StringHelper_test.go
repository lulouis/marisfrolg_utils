package testing

import (
	"fmt"
	"testing"

	"github.com/lulouis/marisfrolg_utils"
)

func TestPadLeft(t *testing.T) {
	r := marisfrolg_utils.PadLeft("marisfrolg", 15, "p")
	fmt.Println(r)
}

func TestPadRight(t *testing.T) {
	r := marisfrolg_utils.PadRight("marisfrolg", 15, "p")
	fmt.Println(r)
}

func TestGetRandomString(t *testing.T) {
	r := marisfrolg_utils.GetRandomString(5)
	fmt.Println(r)
}
func TestValidateMobile(t *testing.T) {
	r := marisfrolg_utils.ValidateMobile("123456789")
	fmt.Println(r)
}

func TestRemoveRepeatedElement(t *testing.T) {
	var OldString []string
	OldString = append(OldString, "ABC")
	OldString = append(OldString, "ABC")
	OldString = append(OldString, "DEF")
	r := marisfrolg_utils.RemoveRepeatedElement(OldString)
	fmt.Println(r)
}

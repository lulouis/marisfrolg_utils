package testing

import (
	"fmt"
	"image"
	"os"
	"testing"

	"github.com/lulouis/marisfrolg_utils"
)

func TestPathlogExistsFile(t *testing.T) {
	file, err := marisfrolg_utils.PathlogExistsFile("./Log")
	fmt.Println(file, err)
}

func TestCreateImg(t *testing.T) {
	var img image.Image
	err := marisfrolg_utils.CreateImg("./Log", "qr.png", img)
	fmt.Println(err)
}

func TestFileToByte(t *testing.T) {
	var f *os.File
	file, err := marisfrolg_utils.FileToByte(f)
	fmt.Println(file, err)
}

package main

import (
	"fmt"

	convert "github.com/lulouis/marisfrolg_utils"
)

func main() {
	substring := convert.Substr("marisfrolg_utils", 0, 10)
	fmt.Println(substring)

}

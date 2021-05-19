package testing

import (
	"fmt"
	"net/url"
	"testing"

	"github.com/lulouis/marisfrolg_utils"
)

func TestHttpPostOnlyBody(t *testing.T) {
	r, err := marisfrolg_utils.HttpPostOnlyBody("www.baidu.com", fmt.Sprintf(`{"data":"123"}`), "123456")
	fmt.Println(string(r), err)
}

func TestHttpGet(t *testing.T) {
	param := url.Values{}
	param.Set("wd", "11")
	r, err := marisfrolg_utils.HttpGet("www.baidu.com", param)
	fmt.Println(string(r), err)
}
func TestHttpGetToUrl(t *testing.T) {
	r, err := marisfrolg_utils.HttpGetToUrl("https://www.baidu.com/s?wd=111")
	fmt.Println(string(r), err)
}

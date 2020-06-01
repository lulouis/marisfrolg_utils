package marisfrolg_utils

import (
	"io/ioutil"
	"net/http"
	"strings"
)

/*
HTTP请求
*/

//只传body
func HttpPostOnlyBody(apiURL string, parmbody string, token string) (rs []byte, err error) {

	req, err := http.NewRequest("POST", apiURL, strings.NewReader(parmbody))
	req.Header.Add("authorization", "Bearer "+token)
	req.Header.Add("cache-control", "no-cache")
	req.Header.Add("Content-Type", "application/json")
	res, err := http.DefaultClient.Do(req)
	body, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	return body, err
}

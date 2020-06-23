package marisfrolg_utils

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

/*
HTTP请求
*/

//post请求:只传body
func HttpPostOnlyBody(apiURL string, parmbody string, token string) (rs []byte, err error) {

	req, err := http.NewRequest("POST", apiURL, strings.NewReader(parmbody))
	if err!=nil{
		return nil,err
	}
	req.Header.Add("authorization", "Bearer "+token)
	req.Header.Add("cache-control", "no-cache")
	req.Header.Add("Content-Type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err!=nil{
		return nil,err
	}
	respData, err := ioutil.ReadAll(res.Body)
	if err!=nil{
		return nil,err
	}
	defer res.Body.Close()
	return respData, err
}

//post请求:带url参数和postBodyData数据
func HttpPostWithReqParamAndToken(reqUrl string, params url.Values, bodyData string, token string) (rs []byte, err error) {
	//url参数转义
	data := params.Encode()
	req, err := http.NewRequest("POST", reqUrl+"?"+data, strings.NewReader(bodyData))
	if err!=nil{
		return nil,err
	}
	req.Header.Add("authorization", "Bearer "+token)
	req.Header.Add("cache-control", "no-cache")
	req.Header.Add("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err!=nil{
		return nil,err
	}
	respData, err := ioutil.ReadAll(res.Body)
	if err!=nil{
		return nil,err
	}
	defer res.Body.Close()
	return respData, err
}
// get 网络请求
func HttpGet(ApiURL string, Params url.Values) (rs []byte, err error) {
	var Url *url.URL
	Url, err = url.Parse(ApiURL)
	if err != nil {
		log.Printf("解析url错误:\r\n%v", err)
		return nil, err
	}
	//如果参数中有中文参数,这个方法会进行URLEncode
	Url.RawQuery = Params.Encode()
	urlstr := Url.String()
	resp, err := http.Get(urlstr)
	//fmt.Println(urlstr)
	if err != nil {
		//fmt.Println("err:", err)
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}
func HttpGetToUrl(ApiURL string) (rs []byte, err error) {
	resp, err := http.Get(ApiURL)
	//fmt.Println(urlstr)
	if err != nil {
		//fmt.Println("err:", err)
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}
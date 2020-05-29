package Marisfrolg_utils

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/go-redis/redis"
	"github.com/tealeg/xlsx"
)

func CheckErr(err error) {
	if err != nil {
		//fmt.Println(err)
		//中断 panic(err)
	}

}
//检查制定路径下是否存在文件如果不存在直接创建文件夹(文件性日志)
func PathlogExistsFile(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	err = os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return false, err
	} else {
		return true, nil
	}
	return false, err
}
func GetRedisByKey(key string, yourDate string,REDIS_CONN string) (list string, err error) {
	client := redis.NewClient(&redis.Options{
		Addr:        REDIS_CONN,
		Password:    "", // no password set
		DB:          0,  // use default DB
		ReadTimeout: 240 * time.Second,
	})
	defer client.Close()

	val, err := client.Get(fmt.Sprintf("%s_%s", key, yourDate)).Result()
	if len(val) > 0 && err == nil {
		return val, nil
	} else {
		return "", errors.New("无缓存")
	}
}

//网络文件下载
func DownloadFile(fileName string, url string) (err error) {
	// Create the file
	out, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// Writer the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

//移除重复数据
func RemoveRepeatedElement(arr []string) (newArr []string) {
	newArr = make([]string, 0)
	sort.Strings(arr)
	for i := 0; i < len(arr); i++ {
		repeat := false
		for j := i + 1; j < len(arr); j++ {
			if arr[i] == arr[j] {
				repeat = true
				break
			}
		}
		if !repeat {
			newArr = append(newArr, arr[i])
		}
	}
	return
}

// 获取文件大小的接口
type Size interface {
	Size() int64
}

func XlsxFileReader(mimeFile multipart.File) (*xlsx.File, error) {

	defer mimeFile.Close()
	var size int64
	if sizeInterface, ok := mimeFile.(Size); ok {
		size = sizeInterface.Size()
	}

	xlFile, err := xlsx.OpenReaderAt(mimeFile, size)
	return xlFile, err
}

func MapToURLValues(data map[string]interface{}) url.Values {
	values := url.Values{}
	for k, v := range data {
		values.Set(k, ValueToString(v))
	}
	return values
}

func ValidateMobile(mobileNum string) bool {
	regMobile := `^1([38][0-9]|4[579]|5[0-3,5-9]|6[6]|7[0135678]|9[89])\d{8}$`
	reg := regexp.MustCompile(regMobile)
	return reg.MatchString(mobileNum)
}


//位数不够自动左补全
//str 传入的字符串 ，totalWidth 补的长度 ，paddingChar 补全的字符
func PadLeft(str string, totalWidth int, paddingChar string) (r string) {
	var result string
	if str == "" {
		return ""
	} else {
		l := len(str)
		result = str
		for i := 0; i < totalWidth-l; i++ {
			result = paddingChar + result
		}
	}
	return result
}

//位数不够自动右补全
// str 传入的字符串 ，totalWidth 补的长度 ，paddingChar 补全的字符
func PadRight(str string, totalWidth int, paddingChar string) (r string) {
	var result string
	if str == "" {
		return ""
	} else {
		l := len(str)
		result = str
		for i := 0; i < totalWidth-l; i++ {
			result = result + paddingChar
		}
	}
	return result
}

//辅助函数
func GetSqlList(s string, ziduanming string) string {
	a := strings.Index(s, ",")
	if a == -1 {
		var condition string
		s = strings.Replace(s, " ", "", -1)
		s = strings.Replace(s, "“", "", -1)
		s = strings.Replace(s, "”", "", -1)
		s = strings.Replace(s, "'", "", -1)
		s = strings.Replace(s, "‘", "", -1)
		s = strings.Replace(s, "，", "", -1)
		s = strings.Replace(s, "[", "", -1)
		s = strings.Replace(s, "]", "", -1)
		condition = " AND " + ziduanming + " = " + "'" + s + "'"
		return condition
	} else {
		var list []string
		var fz = ""
		var condition = ziduanming + "="
		var bz = " AND ("
		s += ","
		s = strings.Replace(s, " ", "", -1)
		s = strings.Replace(s, "\"", "", -1)
		s = strings.Replace(s, "“", "", -1)
		s = strings.Replace(s, "”", "", -1)
		s = strings.Replace(s, "'", "", -1)
		s = strings.Replace(s, "‘", "", -1)
		s = strings.Replace(s, "，", "", -1)
		s = strings.Replace(s, "[", "", -1)
		s = strings.Replace(s, "]", "", -1)
		for _, a := range s {
			if string(a) == "," {
				list = append(list, fz)
				fz = ""
			} else {
				fz = fz + string(a)
			}
		}
		for _, a := range list { //循环加
			bz = bz + condition + "'" + a + "'" + " or "
		}
		content := bz[0 : len(bz)-4]
		content = content + ")"
		return content
	}
}

func StringToRuneArr(s string) []string {
	s = strings.Replace(s, " ", "", -1)
	s = strings.Replace(s, "\"", "", -1)
	s = strings.Replace(s, "“", "", -1)
	s = strings.Replace(s, "”", "", -1)
	s = strings.Replace(s, "'", "", -1)
	s = strings.Replace(s, "‘", "", -1)
	// s = strings.Replace(s, ",", "", -1)
	s = strings.Replace(s, "，", "", -1)
	s = strings.Replace(s, "[", "", -1)
	s = strings.Replace(s, "]", "", -1)
	aa := strings.Split(s, ",")
	return aa
}

//判断是否存在
func Contains(array interface{}, val interface{}) (index int) {
	index = -1
	switch reflect.TypeOf(array).Kind() {
	case reflect.Slice:
		{
			s := reflect.ValueOf(array)
			for i := 0; i < s.Len(); i++ {
				if reflect.DeepEqual(val, s.Index(i).Interface()) {
					index = i
					return
				}
			}
		}
	}
	return
}

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

//SQL IN()的查询里不能超过1000列，将大于1000列的以900为间隔分开组装
//IdList:数据列 例如："'987654321','123456789'"
//Field:字段名
func AssemblyParameters(IdList, Field string) (condition string) {
	var ApplyIdList [100]string
	mbct := strings.Trim(IdList, ",")
	nameArr := strings.Split(mbct, ",")
	if len(nameArr) != 0 {
		for index, a := range nameArr {
			i := index / 900
			ApplyIdList[i] += a + ","
		}
		if ApplyIdList[0] != "" {
			if ApplyIdList[0] != "" && ApplyIdList[1] == "" {
				condition += fmt.Sprintf(" AND (%s IN (%s) ) ", Field, strings.TrimRight(ApplyIdList[0], ","))
			} else {
				for index, itm := range ApplyIdList {
					if itm != "" {
						if index == 0 {
							condition += fmt.Sprintf(" AND (%s IN (%s)  ", Field, strings.TrimRight(itm, ","))
						} else {
							condition += fmt.Sprintf(" OR %s IN (%s) ", Field, strings.TrimRight(itm, ","))
						}
					}
				}
				condition += fmt.Sprintf(")")
			}
		}
		return condition
	} else {
		return ""
	}
}

package Marisfrolg_utils

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

/*
添加操作日志
Type:控制器名称
Title:函数名称
Message:要记录在日志的内容
*/
func AddOperationLog(Type string, Title string, Message string,Filepath string) {
	var (
		Data    = time.Now().Format(`20060102`)
		LogPath = "./"+Filepath+"/" + Data
		exist   bool
		err     error
		path    string
		file    *os.File
	)

	exist, err = PathlogExistsFile(LogPath)
	if err != nil {
		return
	}
	if exist {
		// 以追加模式打开文件，当文件不存在时生成文件
		if Type == "" {
			path = fmt.Sprintf(`%s/%s.log`, LogPath, Data)
		} else {

			path = fmt.Sprintf(`%s/%s.log`, LogPath, Data+"-"+Type)
		}
		file, err = os.OpenFile(path, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
		if err != nil {
			return
		}
		defer file.Close()
		S := "        " + "\r\n" + "日期：[" + time.Now().Format(`2006-01-02 15:04:05`) + "]" + "    " + "IP：" + GetIP() + "\r\n" +
			"标题：" + Title + "\r\n" +
			"内容：" + Message + "\r\n"
		n, err := io.WriteString(file, S)
		if err != nil {
			log.Println(n, err)
		}

	}
}

/*
添加操作日志
Type:控制器名称
Title:函数名称
Message:要记录在日志的内容
*/
func AddOperationLogFromFileName(action, message, fileName string,Filepath string) {
	var (
		exist bool
		err   error
		path  string
		file  *os.File
	)
	nowDate := time.Now().Format(`20060102`)
	LogPath := "./"+Filepath+"/" + nowDate
	exist, err = PathlogExistsFile(LogPath)
	if err != nil {
		return
	}
	if exist {
		// 以追加模式打开文件，当文件不存在时生成文件
		path = fmt.Sprintf(`%s/%s.log`, LogPath, action+"_"+fileName)
		file, err = os.OpenFile(path, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
		if err != nil {
			return
		}
		defer file.Close()
		S := "        " + "\r\n" + "日期：[" + time.Now().Format(`2006-01-02 15:04:05`) + "]" + "    " + "IP：" + GetIP() + "\r\n" +
			"内容：" + message + "\r\n"
		n, err := io.WriteString(file, S)
		if err != nil {
			log.Println(n, err)
		}

	}
}
//创建.log日志
func CreateLog(Path string, fileName string, Data string,Filepath string) {
	LogPath := "./"+Filepath+"/" + time.Now().Format(`20060102`) + "/" + Path
	exist, err := PathlogExistsFile(LogPath)
	if err != nil {
		return
	}
	if exist {
		// 以追加模式打开文件，当文件不存在时生成文件
		path := fmt.Sprintf(`%s/%s.log`, LogPath, fileName)
		file, error := os.OpenFile(path, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
		if error != nil {
			return
		}
		defer file.Close()
		n, err := io.WriteString(file, Data)
		if err != nil {
			log.Println(n, err)
		}
	}
}



func GetIP() string {
	addrs, err := net.InterfaceAddrs()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, address := range addrs {

		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}

		}
	}
	return ""
}

///获取当前函数名
func PrintMyName() string {
	pc, _, _, _ := runtime.Caller(2)
	var funcname = runtime.FuncForPC(pc).Name()
	funcname = filepath.Ext(funcname)
	return strings.TrimPrefix(funcname, ".")
}

///获取当前文件名称
func GetFileName() string {
	_, file, _, _ := runtime.Caller(1)
	filenameWithSuffix := path.Base(file)                              //获取文件名带后缀
	fileSuffix := path.Ext(filenameWithSuffix)                         //获取文件后缀
	filenameOnly := strings.TrimSuffix(filenameWithSuffix, fileSuffix) //获取文件名
	return filenameOnly
}

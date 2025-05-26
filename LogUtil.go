package marisfrolg_utils

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path"
	"path/filepath"
	"reflect"
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
func AddOperationLog(Type string, Title string, Message string) {
	var (
		Data    = time.Now().Format(`20060102`)
		LogPath = "./logs/" + Data
		exist   bool
		err     error
		path    string
		file    *os.File
	)

	exist, err = PathlogExists(LogPath)
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
		S := "        " + "\r\n" + "日期:[" + time.Now().In(CstZone).Format("2006-01-02 15:04:05") + "]" + "    " + "IP:" + GetIP() + "\r\n" +
			"标题:" + Title + "\r\n" +
			"内容:" + Message + "\r\n"
		n, err := io.WriteString(file, S)
		if err != nil {
			log.Println(n, err)
		}

	}
}

// 检查制定路径下是否存在文件如果不存在直接创建文件夹(文件性日志)
func PathlogExists(path string) (bool, error) {
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

// /获取当前函数名
func PrintMyName() string {
	pc, _, _, _ := runtime.Caller(2)
	var funcname = runtime.FuncForPC(pc).Name()
	funcname = filepath.Ext(funcname)
	return strings.TrimPrefix(funcname, ".")
}

// /获取当前文件名称
func GetFileName() string {
	_, file, _, _ := runtime.Caller(1)
	filenameWithSuffix := path.Base(file)                              //获取文件名带后缀
	fileSuffix := path.Ext(filenameWithSuffix)                         //获取文件后缀
	filenameOnly := strings.TrimSuffix(filenameWithSuffix, fileSuffix) //获取文件名
	return filenameOnly
}

/*
添加操作日志
FileName:文件名称
MethodName:方法名称
UserNo:员工工号
Message:日志内容，方法参数
*/
func AddLog(FileName string, MethodName string, UserNo string, Message string) {
	var (
		Data    = time.Now().In(CstZone).Format(`20060102`)
		LogPath = "./logs/" + Data
		exist   bool
		err     error
		path    string
		file    *os.File
	)

	exist, err = PathlogExists(LogPath)
	if err != nil {
		return
	}
	if exist {
		// 以追加模式打开文件，当文件不存在时生成文件
		if FileName == "" {
			path = fmt.Sprintf(`%s/%s.log`, LogPath, Data)
		} else {

			path = fmt.Sprintf(`%s/%s.log`, LogPath, Data+"-"+FileName)
		}
		file, err = os.OpenFile(path, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
		if err != nil {
			return
		}
		defer file.Close()
		S := "        " + "\r\n" + "日期:[" + time.Now().In(CstZone).Format(`2006-01-02 15:04:05`) + "]" + "    " + "员工:" + UserNo + "    " + "IP:" + GetIP() + "\r\n" +
			"方法:" + MethodName + " \n" + Message + "\r\n"
		n, err := io.WriteString(file, S)
		if err != nil {
			log.Println(n, err)
		}

	}
}

// 获取当前执行函数名
func GetMethodName() (funcName string) {
	pc := make([]uintptr, 1)
	runtime.Callers(2, pc)
	f := runtime.FuncForPC(pc[0])
	funcNameList := strings.Split(f.Name(), ".")
	funcName = funcNameList[len(funcNameList)-1]
	return
}

type AuditLogInfo struct {
	FileName       string //文件名称
	MethodName     string //方法名称
	UserCode       string //工号
	Msg            string //日志详细内容
	RequestMsg     string //请求内容
	ResponseMsg    string //响应内容
	SapRequestMsg  string //SAP请求内容
	SapResponseMsg string //SAP响应内容
	Requests       map[string]string
	Responses      map[string]string
}

func NewAuditLogInfo(fileName string, methodName string, userCode string) *AuditLogInfo {
	return &AuditLogInfo{
		FileName:   fileName,
		MethodName: methodName,
		UserCode:   userCode,
		Requests:   make(map[string]string),
		Responses:  make(map[string]string),
	}
}
func (logInfo *AuditLogInfo) AppendMsg(msg string) {
	logInfo.Msg += msg
}
func (logInfo *AuditLogInfo) SetRequest(request interface{}) {
	requestBytes, _ := json.Marshal(request)
	logInfo.RequestMsg = string(requestBytes)
}
func (logInfo *AuditLogInfo) SetResponse(response interface{}) {
	responseBytes, _ := json.Marshal(response)
	logInfo.ResponseMsg = string(responseBytes)
}
func (logInfo *AuditLogInfo) SetSapRequest(sapRequest interface{}) {
	sapRequestBytes, _ := json.Marshal(sapRequest)
	logInfo.SapRequestMsg = string(sapRequestBytes)
}
func (logInfo *AuditLogInfo) SetSapResponse(sapResponse interface{}) {
	sapResponseBytes, _ := json.Marshal(sapResponse)
	logInfo.SapResponseMsg = string(sapResponseBytes)
}
func (logInfo *AuditLogInfo) AppendItem(name string, logItem interface{}) {
	xType := reflect.TypeOf(logItem)
	// 检查类型是否为string
	if xType.Kind() == reflect.String {
		logInfo.Requests[name] = logItem.(string)
	} else {
		requestBytes, _ := json.Marshal(logItem)
		logInfo.Requests[name] = string(requestBytes)
	}
}

func AuditLog(logInfo *AuditLogInfo, err error) {
	var logMsg string
	if logInfo.RequestMsg != "" {
		logMsg += fmt.Sprintf("请求参数:\n%s\n", logInfo.RequestMsg)
	}
	if logInfo.ResponseMsg != "" {
		logMsg += fmt.Sprintf("响应结果:\n%s\n", logInfo.ResponseMsg)
	}
	if logInfo.SapRequestMsg != "" {
		logMsg += fmt.Sprintf("SAP请求:\n%s\n", logInfo.SapRequestMsg)
	}
	if logInfo.SapResponseMsg != "" {
		logMsg += fmt.Sprintf("SAP结果:\n%s\n", logInfo.SapResponseMsg)
	}
	if len(logInfo.Requests) > 0 {
		logMsg += "详细日志:\n"
		for k, v := range logInfo.Requests {
			logMsg += fmt.Sprintf("%s:%s\n", k, v)
		}
	}
	if len(logInfo.Responses) > 0 {
		logMsg += "响应列表:\n"
		for k, v := range logInfo.Responses {
			logMsg += fmt.Sprintf("%s:%s\n", k, v)
		}
	}
	if logInfo.Msg != "" {
		logMsg += fmt.Sprintf("补充记录:\n%s\n", logInfo.Msg)
	}
	if err != nil {
		logMsg += fmt.Sprintf("方法%s,执行错误:\n%s\n", logInfo.MethodName, err.Error())
	} else {
		logMsg += fmt.Sprintf("方法%s,执行成功\n", logInfo.MethodName)
	}
	AddLog(logInfo.FileName, logInfo.MethodName, logInfo.UserCode, logMsg)
}

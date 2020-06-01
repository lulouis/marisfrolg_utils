package main

import (
	"fmt"
	"github.com/lulouis/marisfrolg_utils"
	"testing"
)

func TestSendMail(t *testing.T) {
	var (
		mailTo []string
		FileNameList    []string
	)
	mailTo=append(mailTo,"test@qq.com")
	emailServer := new(marisfrolg_utils.EmailServer)
	emailServer.User = "test@qq.com"
	emailServer.PassWord = "123456"
	emailServer.Host = "smtp.exmail.qq.com"
	emailServer.Port = 465
	emailRequest := *new(marisfrolg_utils.EmailRequest)
	emailRequest.MailTo = mailTo
	emailRequest.SetHeader = "测试"
	emailRequest.Subject = "主题"
	emailRequest.Body = "正文"
	emailRequest.FileUrl = FileNameList
	var r=emailServer.SendMail(emailRequest)
	fmt.Println(r)
}

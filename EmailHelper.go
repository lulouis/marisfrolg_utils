package Marisfrolg_utils

import (
	"fmt"
	_ "io/ioutil"
	"mime"
	"os"
	"path"
	_ "runtime"
	_ "strconv"

	"gopkg.in/gomail.v2"
)

type Attachment struct {
	name        []string
	contentType string
	withFile    bool
}

//定义邮箱服务器连接信息，如果是阿里邮箱 pass填密码，qq邮箱填授权码
type EmailServer struct {
	User     string //邮箱
	PassWord string //密码（QQ邮箱为授权码）
	Host     string
	Port     int //端口号
}

type EmailRequest struct {
	MailTo     []string //收件人可以是多个
	SetHeader  string   //标题
	Subject    string   //邮件主题
	Body       string   //邮件正文
	FileUrl    []string //附件路径可以是多个
	attachment Attachment
}

/*
发送邮件
*/
func (emailServer EmailServer) SendMail(emailRequest EmailRequest) error {
	mail := gomail.NewMessage()
	mail.SetHeader("From", mail.FormatAddress(emailServer.User, emailRequest.SetHeader))
	mail.SetHeader("To", emailRequest.MailTo...)    //发送给多个用户
	mail.SetHeader("Subject", emailRequest.Subject) //设置邮件主题
	mail.SetBody("text/html", emailRequest.Body)    //设置邮件正文
	if len(emailRequest.FileUrl) > 0 {
		for i := 0; i < len(emailRequest.FileUrl); i++ {
			fileSize := getFileSize(emailRequest.FileUrl[i])
			if fileSize == 0 || fileSize > 10000000 {
				error := fmt.Errorf(`附件大小异常！名称:%s ,大小为:%d`, emailRequest.FileUrl[i], fileSize)
				return error
			}
			mail.Attach(emailRequest.FileUrl[i],
				gomail.Rename(path.Base(emailRequest.FileUrl[i])),
				gomail.SetHeader(map[string][]string{
					"Content-Disposition": {
						fmt.Sprintf(`attachment; filename="%s"`, mime.BEncoding.Encode("UTF-8", path.Base(emailRequest.FileUrl[i]))),
					},
				}),
			)
		}
	}
	send := gomail.NewDialer(emailServer.Host, emailServer.Port, emailServer.User, emailServer.PassWord)
	err := send.DialAndSend(mail)
	return err
}

//判断文件大小
func getFileSize(path string) int64 {
	if !exists(path) {
		return 0
	}
	fileInfo, err := os.Stat(path)
	if err != nil {
		return 0
	}
	return fileInfo.Size()
}

//判断是否存在文件
func exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

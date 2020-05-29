package Marisfrolg_utils

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"
)

//判断目录是否存在
func IsDirExists(path string) bool {
	fi, err := os.Stat(path)
	if err != nil {
		return os.IsExist(err)
	} else {
		return fi.IsDir()
	}
}

//是否创建目录
func IsCreateDir(path string) (err error) {
	if !IsDirExists(path) {
		err = os.MkdirAll(path, 0766)
	}
	return err
}

func SaveFileToTempDirectory(isNeedPrefix bool, file *multipart.FileHeader,UploadFile string) (fileName, filePath string, err error) {
	var (
		dir       string
		existsDir bool
		out       *os.File
		reader    multipart.File
	)
	dir, _ = filepath.Abs(filepath.Dir(os.Args[0]))
	filePath = strings.Replace(dir, "\\", "/", -1) + "/"+UploadFile
	existsDir, _ = PathlogExistsFile(filePath)
	if !existsDir {
		os.Mkdir(filePath, os.ModePerm)
	}
	if isNeedPrefix {
		fileName = time.Now().Format("20060102150405") + "_" + file.Filename
	} else {
		fileName = file.Filename
	}

	filePath = filePath + "/" + fileName

	if out, err = os.Create(filePath); err != nil {
		goto ERR
	}
	defer out.Close()
	if reader, err = file.Open(); err != nil {
		goto ERR
	}
	defer reader.Close()
	_, err = io.Copy(out, reader)
	return
ERR:
	return
}

/*
  创建人：李奇峰
  功能：读取文件反序列化
*/
func LoadFile(filename string, v interface{}) (err error) {
	var data []byte
	if data, err = ioutil.ReadFile(filename); err != nil {
		return
	}

	if err = json.Unmarshal(data, v); err != nil {
		return
	}
	return
}

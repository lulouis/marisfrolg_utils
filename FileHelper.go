package marisfrolg_utils

import (
	"encoding/json"
	"fmt"
	"github.com/tealeg/xlsx"
	"image"
	"image/png"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

/* 文件相关操作*/

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

//将文件保存到指定目录
func SaveFileToTempDirectory(isNeedPrefix bool, file *multipart.FileHeader, UploadFile string) (fileName, filePath string, err error) {
	var (
		dir       string
		existsDir bool
		out       *os.File
		reader    multipart.File
	)
	dir, _ = filepath.Abs(filepath.Dir(os.Args[0]))
	filePath = strings.Replace(dir, "\\", "/", -1) + "/" + UploadFile
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

//检查制定路径下是否存在文件如果不存在直接创建文件夹
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

// 获取文件大小的接口
type Size interface {
	Size() int64
}

//读取Excel文件流
func XlsxFileReader(mimeFile multipart.File) (*xlsx.File, error) {

	defer mimeFile.Close()
	var size int64
	if sizeInterface, ok := mimeFile.(Size); ok {
		size = sizeInterface.Size()
	}

	xlFile, err := xlsx.OpenReaderAt(mimeFile, size)
	return xlFile, err
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
//创建图片
func CreateImg(filename string, img image.Image) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	err = png.Encode(file, img)
	if err != nil {
		return err
	}
	file.Close()
	return nil
}

//将文件转换成byte[]
func FileToByte(f *os.File) ([]byte,error) {
	var payload []byte
	for {
		buf := make([]byte, 1024)
		switch nr, err := f.Read(buf[:]); true {
		case nr < 0:
			return nil,err
			os.Exit(1)
		case nr == 0: // EOF
			return payload,nil
		case nr > 0:
			payload = append(payload, buf...)
		}
	}
}

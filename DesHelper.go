package Marisfrolg_utils

import (
	"crypto/cipher"
	"crypto/des"
	"encoding/base64"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

//身份证加密
func FileEncrypt(Path []string) (err error) {
	key := []byte("12345678")
	key[0] = 4
	key[1] = 122
	key[2] = 38
	key[3] = 13
	key[4] = 100
	key[5] = 80
	key[6] = 92
	key[7] = 249
	for _, v1 := range Path {
		dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
		ServerPath := strings.Replace(dir, "\\", "/", -1) + v1
		fileByte, _ := ioutil.ReadFile(ServerPath)
		result, _ := DesEncrypt(fileByte, key)
		err = ioutil.WriteFile(ServerPath, result, 0644)
	}

	return
}

//身份证本地加密
func ClientFileEncrypt(Path []string) (err error) {
	key := []byte("12345678")
	key[0] = 4
	key[1] = 122
	key[2] = 38
	key[3] = 13
	key[4] = 100
	key[5] = 80
	key[6] = 92
	key[7] = 249
	for _, v1 := range Path {
		//dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
		//ServerPath := strings.Replace(dir, "\\", "/", -1) + v1
		fileByte, _ := ioutil.ReadFile(v1)
		result, _ := DesEncrypt(fileByte, key)
		err = ioutil.WriteFile(v1, result, 0644)
	}

	return
}

//身份证解密
func FileDecrypt(Path []string) (Base64File []string, err error) {
	key := []byte("12345678")
	key[0] = 4
	key[1] = 122
	key[2] = 38
	key[3] = 13
	key[4] = 100
	key[5] = 80
	key[6] = 92
	key[7] = 249
	for _, v1 := range Path {
		dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
		ServerPath := strings.Replace(dir, "\\", "/", -1) + v1
		fileByte, _ := ioutil.ReadFile(ServerPath)
		result, _ := DesDecrypt(fileByte, key)
		strpsw := base64.StdEncoding.EncodeToString(result)
		Base64File = append(Base64File, strpsw)
	}

	return
}

//CBC加密
func DesEncrypt(origData, key []byte) ([]byte, error) {
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}
	IV := []byte("12345678")
	IV[0] = 174
	IV[1] = 141
	IV[2] = 59
	IV[3] = 12
	IV[4] = 153
	IV[5] = 210
	IV[6] = 126
	IV[7] = 217
	origData = PKCS5Padding(origData, block.BlockSize())
	//origData = ZeroPadding(origData, block.BlockSize())
	blockMode := cipher.NewCBCEncrypter(block, IV)
	crypted := make([]byte, len(origData))
	// 根据CryptBlocks方法的说明，如下方式初始化crypted也可以
	// crypted := origData
	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

// func Pkcs5Padding(ciphertext []byte, blockSize int) []byte {
// 	padding := blockSize - len(ciphertext)%blockSize
// 	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
// 	return append(ciphertext, padtext...)
// }

//CBC解密
func DesDecrypt(crypted, key []byte) ([]byte, error) {
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}
	IV := []byte("12345678")
	IV[0] = 174
	IV[1] = 141
	IV[2] = 59
	IV[3] = 12
	IV[4] = 153
	IV[5] = 210
	IV[6] = 126
	IV[7] = 217
	blockMode := cipher.NewCBCDecrypter(block, IV)
	origData := make([]byte, len(crypted))
	// origData := crypted
	blockMode.CryptBlocks(origData, crypted)
	origData = PKCS5UnPadding1(origData)
	// origData = ZeroUnPadding(origData)
	return origData, nil
}

func PKCS5UnPadding1(origData []byte) []byte {
	length := len(origData)
	// 去掉最后一个字节 unpadding 次
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

type Stream interface {
	// XORKeyStream XORs each byte in the given slice with a byte from the
	// cipher's key stream. Dst and src may point to the same memory.
	XORKeyStream(dst, src []byte)
}

//ECB加密
func DesECBEncrypt(data, key []byte) ([]byte, error) {
	//NewCipher创建一个新的加密块
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}

	bs := block.BlockSize()
	data = PKCS5Padding(data, bs)
	if len(data)%bs != 0 {
		return nil, errors.New("need a multiple of the blocksize")
	}

	out := make([]byte, len(data))
	dst := out
	for len(data) > 0 {
		//Encrypt加密第一个块，将其结果保存到dst
		block.Encrypt(dst, data[:bs])
		data = data[bs:]
		dst = dst[bs:]
	}
	return out, nil
}

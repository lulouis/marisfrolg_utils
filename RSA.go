package Marisfrolg_utils

import (
	"bytes"
	"crypto"
	"crypto/cipher"
	"crypto/des"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"
)

//LHH
type PriKeyType uint

const (
	PKCS1 PriKeyType = iota
	PKCS8
)

// 可通过openssl产生
//openssl genrsa -out rsa_private_key.pem 1024
var privateKey = []byte(`  
-----BEGIN RSA PRIVATE KEY-----
MIICdwIBADANBgkqhkiG9w0BAQEFAASCAmEwggJdAgEAAoGBAMKQYwXKYC9csfyd
i29JmpLFf1tHKlfftsv8nUGhmxCm55Wue1phh8gl8ab8DGRw1zAH1+MDNdwNDFCX
IvnWh9nRiN6pyYFf9mqKt7QYQdZ0uRYnj+JgJ+iYhJGY8AFcSYcDmZkHV8E/QCDz
0HPb7ulr6Z9p8xKHZ7iIYhDw10FtAgMBAAECgYEAn0Dl/Jhg0IOUIEyoE+hwQFCt
5P3EN/civadA5LatoRyslEUkLJ+GL5p3SRIn5pLCYEsbN3KqRDrd6J09ALjLqlwK
ZKXbi4hsPgTSdd/bWiAIdf+hqr8vWGvCHPIVWnkQcNBTLlKtFc7JNgIKzbCZRxQn
WtUPHH+k8hc+Ob6OZwECQQD5gNWCB24n0ZTI6xUYJS+woLAbXlKGX8+ANAg4X24m
qkWxI6PmKVZej8aUMiE6YtNVtdDtWrBYAH5V9IG+tPytAkEAx6FVMOahA4aS/NNS
paR8OByWqtN+VRP+X9p8Ytb2sz+P0CdE14ThKhARK428ToYZPTkrs66cNm5CGB8g
dyUvwQJBAJCx5aCGHJ0dD1NB+jbJghHF7rvAhM2HDPiFtGq09VWZE9e6GpglSwCG
ExzowZpxq6weSC8OlAxFJP9GUGQ/4/UCQDk5+3ToODINiudlIOUREPb44wwXUrjK
4XnS5SNkYhYiW3SdPTPXCMEJGBL3L4sHEAcn82ov3OIRm2rUyXa+N0ECQCBcQwku
3IIin7n4mJqPue4Y80I4WoYYrfjLOD96D0YCQFz9q4sLsu6WwzY1bNuzxqyqRdLN
x2lkD+DZuKz5BTY=
-----END RSA PRIVATE KEY-----
`)

//openssl
//openssl rsa -in rsa_private_key.pem -pubout -out rsa_public_key.pem
var publicKey = []byte(`  
-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDCkGMFymAvXLH8nYtvSZqSxX9b
RypX37bL/J1BoZsQpueVrntaYYfIJfGm/AxkcNcwB9fjAzXcDQxQlyL51ofZ0Yje
qcmBX/Zqire0GEHWdLkWJ4/iYCfomISRmPABXEmHA5mZB1fBP0Ag89Bz2+7pa+mf
afMSh2e4iGIQ8NdBbQIDAQAB
-----END PUBLIC KEY-----    
`)

//LHH
//私钥签名
func RSASign(data, privateKey []byte, keyType PriKeyType) ([]byte, error) {
	h := sha256.New()
	h.Write(data)
	hashed := h.Sum(nil)
	priv, err := RSAgetPriKey(privateKey, keyType)
	if err != nil {
		return nil, err
	}
	return rsa.SignPKCS1v15(rand.Reader, priv, crypto.SHA256, hashed)
}

//公钥验证
func RSASignVer(data, signature, publicKey []byte) error {
	hashed := sha256.Sum256(data)
	//获取公钥
	pub, err := RSAgetPubKey(publicKey)
	if err != nil {
		return err
	}
	//验证签名
	return rsa.VerifyPKCS1v15(pub, crypto.SHA256, hashed[:], signature)
}

// 公钥加密
func RSAEncrypt(data []byte) ([]byte, error) {
	//获取公钥
	pub, err := RSAgetPubKey(publicKey)
	if err != nil {
		return nil, err
	}
	//加密
	return rsa.EncryptPKCS1v15(rand.Reader, pub, data)
}

// 私钥解密,privateKey为pem文件里的字符
func RSADecrypt(encData []byte, keyType PriKeyType) ([]byte, error) {
	//解析PKCS1a或者PKCS8格式的私钥
	priv, err := RSAgetPriKey(privateKey, keyType)
	if err != nil {
		return nil, err
	}
	// 解密
	return rsa.DecryptPKCS1v15(rand.Reader, priv, encData)
}

func RSAgetPubKey(publicKey []byte) (*rsa.PublicKey, error) {
	//解密pem格式的公钥
	block, _ := pem.Decode(publicKey)
	if block == nil {
		return nil, errors.New("public key error")
	}
	// 解析公钥
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	// 类型断言
	if pub, ok := pubInterface.(*rsa.PublicKey); ok {
		return pub, nil
	} else {
		return nil, errors.New("public key error")
	}
}

func RSAgetPriKey(privateKey []byte, keyType PriKeyType) (*rsa.PrivateKey, error) {
	//获取私钥
	block, _ := pem.Decode(privateKey)
	if block == nil {
		return nil, errors.New("private key error!")
	}
	var priKey *rsa.PrivateKey
	var err error
	switch keyType {
	case PKCS1:
		{
			priKey, err = x509.ParsePKCS1PrivateKey(block.Bytes)
			if err != nil {
				return nil, err
			}
		}
	case PKCS8:
		{
			prkI, err := x509.ParsePKCS8PrivateKey(block.Bytes)
			if err != nil {
				return nil, err
			}
			priKey = prkI.(*rsa.PrivateKey)
		}
	default:
		{
			return nil, errors.New("unsupport private key type")
		}
	}
	return priKey, nil
}

//LHH
func Encrypt(password, deskey string) (reuslt string, err error) {
	origData := []byte(password)
	key := []byte(deskey)
	block, err := des.NewCipher(key)
	if err != nil {
		return "error", err
	}
	origData = PKCS5Padding(origData, block.BlockSize())
	blockMode := cipher.NewCBCEncrypter(block, key)
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)

	reuslt = base64.StdEncoding.EncodeToString(crypted) //一般64位编码处理
	reuslt = fmt.Sprintf("%X", crypted)                 //ODS的处理，16进制编码
	return reuslt, nil
}

func Decrypt(password, deskey string) ([]byte, error) {
	//针对ODS的反解处理，16进制字符串编码变为二进制
	crypted, _ := hex.DecodeString(password)
	// crypted := []byte(password)

	key := []byte(deskey)
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockMode := cipher.NewCBCDecrypter(block, key)
	origData := make([]byte, len(crypted))
	// origData := crypted
	blockMode.CryptBlocks(origData, crypted)
	origData = PKCS5UnPadding(origData)
	// origData = ZeroUnPadding(origData)
	return origData, nil
}

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	// 去掉最后一个字节 unpadding 次
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

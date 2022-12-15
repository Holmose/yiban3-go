package helper

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func Md5(sig string) string {
	newSig := md5.Sum([]byte(sig)) //转成加密编码
	// 将编码转换为字符串
	newArr := fmt.Sprintf("%x", newSig)
	//输出字符串字母都是小写，转换为大写
	return strings.ToTitle(newArr)
}

func Rand() string {
	b := make([]byte, 5)
	_, err := rand.Read(b)
	if err != nil {
		Rand()
	}
	return string(b)
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

const AesKey string = "C69E7046C69E7046"

func AesEncrypt(orig string) string {
	// 转成字节数组
	origData := []byte(orig)
	k := []byte(AesKey)

	// 分组秘钥
	block, _ := aes.NewCipher(k)
	// 获取秘钥块的长度
	blockSize := block.BlockSize()
	// 补全码
	origData = PKCS7Padding(origData, blockSize)
	// 加密模式
	blockMode := cipher.NewCBCEncrypter(block, k[:blockSize])
	// 创建数组
	cryted := make([]byte, len(origData))
	// 加密
	blockMode.CryptBlocks(cryted, origData)

	return base64.StdEncoding.EncodeToString(cryted)

}

func AesDecrypt(cryted string) string {
	// 转成字节数组
	crytedByte, _ := base64.StdEncoding.DecodeString(cryted)
	k := []byte(AesKey)

	// 分组秘钥
	block, _ := aes.NewCipher(k)
	// 获取秘钥块的长度
	blockSize := block.BlockSize()
	// 加密模式
	blockMode := cipher.NewCBCDecrypter(block, k[:blockSize])
	// 创建数组
	orig := make([]byte, len(crytedByte))
	// 解密
	blockMode.CryptBlocks(orig, crytedByte)
	// 去补全码
	orig = PKCS7UnPadding(orig)
	return string(orig)
}

// PKCS7Padding 补码
func PKCS7Padding(data []byte, blockSize int) []byte {
	//判断缺少几位长度。最少1，最多 blockSize
	padding := blockSize - len(data)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padText...)
}

// PKCS7UnPadding 去码
func PKCS7UnPadding(data []byte) []byte {
	length := len(data)
	//获取填充的个数
	unPadding := int(data[length-1])
	return data[:(length - unPadding)]
}

func Division(a string, b string) float64 {
	i1, _ := strconv.Atoi(a)
	i2, _ := strconv.Atoi(b)
	num1 := float64(i1) / float64(i2)
	return num1
}

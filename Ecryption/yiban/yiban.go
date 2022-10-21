package yiban

import (
	"Yiban3/Ecryption/aes"
	"Yiban3/Ecryption/rsa"
	"encoding/base64"
)

/*
	易班相关的加密解密工具
*/

type Params struct {
	WFId   string
	Data   interface{}
	Extend interface{}
}

// KEY 秘钥
var securityKey16 = "2knV5VGRTScU7pOq"

// IV 偏移量
var iv = "UmNWaNtM0PUdtFCs"

var Aes = aes.AesTool(securityKey16, iv)

// FormDecrypt 先进行Base64解密 然后对其解密Aes
func FormDecrypt(decryptData string) (string, error) {
	decodeString, err := base64.StdEncoding.DecodeString(decryptData)
	if err != nil {
		return "", err
	}
	decrypt, err := Aes.Decrypt(string(decodeString))
	if err != nil {
		return "", err
	}
	return decrypt, nil
}

// FormEncrypt 先进行Aes加密，然后对其进行Base64加密
func FormEncrypt(encryptData string) (string, error) {
	encrypt, err := Aes.Encrypt(encryptData)
	if err != nil {
		return "", err
	}
	toString := base64.StdEncoding.EncodeToString([]byte(encrypt))
	return toString, nil
}

// LoginEncrypt 先进行Rsa加密，然后使用Base64加密
func LoginEncrypt(pwd []byte, pubKey []byte) string {
	encrypt := rsa.RSA_Encrypt(pwd, pubKey)
	toString := base64.StdEncoding.EncodeToString(encrypt)
	return toString
}

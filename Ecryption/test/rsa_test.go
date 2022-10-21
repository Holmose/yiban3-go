package test

import (
	"Yiban3/Ecryption/rsa"
	"fmt"
	"os"
	"testing"
)

func TestRSA(t *testing.T) {
	//生成密钥对，保存到文件
	rsa.GenerateRSAKey(2048)
	message := []byte("hello world")
	//加密
	file, _ := os.ReadFile("pem/public.pem")
	cipherText := rsa.RSA_Encrypt(message, file)
	fmt.Println("加密后为：", string(cipherText))
	//解密
	plainText := rsa.RSA_Decrypt(cipherText, "pem/private.pem")
	fmt.Println("解密后为：", string(plainText))
}

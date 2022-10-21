package aes

import (
	"bytes"
	basicAES "crypto/aes"
	"crypto/cipher"
	"encoding/base64"
)

/**
 * run test
 */

type aes struct {
	securityKey []byte
	iv          []byte
}

/**
 * constructor
 */
func AesTool(securityKey string, iv string) *aes {
	return &aes{[]byte(securityKey), []byte(iv)}
}

/**
 * 加密
 * @param string $plainText 明文
 * @return bool|string
 */
func (a aes) Encrypt(plainText string) (string, error) {
	block, err := basicAES.NewCipher(a.securityKey)
	if err != nil {
		return "", err
	}
	plainTextByte := []byte(plainText)
	blockSize := block.BlockSize()
	plainTextByte = addPKCS7Padding(plainTextByte, blockSize)
	cipherText := make([]byte, len(plainTextByte))
	mode := cipher.NewCBCEncrypter(block, a.iv)
	mode.CryptBlocks(cipherText, plainTextByte)
	return base64.StdEncoding.EncodeToString(cipherText), nil
}

/**
 * 解密
 * @param string $cipherText 密文
 * @return bool|string
 */
func (a aes) Decrypt(cipherText string) (string, error) {
	block, err := basicAES.NewCipher(a.securityKey)
	if err != nil {
		return "", err
	}
	cipherDecodeText, decodeErr := base64.StdEncoding.DecodeString(cipherText)
	if decodeErr != nil {
		return "", decodeErr
	}
	mode := cipher.NewCBCDecrypter(block, a.iv)
	originCipherText := make([]byte, len(cipherDecodeText))
	mode.CryptBlocks(originCipherText, cipherDecodeText)
	originCipherText = stripPKSC7Padding(originCipherText)
	return string(originCipherText), nil
}

/**
 * 填充算法
 * @param string $source
 * @return string
 */
func addPKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	paddingText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, paddingText...)
}

/**
 * 移去填充算法
 * @param string $source
 * @return string
 */
func stripPKSC7Padding(cipherText []byte) []byte {
	length := len(cipherText)
	unpadding := int(cipherText[length-1])
	return cipherText[:(length - unpadding)]
}

package test

import (
	aes2 "Yiban3/Ecryption/aes"
	"fmt"
	"testing"
)

func TestAes(t *testing.T) {

	securityKey16 := "2knV5VGRTScU7pOq"
	iv := "UmNWaNtM0PUdtFCs"
	aes := aes2.AesTool(securityKey16, iv)

	// 解密
	word := "KMi96YvAD7uVY2qN6O/uAZmGxZJ3iH0FVWmgZWKlz6qtrXWxWGRahZ7B2c4Typ/s9pxLJrSLQ5231Kn/VqH6MkZ8bCdDouoGrCn+Q4nhpk+ATHcyl82TIbRHG5de9Xyu0FI3DLyc0IPx7hB4FP0hKIliaGUJf67KsXeYtciIS7vh0BK51NZCET7Tqw0nmxCfqZNZtDUgDQclkFE4J0QC6v7b2Yco5lo0NbluuTTNuhNwbNZr2VlH3FLuzV4XZUXJT2eHCLR17o5GZ+z4pJrrF2REael3lJ+lBo6c37ds1TaGifBTzJP80/8QwmzaKYxqRdxqPenlNKiX7RJ8kgfY1rO80G9PzRRBCrYuMblt9tsqddusmDscqv+HtjqIGPIQN+Q5WIIGob52OoxVCKkakMxlO3Krrzw9lVR5n05DF4chR4UEzWIoZXmnDxuUS5ZLRVhfpsKEJm+7MgoP77WlCc33MmBlDN2QVFJOiR+N4/IXl0scL73S7Tz0PhxO8u7R9LXwTxpPjsFhsLBagO6eQe2ug2YPc9LwHXbikcHdDLJpQrEUGNTBSScCsMQ7p5qpKES2VdkTCJExeRXbgjvtcIdthoRGrvpkf6nnRtsCRBaISUW1wdQxEC/leL+myX/A8ryOIJKzBhsmvSQyqLPJjldC35p9581hvFhg+A8zLkbT/GLWUKADRGkgDXJqWBQqFl9qRNBeD6nvouqWTbjJqhApro4pVMBS4LZbLXuC0XLIJqogjQPCQ3T5blrgIo2vKjOCX+L8Gzx68TkhhDhUJbYFnpdCKJoKPLcsNIXuRNJze3rzBPmV5YkFKFivUUGZiTgEcyXqhyH1w1DeCVV5dx/TTQVGjzp+/UBSOyUoXIZDK5HHIWtcSwjMRNzazxpr8Io63Y1by76MOoTLkpXeVerL0UmJYOBdsPy1GRZs07nHRxGP6L+sCzzSgtX9P7SftZ4hOTNLSOLdtQYxz+xRyzzi9AQhmx8ncmjRHKf3zkl7VGg6uFJygDCe1Y6OYJlZ7RNqAAuoztz6y3nP+n5+zFZq3I9PM59w0E7t/PqHbPo="
	outPlainText, _ := aes.Decrypt(word)
	fmt.Println("解密后明文：" + outPlainText)

	// 加密
	plainText, _ := aes.Encrypt(outPlainText)
	fmt.Println("加密后的密文：" + plainText)

	if plainText != word {
		t.Error("加密不对称")
	}
}

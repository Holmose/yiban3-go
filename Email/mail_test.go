package Email

import (
	"Yiban3/Browser/config"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"testing"
)

func TestMail(t *testing.T) {
	// 邮件接收方
	mailTo := []string{
		//可以是多个接收人
		"holmose@qq.com",
	}
	subject := "Hello World!" // 邮件主题
	body := "测试发送邮件"          // 邮件正文

	// 载入配置文件
	// 读取配置文件
	file, err := os.Open("config/config.json")
	defer file.Close()
	if err != nil {
		log.Panic(err)
	}
	var conf config.ConfigS
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&conf)
	if err != nil {
		log.Println("读取系统配置文件失败", err.Error())
	} else {
		log.Println("读取系统配置文件成功")
		config.MailUser = conf.MailUser
		config.MailPass = conf.MailPass
		config.MailHost = conf.MailHost
	}
	err = SendMail(mailTo, subject, body)
	if err != nil {
		fmt.Println("Send fail! - ", err)
		return
	}
	fmt.Println("Send successfully!")
}

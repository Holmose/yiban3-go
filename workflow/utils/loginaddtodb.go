package utils

import (
	"Yiban3/browser/tasks/login"
	browser "Yiban3/browser/types"
	"Yiban3/mysqlcon"
	"fmt"
	"log"
	"time"
)

// LoginAddVerifyToMysql 登录并添加数据到数据库
func LoginAddVerifyToMysql(b *browser.Browser) error {
	verifyISNil := false
	if b.User.Verify == "" {
		verifyISNil = true
	}
retry:
	_, err := login.Login(b)
	if err != nil {
		log.Println(err, "重试中。。。")
		time.Sleep(time.Second)
		goto retry
	}
	if verifyISNil {
		sql := fmt.Sprintf(
			"UPDATE yiban_yiban set verify=\"%s\" where username=\"%s\" ", b.User.Verify, b.User.Username)
		mysqlcon.Exec(sql)
	}
	return nil

}

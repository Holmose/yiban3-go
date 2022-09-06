package email

import (
	"Yiban3/browser/config"
	"Yiban3/browser/types"
	"crypto/tls"
	"fmt"
	"gopkg.in/gomail.v2"
	"log"
	"time"
)

/*
go邮件发送
*/

func SendMail(mailTo []string, subject string, body string) error {
	if config.MailUser == "" || config.MailPass == "" || config.MailHost == "" {
		return fmt.Errorf("未设置邮件发送服务配置，或参数配置错误")
	}
	// 设置邮箱主体
	mailConn := map[string]string{
		"user": config.MailUser, //发送人邮箱（邮箱以自己的为准）
		"pass": config.MailPass, //发送人邮箱的密码，现在可能会需要邮箱 开启授权密码后在pass填写授权码
		"host": config.MailHost, //邮箱服务器（此时用的是qq邮箱）
	}

	m := gomail.NewMessage(
		//发送文本时设置编码，防止乱码。 如果txt文本设置了之后还是乱码，那可以将原txt文本在保存时
		//就选择utf-8格式保存
		gomail.SetEncoding(gomail.Base64),
	)
	m.SetHeader("From", m.FormatAddress(mailConn["user"], "无敌打卡反馈")) // 添加别名
	m.SetHeader("To", mailTo...)                                     // 发送给用户(可以多个)
	m.SetHeader("Subject", subject)                                  // 设置邮件主题
	m.SetBody("text/html", body)                                     // 设置邮件正文

	//一个文件（加入发送一个 txt 文件）：/tmp/foo.txt，我需要将这个文件以邮件附件的方式进行发送，同时指定附件名为：附件.txt
	//同时解决了文件名乱码问题
	//name := "附件.txt"
	//m.Attach("E:/GoCode/src/goMail/gomail.txt",
	//	gomail.Rename(name), //重命名
	//	gomail.SetHeader(map[string][]string{
	//		"Content-Disposition": []string{
	//			fmt.Sprintf(`attachment; filename="%s"`, mime.QEncoding.Encode("UTF-8", name)),
	//		},
	//	}),
	//)

	/*
	   创建SMTP客户端，连接到远程的邮件服务器，需要指定服务器地址、端口号、用户名、密码，如果端口号为465的话，
	   自动开启SSL，这个时候需要指定TLSConfig
	*/
	d := gomail.NewDialer(mailConn["host"], 465, mailConn["user"], mailConn["pass"]) // 设置邮件正文
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	err := d.DialAndSend(m)
	return err
}

func YiSend(b *browser.Browser, detail browser.Detail,
	clockResult string, position browser.Position) {
	// 邮件接收方
	mailTo := []string{
		//可以是多个接收人
		b.User.Mail,
	}
	htmlBody := `<!DOCTYPE html>
        <html lang="en">
        <head>
            <meta charset="UTF-8">
            <title>%s</title>
            <style>
                body {
                    background-color: #FFFFCC;
                }
                h3 {
                    color: crimson;
                }
                h4 {
                    color: #ef3030;
                }
                a {
                    text-decoration: none;
                    color: #590b81;
                }
                span {
                    color: darkslateblue;
                    font-size: 15px;
                }
            </style>
        </head>
        <body>
            <h4>任务名称：<span>%s</span></h4>
            <h4>打卡时间：<span>%s</span></h4>
            <h4>账&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;户：<span>%s</span></h4>
            <h4>打卡结果：<span>%s</span></h4>
            <h4>签到位置: <span>&nbsp;&nbsp;%s</span></h4>
			<h4>剩余天数：<span>%v天</span></h4>
        </body>
        </html>`
	mailBody := fmt.Sprintf(htmlBody, "无敌打卡",
		detail.Data.Title,
		time.Now().Format("2006-01-02 15:04"),
		b.User.Username,
		clockResult,
		fmt.Sprintf("%v：(%v, %v)", position.Address,
			position.Longitude, position.Latitude),
		b.User.Day)
	subject := "无敌打卡" // 邮件主题

retry:
	err := SendMail(mailTo, subject, mailBody)
	if err != nil {
		// 重试
		log.Printf("用户：%v 邮件发送失败：%v，重试中...", b.User.Username, err)
		time.Sleep(time.Second * 60)
		goto retry
	} else {
		log.Printf("[用户：%v 邮件发送成功]", b.User.Username)
	}
}

// YiTips 发送提醒消息
func YiTips(mailTo []string, mailBody string) {
	htmlBody := `<!DOCTYPE html>
        <html lang="en">
        <head>
            <meta charset="UTF-8">
            <title>%s</title>
            <style>
                body {
                    background-color: #FFFFCC;
                }
                h3 {
                    color: crimson;
                }
                h4 {
                    color: #ef3030;
                }
                a {
                    text-decoration: none;
                    color: #590b81;
                }
                span {
                    color: darkslateblue;
                    font-size: 15px;
                }
            </style>
        </head>
        <body>
			%v
        </body>
        </html>`
	mailBody = fmt.Sprintf(htmlBody, "无敌打卡", mailBody)
	subject := "无敌打卡" // 邮件主题

retry:
	err := SendMail(mailTo, subject, mailBody)
	if err != nil {
		// 重试
		log.Printf("用户：%v 邮件发送失败：%v，重试中...", mailTo, err)
		time.Sleep(time.Second * 60)
		goto retry
	} else {
		log.Printf("用户：%v 邮件发送成功", mailTo)
	}
}

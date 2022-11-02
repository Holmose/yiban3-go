package config

import (
	"bytes"
	"encoding/json"
	"log"
	"os"
	"time"
)

type GlobalConfiger interface {
	InitConfig(filePath string) // 初始化全局变量
	SaveConfig(filePath string) // 保存全局配置到文件
}

// 全局变量
var (
	// 打卡配置
	CSRF                     = "302ec8bc6a82d6bbf0ce37b7392d429e"
	MaxNum                   = 10
	ShowSecond time.Duration = 5
	// CompleteTemplateDelta 获取最近几天的打卡模板
	CompleteTemplateDelta = 8
	MysqlConStr           = ""
	// 邮件服务配置
	MailUser = ""
	MailPass = ""
	MailHost = ""
	// 配置打卡模板匹配字符串
	SubString = map[string]string{
		"Holiday": "学生身体状况",
		"Daily":   "体温报备"}
	// 定时任务配置
	PerMinute []int
	PerHour   []int
)

type ConfigS struct {
	CSRF                  string
	MaxNum                int
	ShowSecond            time.Duration
	CompleteTemplateDelta int
	MysqlConStr           string
	MailUser              string
	MailPass              string
	MailHost              string
	SubString             map[string]string
	PerMinute             []int
	PerHour               []int
}

func (c *ConfigS) InitConfig(filePath string) error {
	file, err := os.Open(filePath)
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Println("读取配置文件失败: ", err)
		}
	}(file)
	if err != nil {
		return err
	}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&c)
	if err != nil {
		return err
	} else {
		CSRF = c.CSRF
		MaxNum = c.MaxNum
		ShowSecond = c.ShowSecond
		CompleteTemplateDelta = c.CompleteTemplateDelta
		MysqlConStr = c.MysqlConStr
		MailUser = c.MailUser
		MailPass = c.MailPass
		MailHost = c.MailHost
		if c.SubString != nil {
			SubString = c.SubString
		} else {
			c.SubString = SubString
		}
		PerMinute = c.PerMinute
		PerHour = c.PerHour

		return nil
	}
}
func (c *ConfigS) SaveConfig(filePath string) error {
	by, err := json.Marshal(c)
	if err != nil {
		return err
	}
	// json格式化
	var out bytes.Buffer
	err = json.Indent(&out, by, "", "\t")
	if err != nil {
		return err
	}
	err = os.WriteFile(filePath, out.Bytes(), 777)
	if err != nil {
		return err
	}
	return nil
}

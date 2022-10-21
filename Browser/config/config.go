package config

import (
	"Yiban3/Browser/fetcher"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
	"time"
)

var (
	// 打卡配置
	CSRF                     = "302ec8bc6a82d6bbf0ce37b7392d429e"
	MaxNum                   = 3
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

func WriteSysconf(conf ConfigS) {
	by, err := json.Marshal(conf)
	if err != nil {
		log.Println(err)
	}
	// json格式化
	var out bytes.Buffer
	err = json.Indent(&out, by, "", "\t")
	if err != nil {
		log.Println(err)
	}
	err = os.WriteFile("config/config.json", out.Bytes(), 777)
	if err != nil {
		log.Println(err)
	}
}

// FetchConfig 从蓝奏云上获取配置信息
func FetchConfig(url string) (string, error) {
	client := &http.Client{
		Timeout: 18 * time.Second,
	}
	fetch, err := fetcher.Fetch(client, url)
	if err != nil {
		return "", errors.New("获取相关信息失败")
	}
	fileTeta := `<span class="teta tetb">说</span><span id="filename">([^<]*)</span></div><div class="d2">`
	re := regexp.MustCompile(fileTeta)
	matches := re.FindAllStringSubmatch(string(fetch), -1)

	if matches == nil || len(matches[0]) < 1 {
		return "", errors.New("没有获取到配置文件")
	} else {
		return matches[0][1], nil
	}
}

// 执行命令打卡浏览器
var commands = map[string]string{
	"windows": "start",
}

func Open(uri string) error {
	run, ok := commands[runtime.GOOS]
	if !ok {
		return fmt.Errorf("don't know how to open things on %s platform", runtime.GOOS)
	}

	cmd := exec.Command("cmd", "/C", run, uri)
	return cmd.Start()
}

func GetVersionMsg(url string, currentVersion string) error {
	versionMsg, err := FetchConfig(url)
	if err != nil {
		return err
	}
	if strings.Contains(versionMsg, currentVersion) {
		log.Println("当前版本为最新版")
	} else {
		log.Println("有最新版了，请更新")
		versionMsg = strings.Replace(versionMsg, "\n", " ", -1)
		fileVersion := `版本号：\[(.*)\] 下载链接：\[(https://.*)\]`
		re := regexp.MustCompile(fileVersion)
		matches := re.FindAllStringSubmatch(versionMsg, -1)
		if matches == nil || len(matches[0]) < 1 {
			return errors.New("没有获取到版本信息")
		}
		log.Printf("版本号：%v 下载链接：%v", matches[0][1], matches[0][2])
		err := Open(matches[0][2])
		if err != nil {
			return err
		}
	}
	return nil
}

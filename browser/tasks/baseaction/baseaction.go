package baseaction

import (
	"Yiban3/browser/config"
	"Yiban3/browser/types"
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strings"
	"time"
)

// GetCompleteList 获取打卡完成的列表 delta: 前几天的数据
func GetCompleteList(b *browser.Browser, delta int) (t browser.Tasks, err error) {
	currentTime := time.Now()
	oldTime := currentTime.AddDate(0, 0, -delta)
	currentTimeString := currentTime.Format("2006-01-02") + " 23:59"
	oldTimeString := oldTime.Format("2006-01-02") + " 00:00"

	completeUrl := "https://api.uyiban.com/officeTask/client/index/completedList?StartTime=%s&EndTime=%s&CSRF=%s"
	completeUrl = fmt.Sprintf(completeUrl, oldTimeString, currentTimeString, config.CSRF)

	resp, err := b.ClientGet(completeUrl)
	if err != nil {
		return t, err
	}
	defer resp.Body.Close()
	bodyReader := bufio.NewReader(resp.Body)
	bytes, err := io.ReadAll(bodyReader)
	if err != nil {
		return t, err
	}

	err = json.Unmarshal(bytes, &t)
	if err != nil {
		return t, err
	}

	return t, nil
}

// GetDetail 获取表单信息
func GetDetail(b *browser.Browser, d browser.Data) (dt browser.Detail, err error) {
	detailUrl := "https://api.uyiban.com/officeTask/client/index/detail?TaskId=%s&CSRF=%s"
	detailUrl = fmt.Sprintf(detailUrl, d.TaskID, config.CSRF)

	resp, err := b.ClientGet(detailUrl)
	defer resp.Body.Close()
	if err != nil {
		return dt, err
	}

	bodyReader := bufio.NewReader(resp.Body)
	bytes, err := io.ReadAll(bodyReader)

	if err != nil {
		return dt, err
	}

	err = json.Unmarshal(bytes, &dt)
	if err != nil {
		return dt, err
	}
	return dt, nil
}

// GetMessage 获取完成打卡的详细信息
func GetMessage(b *browser.Browser, d browser.Data) (comDet browser.CompleteDetail, err error) {
	detail, err := GetDetail(b, d)
	if err != nil {
		return comDet, err
	}
	messageUrl := "https://api.uyiban.com/workFlow/c/work/show/view/%s?CSRF=%s"
	messageUrl = fmt.Sprintf(messageUrl, detail.Data.InitiateId, config.CSRF)

	resp, err := b.ClientGet(messageUrl)
	defer resp.Body.Close()
	if err != nil {
		return comDet, err
	}
	bodyReader := bufio.NewReader(resp.Body)
	bytes, err := io.ReadAll(bodyReader)

	err = json.Unmarshal(bytes, &comDet)
	if err != nil {
		return comDet, err
	}

	return comDet, nil
}

// GetUnCompleteList 获取未打卡的列表
func GetUnCompleteList(b *browser.Browser) (t browser.Tasks, err error) {
	currentTime := time.Now()
	currentTimeString := currentTime.Format("2006-01-02") + " 23:59"
	oldTimeString := currentTime.Format("2006-01-02") + " 00:00"

	unCompleteUrl := "https://api.uyiban.com/officeTask/client/index/uncompletedList?StartTime=%s&EndTime=%s&CSRF=%s"
	unCompleteUrl = fmt.Sprintf(unCompleteUrl, oldTimeString, currentTimeString, config.CSRF)

	resp, err := b.ClientGet(unCompleteUrl)
	if err != nil {
		return t, err
	}
	defer resp.Body.Close()
	bodyReader := bufio.NewReader(resp.Body)
	bytes, err := io.ReadAll(bodyReader)

	err = json.Unmarshal(bytes, &t)
	if err != nil {
		return t, err
	}

	return t, nil
}

// FetchUnComplete 从未打卡列表中拿出一条
func FetchUnComplete(b *browser.Browser, t browser.Tasks) (d browser.Data, err error) {
	for _, d := range t.Data {
		timeUnix := time.Now().Unix()

		if timeUnix >= d.StartTime && strings.Contains(d.Title, "体温报备") {
			log.Printf("(** 用户：%v %v 未打卡! **)", b.User.Username, d.Title)
			return d, nil
		} else if strings.Contains(d.Title, "学生身体状况采集") {
			log.Printf("(** 用户：%v %v 未打卡! **)", b.User.Username, d.Title)
			return d, nil
		}

	}

	return d, fmt.Errorf("没有未打卡数据！")
}

// CreateForm 创建打卡表单（获取在服务器已经生成的）
func CreateForm(b *browser.Browser, d browser.Data) (form browser.FormTask, err error) {
	// 获取表单信息
	detail, err := GetDetail(b, d)
	if err != nil {
		return form, err
	}

	formUrl := "https://api.uyiban.com/workFlow/c/my/form/%s?CSRF=%s"
	formUrl = fmt.Sprintf(formUrl, detail.Data.WFId, config.CSRF)

	resp, err := b.ClientGet(formUrl)
	defer resp.Body.Close()
	if err != nil {
		return form, err
	}
	bodyReader := bufio.NewReader(resp.Body)
	bytes, err := io.ReadAll(bodyReader)

	if err != nil {
		return form, err
	}

	json.Unmarshal(bytes, &form)

	return form, nil
}

// SubmitForm 提交表单
func SubmitForm(b *browser.Browser, formEncrypt string) (bytes []byte, err error) {
	submitUrl := "https://api.uyiban.com/workFlow/c/my/apply?CSRF=%s"
	submitUrl = fmt.Sprintf(submitUrl, config.CSRF)

	params := fmt.Sprintf("Str=%s", formEncrypt)
	// 发送请求
	resp, err := b.ClientPost(submitUrl, params)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}
	bodyReader := bufio.NewReader(resp.Body)
	bytes, err = io.ReadAll(bodyReader)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

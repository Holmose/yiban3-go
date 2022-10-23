package action

import (
	"Yiban3/Browser/fetcher"
	"Yiban3/Browser/tasks/baseaction"
	"Yiban3/Browser/tasks/clock"
	browser "Yiban3/Browser/types"
	"Yiban3/Email"
	"Yiban3/Workflow/utils"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"
)

/*
	详细打卡任务节点
*/

// 获取传入的浏览器对象
func getBrowser(i interface{}) browser.Browser {
	datas := i.(map[string]interface{})
	for {
		if datas["loginBrowser"] != nil {
			break
		}
		time.Sleep(time.Millisecond * 50)
	}
	b := datas["loginBrowser"].(browser.Browser)
	return b
}

// PositionTemplateAction 根据位置获取一条已经打卡成功的数据
type PositionTemplateAction struct{}

func (a *PositionTemplateAction) Run(i interface{}) {
	b := getBrowser(i)

	retryCount := 0
retry:
	completeDetail, err := clock.GetCompleteByLocation(&b, b.User.Position)
	if err != nil {
		log.Printf("[用户：%v %v 重试中...]", b.User.Username, err)
		retryCount++
		if retryCount <= 13 {
			time.Sleep(time.Second)
			goto retry
		} else {
			// 发送邮件传递执行结果
			datas := i.(map[string]interface{})
			bodyMail := fmt.Sprintf(
				"<h2>账号：%v 获取一个在%v 的打卡模板失败，请至少保证最近几天在该位置有打卡记录！</h2>",
				b.User.Username, b.User.Position)

			datas["SendTipsMsg"] = map[string]interface{}{
				"mailTo":   b.User.Mail,
				"mailBody": bodyMail,
			}
		}
	}
	b.ChanData.CompleteDetailChan <- completeDetail
}

// UnClockListAction 获取未打卡的列表
type UnClockListAction struct{}

func (a *UnClockListAction) Run(i interface{}) {
	b := getBrowser(i)

	unCompleteList, err := baseaction.GetUnCompleteList(&b)
retry:
	if err != nil {
		log.Println(err, "重试中。。。")
		time.Sleep(time.Second)
		goto retry
	}

	unComplete, err := baseaction.FetchUnComplete(&b, unCompleteList)
	if err != nil {
		if !strings.Contains(err.Error(), "未打卡数据不存在") {
			log.Println(err)
		}
		log.Printf("[[ 用户：%v %v ]]", b.User.Username, err)
		// TODO 会重复关闭
		close(b.ChanData.UnCompleteChan)
		return
	} else {
		// 用于获取打卡表单信息
		b.ChanData.UnCompleteChan <- unComplete
		// 用于获取更为详细的表单信息
		b.ChanData.UnCompleteChan <- unComplete
	}
}

// CreateFormAction 获取打卡表单信息
type CreateFormAction struct{}

func (a *CreateFormAction) Run(i interface{}) {
	b := getBrowser(i)

	unComplete, ok := <-b.ChanData.UnCompleteChan
	if !ok {
		close(b.ChanData.FormChan)
		return
	} else {
		form, err := baseaction.CreateForm(&b, unComplete)
		if err != nil {
			log.Println(err)
		}
		b.ChanData.FormChan <- form
	}
}

// GetDetailFormAction 获取更为详细的表单信息
type GetDetailFormAction struct {
	once sync.Once // 确保只会关闭一次
}

func (a *GetDetailFormAction) Run(i interface{}) {
	b := getBrowser(i)

	unComplete, ok := <-b.ChanData.UnCompleteChan
	if ok {
		detail, err := baseaction.GetDetail(&b, unComplete)
		if err != nil {
			log.Println(err)
		}
		b.ChanData.DetailChan <- detail
	} else {
		close(b.ChanData.DetailChan)
	}
}

// FillFormSubmitAction 填写打卡表单并提交
type FillFormSubmitAction struct{}

func (a *FillFormSubmitAction) Run(i interface{}) {
	b := getBrowser(i)

	// 获取chan数据
	form, ok1 := <-b.ChanData.FormChan
	detail, ok2 := <-b.ChanData.DetailChan
	completeDetail, ok3 := <-b.ChanData.CompleteDetailChan
	// 没有数据直接返回
	if !(ok1 && ok2 && ok3) {
		return
	}

	// 填写打卡表单
	var formEncrypt string
	var position browser.Position
	var err error
	if b.User.IsHoliday {
		formEncrypt, position, err = clock.FillHolidayForm(form, detail, completeDetail)
		if err != nil {
			log.Println(err)
		}
		log.Printf("用户：%v 表单填写完成！", b.User.Username)
	} else {
		formEncrypt, position, err = clock.FillForm(form, detail, completeDetail)
		if err != nil {
			log.Println(err)
		}
		log.Printf("用户：%v 表单填写完成！", b.User.Username)
	}

	// 信息不对，停止提交
	if completeDetail.Data.WFName == "" {
		return
	}
	// 提交表单
	bytes, err := baseaction.SubmitForm(&b, formEncrypt)
	unicode, err := fetcher.ZhToUnicode(bytes)
	if err != nil {
		log.Println(err)
	}
	// 检查打卡结果
	var ret map[string]interface{}
	err = json.Unmarshal(unicode, &ret)
	if err != nil {
		log.Println(err)
	}
	// 插入任务结果到数据库中
	go utils.InsertTask(&b, detail, position, ret)
	// 插入打卡模板到数据库中 TODO 便于后期无法获取模板时使用
	go utils.InsertForm(&b, position, completeDetail)

	// 传递执行结果
	datas := i.(map[string]interface{})
	datas["SendMailMsg"] = map[string]interface{}{
		"ret":      ret,
		"detail":   detail,
		"position": position,
	}
	log.Printf("[用户：%v, 打卡结束！------]", b.User.Username)
}

// SendMailAction 发送打卡结果邮件
type SendMailAction struct{}

func (a *SendMailAction) Run(i interface{}) {
	b := getBrowser(i)
	datas := i.(map[string]interface{})
	SendMailMsg := datas["SendMailMsg"]
	if SendMailMsg == nil {
		return
	}
	SendMail := SendMailMsg.(map[string]interface{})
	ret := SendMail["ret"].(map[string]interface{})
	detail := SendMail["detail"].(browser.Detail)
	position := SendMail["position"].(browser.Position)

	var clockResult string
	if ret["data"] == nil {
		clockResult = fmt.Sprintf("打卡失败 %v", ret["msg"])
	} else {
		clockResult = fmt.Sprintf("打卡成功！")
	}

	// 发送邮件
retry:
	err := Email.YiSend(&b, detail, clockResult, position)
	if err != nil {
		// 重试
		log.Printf("用户：%v 邮件发送失败：%v，重试中...", b.User.Username, err)
		time.Sleep(time.Second * 60)
		goto retry
	}
	log.Printf("[用户：%v 邮件发送成功]", b.User.Username)
}

// SendTipsMailAction 发送提示邮件
type SendTipsMailAction struct{}

func (a *SendTipsMailAction) Run(i interface{}) {
	// 判断是否发送
	datas := i.(map[string]interface{})
	SendTips := datas["SendTipsMsg"]
	if SendTips == nil {
		return
	}
	SendTipsMsg := SendTips.(map[string]interface{})
	mailTo := SendTipsMsg["mailTo"].(string)
	mailBody := SendTipsMsg["mailBody"].(string)

	// 确保邮件发送成功
retry:
	err := Email.YiTips([]string{mailTo},
		mailBody)
	if err != nil {
		// 重试
		log.Printf("用户：%v 邮件发送失败：%v，重试中...", mailTo, err)
		time.Sleep(time.Second * 60)
		goto retry
	}
	log.Printf("[用户：%v 邮件发送成功]", mailTo)
}

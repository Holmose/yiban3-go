package schedule

import (
	"Yiban3/browser/config"
	"Yiban3/browser/fetcher"
	"Yiban3/browser/tasks/baseaction"
	"Yiban3/browser/tasks/clock"
	"Yiban3/browser/tasks/login"
	"Yiban3/browser/types"
	"Yiban3/ecryption/yiban"
	"Yiban3/email"
	"Yiban3/mysqlcon"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
)

var wg sync.WaitGroup
var isFirst = true

// ChanListRun 运行入口
func ChanListRun() {
	// 载入配置文件
	AddUserToQ()
	if isFirst {
		if config.CSRF == "" {
			log.Println("系统配置CSRF必须配置")
			os.Exit(99)
		}
		wg.Add(1)
		go CreateBrowserByUserQ()
		wg.Wait()
	}
	ScheduleSim(browserQ)
	if isFirst {
		log.Println("用户配置格式化写入成功！")
		WriteUserinfo(userConfigList)
	}
	isFirst = false
}

// ChanListRunMysql 运行入口
func ChanListRunMysql() {
	// 从数据库中载入数据
	browserQ = []*browser.Browser{}
	AddUserToQByMysql()
	if isFirst {
		if config.CSRF == "" {
			log.Println("系统配置CSRF必须配置")
			os.Exit(99)
		}
		isFirst = false
	}
	wg.Add(1)
	go CreateBrowserByUserQ()
	wg.Wait()
	ScheduleSim(browserQ)
}

// 所有用户剩余天数减一
func DayReduce() {
	reduceSql := "UPDATE yiban_yiban set day=day-1 where day>0"
	mysqlcon.Exec(reduceSql)
}

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

var userConfigList []browser.User

func Run(b *browser.Browser) {
	defer wg.Done()
	<-pool
	// 登录并添加数据到数据库
	err := LoginAddVerifyToMysql(b)

	// 获取一条打卡模板
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-pool
		retryCount := 0
	retry:
		completeDetail, err := clock.GetCompleteByLocation(b, b.User.Position)
		if err != nil {
			log.Printf("用户：%v %v 重试中...", err, b.User.Username)
			retryCount++
			if retryCount <= 13 {
				goto retry
			} else {
				// 发送邮件
				bodyMail := fmt.Sprintf(
					"<h2>账号：%v 获取一个在%v 的打卡模板失败，请至少保证最近几天在该位置有打卡记录！</h2>",
					b.User.Username, b.User.Position)
				email.YiTips([]string{b.User.Mail},
					bodyMail)
			}
		}
		b.ChanData.CompleteDetailChan <- completeDetail
	}()

	// 获取未打卡记录
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-pool
		unCompleteList, err := baseaction.GetUnCompleteList(b)
	retry:
		if err != nil {
			log.Println(err, "重试中。。。")
			time.Sleep(time.Second)
			goto retry
		}

		unComplete, err := baseaction.FetchUnComplete(b, unCompleteList)
		if err != nil {
			if strings.Contains(err.Error(), "没有未打卡数据") {
				log.Printf("用户：%v %v", b.User.Username, err)
				wg.Add(-4)
				runtime.Goexit()
			} else {
				log.Println(err)
			}
		}

		b.ChanData.UnCompleteChan <- unComplete
		b.ChanData.UnCompleteChan <- unComplete
	}()

	// 创建打卡表单
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-pool
		unComplete := <-b.ChanData.UnCompleteChan
		form, err := baseaction.CreateForm(b, unComplete)
		if err != nil {
			log.Println(err)
		}
		b.ChanData.FormChan <- form

	}()

	// 获取表单信息
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-pool
		unComplete := <-b.ChanData.UnCompleteChan
		detail, err := baseaction.GetDetail(b, unComplete)
		if err != nil {
			log.Println(err)
		}
		b.ChanData.DetailChan <- detail
	}()

	// 获取chan数据
	form := <-b.ChanData.FormChan
	detail := <-b.ChanData.DetailChan
	completeDetail := <-b.ChanData.CompleteDetailChan

	// 填写打卡表单
	var formEncrypt string
	var position browser.Position
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

	if completeDetail.Data.WFName == "" {
		return
	}

	// 提交表单
	bytes, err := baseaction.SubmitForm(b, formEncrypt)
	unicode, err := fetcher.ZhToUnicode(bytes)
	if err != nil {
		log.Println(err)
	}

	var ret map[string]interface{}
	err = json.Unmarshal(unicode, &ret)
	if err != nil {
		log.Println(err)
	}

	clockResult := ""
	if ret["data"] == nil {
		log.Printf("用户：%v, 打卡失败：%v", b.User.Username,
			ret["msg"])
		clockResult = fmt.Sprintf("打卡失败 %v", ret["msg"])
	} else {
		log.Printf("用户：%v, 打卡成功！", b.User.Username)
		clockResult = fmt.Sprintf("打卡成功！")
	}
	wg.Add(2)
	// 插入任务结果到数据库中
	go InsertTask(b, detail, position, ret)
	// 发送邮件进行提示
	// 插入打卡模板到数据库中 TODO 便于后期无法获取模板时使用
	go InsertForm(b, position, completeDetail)
	email.YiSend(b, detail, clockResult, position)
	log.Printf("用户：%v, 打卡结束！------", b.User.Username)
}

// InsertTask 插入打卡结果
func InsertTask(b *browser.Browser,
	detail browser.Detail,
	position browser.Position,
	ret map[string]interface{}) {
	defer wg.Done()
	// 查询用户在数据库中的ID
	idSelect := fmt.Sprintf("select id from yiban_yiban where username=%s", b.User.Username)
	rst, err := mysqlcon.Query(idSelect)
	if err != nil {
		log.Println("没有找到数据!")
		return
	} else {
		log.Println("获取用户数据成功！")
	}
	clockStatus := 0
	var clockResult string
	if ret["data"] == nil {
		clockResult = ret["msg"].(string)
	} else {
		clockStatus = 1
		clockResult = "打卡成功！"
	}
	taskInsert := fmt.Sprintf("INSERT INTO yiban_task("+
		"yiban_id, name, position, result,"+
		"position_num,message,datetime)values("+
		"%s,\"%s\",\"%s\",%d,"+
		"\"%s\",\"%s\",\"%s\");",
		rst[0]["id"], detail.Data.Title, position.Address, clockStatus,
		fmt.Sprintf("(%v, %v)", position.Longitude, position.Latitude),
		clockResult, time.Now().Format("2006-01-02 15:04"))
	mysqlcon.Exec(taskInsert)
}

func InsertForm(b *browser.Browser,
	position browser.Position,
	completeDetail browser.CompleteDetail) {
	defer wg.Done()
	// 查询用户在数据库中的ID
	idSelect := fmt.Sprintf("select id from yiban_yiban where username=%s", b.User.Username)
	rst, err := mysqlcon.Query(idSelect)
	if err != nil {
		log.Println("没有找到数据!")
		return
	} else {
		log.Println("获取用户数据成功！")
	}
	// 是假期为0，不是假期为1
	isHolidayNum := 0
	if strings.Contains(completeDetail.Data.WFName, "学生身体状况") {
		isHolidayNum = 1
	}
	marshal, err := json.Marshal(completeDetail)
	if err != nil {
		log.Println(err)
		return
	}
	formEncrypt, err := yiban.FormEncrypt(string(marshal))
	if err != nil {
		log.Println(err)
		return
	}
	formInsert := fmt.Sprintf(
		"insert into yiban_formdata("+
			"yiban_id,name,address,data,datetime,status"+
			")values(%s, \"%s\", \"%s\",\"%s\",\"%s\", %v);",
		rst[0]["id"], completeDetail.Data.WFName, position.Address, formEncrypt,
		time.Now().Format("2006-01-02 15:04"), isHolidayNum)

	mysqlcon.Exec(formInsert)
}
func WriteUserinfo(userConfigList []browser.User) {
	by, err := json.Marshal(userConfigList)
	if err != nil {
		log.Println(err)
	}
	// json格式化
	var out bytes.Buffer
	err = json.Indent(&out, by, "", "\t")
	if err != nil {
		log.Println(err)
	}
	err = os.WriteFile("config/userinfo.json", out.Bytes(), 777)
	if err != nil {
		log.Println(err)
	}
}

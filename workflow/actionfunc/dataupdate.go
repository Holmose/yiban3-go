package actionfunc

import (
	browser "Yiban3/browser/types"
	"Yiban3/ecryption/yiban"
	"Yiban3/mysqlcon"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"
)

// InsertTask 插入打卡结果
func InsertTask(b *browser.Browser,
	detail browser.Detail,
	position browser.Position,
	ret map[string]interface{}) {
	retryCount := 0
retry:
	// 查询用户在数据库中的ID
	idSelect := fmt.Sprintf("select id from yiban_yiban where username=%s", b.User.Username)
	rst, ok := mysqlcon.Query(idSelect)
	if !ok {
		log.Println("[Database connection failed.] [Retry...]")
		retryCount++
		if retryCount <= 10 {
			time.Sleep(time.Second)
			goto retry
		}
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
	retryCount := 0
retry:
	// 查询用户在数据库中的ID
	idSelect := fmt.Sprintf("select id from yiban_yiban where username=%s", b.User.Username)
	rst, ok := mysqlcon.Query(idSelect)
	if !ok {
		log.Println("[Database connection failed.] [Retry...]")
		retryCount++
		if retryCount <= 10 {
			time.Sleep(time.Second)
			goto retry
		}
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

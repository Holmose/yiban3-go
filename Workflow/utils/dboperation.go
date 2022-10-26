package utils

import (
	"Yiban3/Browser/tasks/login"
	browser "Yiban3/Browser/types"
	"Yiban3/Ecryption/yiban"
	"Yiban3/MysqlConnect"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

// LoginAddVerifyToMysql 登录并添加数据到数据库
func LoginAddVerifyToMysql(b *browser.Browser) error {
	hverify := b.User.Verify
retry:
	_, err := login.Login(b)
	if err != nil {
		log.Println(err, "重试中。。。")
		time.Sleep(time.Second)
		goto retry
	}
	if hverify == "" || b.User.Verify != hverify {
		sql := fmt.Sprintf(
			"UPDATE yiban_yiban set verify=\"%s\" where username=\"%s\" ", b.User.Verify, b.User.Username)
		MysqlConnect.Exec(sql)
	}
	return nil

}

// QueryYibanUserToQ 分页查询用户数据
func QueryYibanUserToQ(userChan *YibanChan, userCount *[]int) {
	retryCount := 0
retry:
	// status 0 为否 false，不是假期
	userTotal := "SELECT COUNT(*) FROM yiban_yiban where day>0;"
	rst, err := MysqlConnect.Query(userTotal)
	if err != nil {
		log.Println("[Database connection failed.] [Retry...]")
		retryCount++
		if retryCount <= 10 {
			time.Sleep(time.Second)
			goto retry
		}
	}
	count, err := strconv.Atoi(rst[0]["COUNT(*)"])
	if err != nil {
		log.Println(err)
	}
	*userCount = []int{count}

	// 添加分页查询
	pageNum := 5
	// 获取分多少页
	pageCount := fmt.Sprintf(
		"select ceil(count(*)/%v) as pageTotal from yiban_yiban where day>0;", pageNum)
	rst, err = MysqlConnect.Query(pageCount)
	if err != nil {
		log.Println("[Failed to obtain paging data. ] [Trying again...]")
		retryCount++
		if retryCount <= 10 {
			time.Sleep(time.Second)
			goto retry
		}
	}
	pageTotal, err := strconv.Atoi(rst[0]["pageTotal"])
	for i := 0; i < pageTotal; i++ {
		// 获取每页数据
		pageMsg := fmt.Sprintf(
			"select * from yiban_yiban where day>0 limit %v offset %v;",
			pageNum, i*pageNum)
		rst, err = MysqlConnect.Query(pageMsg)
		if err != nil {
			log.Println("[Failed to obtain paging data. ] [Trying again...]")
			retryCount++
			if retryCount <= 10 {
				time.Sleep(time.Second)
				goto retry
			}
		}
		q, err := GetUserToQ(rst)
		for _, user := range q {
			userChan.C <- user
		}
		if err != nil {
			log.Println("[Failed to obtain the remaining time!]")
		}
	}
	// 安全关闭通道
	userChan.SafeClose()
	log.Println("[后端数据获取完成]")
}

func GetUserToQ(rst []map[string]string) ([]browser.User, error) {
	var q []browser.User
	for _, yiban := range rst {
		status := yiban["status"]
		isHoliday := false
		if status != "0" {
			isHoliday = true
		}
		day, err := strconv.Atoi(yiban["day"])
		if err != nil {
			return nil, err
		}
		userInfo := browser.User{
			Username:   yiban["username"],
			Password:   yiban["password"],
			Verify:     yiban["verify"],
			Position:   yiban["address"],
			Mail:       yiban["e_mail"],
			Crontab:    yiban["clock_crontab"],
			IsHoliday:  isHoliday,
			CreateTime: yiban["create_time"],
			UpdateTime: yiban["update_time"],
			Day:        day,
		}
		q = append(q, userInfo)
	}
	return q, nil
}

// InsertTask 插入打卡结果
func InsertTask(b *browser.Browser,
	detail browser.Detail,
	position browser.Position,
	ret map[string]interface{}) {
	retryCount := 0
retry:
	// 查询用户在数据库中的ID
	idSelect := fmt.Sprintf("select id from yiban_yiban where username=%s", b.User.Username)
	rst, err := MysqlConnect.Query(idSelect)
	if err != nil {
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
	MysqlConnect.Exec(taskInsert)
}

func InsertForm(b *browser.Browser,
	position browser.Position,
	completeDetail browser.CompleteDetail) {
	retryCount := 0
retry:
	// 查询用户在数据库中的ID
	idSelect := fmt.Sprintf("select id from yiban_yiban where username=%s", b.User.Username)
	rst, err := MysqlConnect.Query(idSelect)
	if err != nil {
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

	MysqlConnect.Exec(formInsert)
}

// DayReduce 所有用户剩余天数减一
func DayReduce() {
	reduceSql := "UPDATE yiban_yiban set day=day-1 where day>0"
	MysqlConnect.Exec(reduceSql)
}

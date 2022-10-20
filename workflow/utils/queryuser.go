package utils

import (
	browser "Yiban3/browser/types"
	"Yiban3/mysqlcon"
	"Yiban3/workflow/mychan"
	"fmt"
	"log"
	"strconv"
	"time"
)

// QueryYibanUserToQ 分页查询用户数据
func QueryYibanUserToQ(userChan *mychan.YibanChan, userCount *[]int) {
	retryCount := 0
retry:
	// status 0 为否 false，不是假期
	userTotal := "SELECT COUNT(*) FROM yiban_yiban where day>0;"
	rst, ok := mysqlcon.Query(userTotal)
	if !ok {
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
	rst, ok = mysqlcon.Query(pageCount)
	if !ok {
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
		rst, ok = mysqlcon.Query(pageMsg)
		if !ok {
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
			Username:  yiban["username"],
			Password:  yiban["password"],
			Verify:    yiban["verify"],
			Position:  yiban["address"],
			Mail:      yiban["e_mail"],
			Crontab:   yiban["clock_crontab"],
			IsHoliday: isHoliday,
			Day:       day,
		}
		q = append(q, userInfo)
	}
	return q, nil
}

package schedule

import (
	"Yiban3/browser/config"
	"Yiban3/browser/types"
	"Yiban3/mysqlcon"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
)

func AddUserToQ() {
	loadSysconf()
	//读取JSON
	file, err := os.Open("config/userinfo.json")
	defer file.Close()
	if err != nil {
		log.Panic(err)
	}
	var userInfo []browser.User
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&userInfo)
	if err != nil {
		log.Println("读取用户配置信息失败", err.Error())
	} else {
		log.Println("读取用户配置信息成功")
		userQ = userInfo
		userCount = len(userInfo)
	}
}

// CheckUser 当用户数量发生变化时 创建一个协程
func CheckUser() {
	retryCount := 0
retry:
	// status 0 为否 false，不是假期
	userTotal := "SELECT COUNT(*) FROM yiban_yiban where day>0;"
	rst, ok := mysqlcon.Query(userTotal)
	if ok {
		fmt.Println("用户数量心跳检测！", rst[0]["COUNT(*)"])
	} else {
		log.Println("没有找到数据!")
		retryCount++
		if retryCount <= 10 {
			goto retry
		}
	}
	count, err := strconv.Atoi(rst[0]["COUNT(*)"])
	if err != nil {
		log.Println(err)
	}
	if userCount != count {
		// 执行主程序
		log.Println("用户数量变化，执行打开程序")
		ChanListRunMysql()
	}
}

func AddUserToQByMysql() {
	loadSysconf()
	userQ = []browser.User{}
	retryCount := 0
retry:
	// status 0 为否 false，不是假期
	userTotal := "SELECT COUNT(*) FROM yiban_yiban where day>0;"
	rst, ok := mysqlcon.Query(userTotal)
	if ok {
		log.Println("获取数据成功！")
	} else {
		log.Println("没有找到数据!")
		retryCount++
		if retryCount <= 10 {
			goto retry
		}
	}
	count, err := strconv.Atoi(rst[0]["COUNT(*)"])
	if err != nil {
		log.Println("获取用户总数失败")
	}
	userCount = count

	// 添加分页查询
	pageNum := 5
	// 获取分多少页
	pageCount := fmt.Sprintf(
		"select ceil(count(*)/%v) as pageTotal from yiban_yiban where day>0;", pageNum)
	rst, ok = mysqlcon.Query(pageCount)
	if !ok {
		log.Println("没有找到数据!")
		retryCount++
		if retryCount <= 10 {
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
		if ok {
			log.Println("获取分页数据成功！", i*pageNum)
		} else {
			log.Println("没有找到数据!")
			retryCount++
			if retryCount <= 10 {
				goto retry
			}
		}
		q, err := GetUserToQ(rst)
		userQ = append(userQ, q...)
		if err != nil {
			log.Printf("获取剩余天数失败: %v", err)
		}
	}
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
func loadSysconf() {
	// 读取配置文件
	file, err := os.Open("config/config.json")
	defer file.Close()
	if err != nil {
		log.Panic(err)
	}
	var conf config.ConfigS
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&conf)
	if err != nil {
		log.Println("读取系统配置文件失败", err.Error())
	} else {
		log.Println("读取系统配置文件成功")
		config.CSRF = conf.CSRF
		config.MaxNum = conf.MaxNum
		config.ShowSecond = conf.ShowSecond
		config.CompleteTemplateDelta = conf.CompleteTemplateDelta
		config.MysqlConStr = conf.MysqlConStr
		config.MailUser = conf.MailUser
		config.MailPass = conf.MailPass
		config.MailHost = conf.MailHost
		if conf.SubString != nil {
			config.SubString = conf.SubString
		} else {
			conf.SubString = config.SubString
		}
		config.PerMinute = conf.PerMinute
		config.PerHour = conf.PerHour
		config.WriteSysconf(conf)
	}
}

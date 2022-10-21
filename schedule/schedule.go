package schedule

import (
	"Yiban3/Browser/config"
	"Yiban3/Browser/types"
	"Yiban3/MysqlConnect"
	"fmt"
	"github.com/robfig/cron/v3"
	"log"
	"sync"
	"time"
)

var userCount = 0

// 控制协程数量
var pool chan struct{}

// 全局变量存储定时任务用户
var Global_User_Cron map[string]browser.User

// ScheduleSim 统一分配个人浏览器
func ScheduleSim(browserQ []*browser.Browser) {
	if Global_User_Cron == nil {
		Global_User_Cron = make(map[string]browser.User)
	}
	// 最大协程数
	pool = make(chan struct{}, config.MaxNum)
	go func() {
		for {
			pool <- struct{}{}
		}
	}()

	log.Printf("当前用户总数：%v", userCount)
	for {
		if len(browserQ) < 1 {
			break
		}
		b := browserQ[0]
		browserQ = browserQ[1:]
		if b.User.Crontab != "" ||
			b.User.Crontab != Global_User_Cron[b.User.Username].Crontab {
			// 创建用户的定时打卡
			go CronUserCreate(b)
		} else {
			// 使用默认规则进行打卡
			wg.Add(1)
			go Run(b)
		}
	}
	wg.Wait()
}
func CronUserCreate(b *browser.Browser) {
	var wgcron sync.WaitGroup
	var mutex sync.Mutex
	spec := b.User.Crontab
	if _, ok := Global_User_Cron[b.User.Username]; spec != "" && !ok {
		c := cron.New(cron.WithChain())
		_, err := c.AddFunc(spec, func() {
			log.Printf("%v 用户：%v个人定时打卡任务执行",
				time.Now().Format("2006年01月02日15:04"), b.User.Username)
			user := getUser(b)
			// 定时任务发生变化，重建
			if user.Crontab != Global_User_Cron[b.User.Username].Crontab {
				b.User.Cron = nil
				log.Printf("用户：%v的个人定时打卡任务已停止", b.User.Username)
				c.Stop()
			} else {
				b := browser.Browser{}
				browser.CreateBrowser(&b, user)
				wg.Add(1)
				Run(&b)
			}
		})
		b.User.Cron = c
		mutex.Lock()
		Global_User_Cron[b.User.Username] = b.User
		mutex.Unlock()

		if err != nil {
			log.Printf("用户：%v 个人定时任务创建失败", b.User.Username)
		} else {
			log.Printf("用户：%v 个人定时任务创建成功，等待执行中...", b.User.Username)
			wgcron.Add(1)
			defer wgcron.Done()
			c.Start()
			defer c.Stop()
		}
		wgcron.Wait()
	} else {
		log.Printf("用户：%v 个人定时任务已经创建过了，或使用的是默认规则", b.User.Username)
	}
}
func getUser(b *browser.Browser) browser.User {
	retryCount := 0
retry:
	// 获取每页数据
	userMsg := fmt.Sprintf(
		"select * from yiban_yiban where day>0 and username=%v limit 1;",
		b.User.Username)
	rst, err := MysqlConnect.Query(userMsg)
	if err != nil {
		log.Println("没有找到数据!")
		retryCount++
		if retryCount <= 10 {
			goto retry
		}
	} else {
		log.Println("获取用户数据成功！")
	}
	q, err := GetUserToQ(rst)
	if err != nil {
		log.Printf("获取剩余天数失败: %v", err)
	}
	return q[0]
}

package workflow

import (
	"Yiban3/browser/config"
	"Yiban3/workflow/action/utils"
	"fmt"
	"github.com/robfig/cron/v3"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"
)

/*
	定时任务系统
*/

type CronTaskSystemAction struct {
	once sync.Once // 限制只能被执行一次
}

func (a *CronTaskSystemAction) Run(i interface{}) {
	// 执行定时任务 只能执行一次
	a.once.Do(func() {
		fmt.Println("[定时任务系统启动...]")
		cronRun()
	})
}

func cronTaskCreate() (*cron.Cron, error) {
	c := cron.New(cron.WithChain())

	// 转换为字符串数组
	var perMinuteStrArr []string
	var perHourStrArr []string

	if config.PerMinute == nil || config.PerHour == nil {
		return c, fmt.Errorf("未设置定时任务配置")
	}
	for _, minute := range config.PerMinute {
		minuteStr := strconv.Itoa(minute)
		perMinuteStrArr = append(perMinuteStrArr, minuteStr)
	}
	perMinuteStr := strings.Join(perMinuteStrArr, ",")
	for _, minute := range config.PerHour {
		minuteStr := strconv.Itoa(minute)
		perHourStrArr = append(perHourStrArr, minuteStr)
	}
	perHourStr := strings.Join(perHourStrArr, ",")

	spec := fmt.Sprintf("%v %v * * *", perMinuteStr, perHourStr)
	_, err := c.AddFunc(spec, func() {
		log.Printf("[%v 定时打卡任务执行]\n", time.Now().Format("2006年01月02日15:04"))
		log.Println("执行打卡逻辑。。。。。")
	})
	if err != nil {
		return c, err
	}
	// 每天2点 剩余天数减一
	spec = fmt.Sprintf("0 2 * * *")
	_, err = c.AddFunc(spec, func() {
		log.Println(time.Now().Format("2006年01月02日15:04"), "定时剩余天数减一任务执行")
		actionfunc.actionfunc.DayReduce()
	})
	if err != nil {
		return c, err
	}
	// 用户数变化检查
	spec = fmt.Sprintf("*/10 9-17 * * *")
	_, err = c.AddFunc(spec, func() {
		log.Println(time.Now().Format("2006年01月02日15:04"), "用户数量心跳检测执行")
		log.Println("心跳检测执行")
	})
	if err != nil {
		return c, err
	}
	return c, nil
}

func cronRun() {
	var wg sync.WaitGroup
	// 创建定时任务
	cronTask, err := cronTaskCreate()
	if err != nil {
		log.Printf("定时任务创建失败 %v", err)
	} else {
		log.Println("定时任务创建成功，等待执行中...")
		wg.Add(1)
		defer wg.Done()
		cronTask.Start()
		defer cronTask.Stop()
	}
	wg.Wait()
}

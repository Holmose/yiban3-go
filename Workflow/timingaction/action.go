package timingaction

import (
	"log"
	"sync"
)

// CronTaskByConfigAction 根据配置文件创建定时任务
type CronTaskByConfigAction struct {
	once sync.Once // 限制只能被执行一次
}

func (a *CronTaskByConfigAction) Run(i interface{}) {
	// 执行打卡任务，如果系统中存在用户配置cron则跳过打卡
	a.once.Do(func() {
		log.Println("[打卡任务创建中]")
		var wg sync.WaitGroup
		// 创建定时任务
		cronTask, err := clockTaskCron()
		if err != nil {
			log.Printf("[打卡任务创建失败 %v]", err)
		} else {
			log.Println("[打卡任务创建成功，等待执行中...]")
			wg.Add(1)
			defer wg.Done()
			cronTask.Start()
			defer cronTask.Stop()
		}
		wg.Wait()
	})
}

// DateUpdateAction 日期更新
type DateUpdateAction struct {
	once sync.Once // 限制只能被执行一次
}

func (a *DateUpdateAction) Run(i interface{}) {
	a.once.Do(func() {
		log.Println("[日期更新任务创建中]")
		var wg sync.WaitGroup
		// 创建定时任务
		cronTask, err := dailyReduceCron()
		if err != nil {
			log.Printf("[日期更新任务创建失败 %v]", err)
		} else {
			log.Println("[日期更新任务创建成功，等待执行中...]")
			wg.Add(1)
			defer wg.Done()
			cronTask.Start()
			defer cronTask.Stop()
		}
		wg.Wait()
	})
}

// CronTaskBySingleAction 根据用户cron创建打卡任务
type CronTaskBySingleAction struct{}

func (a *CronTaskBySingleAction) Run(i interface{}) {
	log.Println("根据用户cron创建打卡任务")
	clockFilterExec()
}

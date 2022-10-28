package timingaction

import (
	"Yiban3/Workflow/utils"
	"github.com/robfig/cron/v3"
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
type CronTaskBySingleAction struct {
	once sync.Once // 限制只能被执行一次
}

func (a *CronTaskBySingleAction) Run(i interface{}) {
	log.Println("根据用户cron创建打卡任务")
	c := cron.New(cron.WithChain())
	cronUsers := map[string]utils.CronUser{}
	var wgm sync.WaitGroup

	// 定时监测数据变化 不wait不会停止，因为有其他系统在运行
	wgm.Add(1)
	go a.once.Do(func() {
		log.Println("[心跳检测创建中]")
		var wg sync.WaitGroup
		// 创建定时任务
		monitor, err := monitorData(c, cronUsers)
		if err != nil {
			log.Printf("[心跳检测创建失败 %v]", err)
		} else {
			log.Println("[心跳检测创建成功，等待执行中...]")
			wg.Add(1)
			defer wg.Done()
			monitor.Start()
			defer monitor.Stop()
		}
		wg.Wait()
	})

	wgm.Add(1)
	defer wgm.Done()
	log.Println("个人任务定时管理器启动")
	c.Start()
	defer c.Stop()

	wgm.Wait()
}

package timingaction

import (
	"fmt"
	"log"
	"sync"
)

// CronTaskByConfigAction 根据配置文件创建定时任务
type CronTaskByConfigAction struct {
	once sync.Once // 限制只能被执行一次
}

func (a *CronTaskByConfigAction) Run(i interface{}) {
	// 执行定时任务 只能执行一次
	a.once.Do(func() {
		log.Println("[基础定时任务启动]")
		var wg sync.WaitGroup
		// 创建定时任务
		cronTask, err := clockTaskCron()
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
	})
}

// EndAction 功能拓展占位
type EndAction struct{}

func (a *EndAction) Run(i interface{}) {
	fmt.Println("[定时任务功能拓展占位]")
}

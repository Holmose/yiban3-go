package workflow

import (
	browser "Yiban3/browser/types"
	"Yiban3/workflow/actionfunc"
	"Yiban3/workflow/mychan"
	"log"
	"sync"
	"time"
)

// LoadSystemConfigAction 加载系统配置
type LoadSystemConfigAction struct{}

func (a *LoadSystemConfigAction) Run(i interface{}) {
	err := actionfunc.LoadSystemConfig("config/config.json")
	if err != nil {
		log.Panic("[加载系统配置] [失败]", err.Error())
	} else {
		log.Println("[加载系统配置] [成功]")
	}
}

// NewUserChanAction 从数据库获取用户信息数组
type NewUserChanAction struct{}

func (a *NewUserChanAction) Run(i interface{}) {
	log.Println("[获取用户数据]")
	var userCount []int
	var userChan *mychan.YibanChan

	userChan = mychan.NewYibanChan()
	go actionfunc.QueryYibanUserToQ(userChan, &userCount)

	// 传递数据
	datas := i.(map[string]interface{})
	datas["userCount"] = &userCount
	datas["userChan"] = userChan
}

// NewBrowserChanAction 为每个用户创建浏览器对象
type NewBrowserChanAction struct{}

func (a *NewBrowserChanAction) Run(i interface{}) {
	log.Println("[获取用户创建浏览器对象]")
	datas := i.(map[string]interface{})
	userCount := datas["userCount"].(*[]int)
	for {
		if len(*userCount) > 0 {
			break
		}
		time.Sleep(time.Millisecond * 50)
	}
	userChan := datas["userChan"].(*mychan.YibanChan)

	var wg sync.WaitGroup

	var browserChan *mychan.YibanChan
	browserChan = mychan.NewYibanChan()

	wg.Add(1)
	go func() {
		count := 0
		for {
			if count >= (*userCount)[0] {
				wg.Done()
				break
			}
			select {
			// 1.队列中有数据时取出数组
			case user, ok := <-userChan.C:
				if ok {
					// 2.创建Browser (基于用户信息创建)
					b := browser.Browser{}
					browser.CreateBrowser(&b, user.(browser.User))

					// 3.加入Browser通道
					wg.Add(1)
					go func() {
						browserChan.C <- b
						count++
						wg.Done()
					}()
				}
			}
		}
	}()

	// 传递数据
	datas["browserChan"] = browserChan

	// 等待完成
	wg.Wait()

	// 安全关闭通道
	browserChan.SafeClose()
}

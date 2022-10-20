package action

import (
	"Yiban3/browser/config"
	browser "Yiban3/browser/types"
	"Yiban3/workflow/mychan"
	"Yiban3/workflow/utils"
	"fmt"
	"log"
	"sync"
	"time"
)

/*
	打卡程序流程块
*/

// LoginAction 取出浏览器对象，并执行登录操作
type LoginAction struct{}

func (a *LoginAction) Run(i interface{}) {
	datas := i.(map[string]interface{})
	for {
		if datas["browserChan"] != nil {
			break
		}
		time.Sleep(time.Millisecond * 50)
	}

	browserChan := datas["browserChan"].(*mychan.YibanChan)
	userCount := *(datas["userCount"].(*[]int))

	var wg sync.WaitGroup
	var loginChan *mychan.YibanChan
	loginChan = mychan.NewYibanChan()

	wg.Add(1)
	go func() {
		count := 0
		for {
			if count >= userCount[0] {
				wg.Done()
				break
			}
			select {
			case brow, ok := <-browserChan.C:
				if ok {
					wg.Add(1)
					go func() {
						// 登录并添加数据到数据库
						b := brow.(browser.Browser)
						err := utils.LoginAddVerifyToMysql(&b)
						if err != nil {
							fmt.Println(err)
						} else {
							loginChan.C <- b
							count++
						}
						wg.Done()
					}()
				}
			}
		}
	}()
	log.Println("[创建登录进程] [成功]")
	// 传递数据
	datas["loginChan"] = loginChan

	// 等待完成
	wg.Wait()

	// 安全关闭通道
	loginChan.SafeClose()
}

// GetLoginBrowserAction 获取浏览器对象执行打卡任务
type GetLoginBrowserAction struct {
	ClockWorkflow func(interface{}) // 执行创建的打卡流
}

func (a *GetLoginBrowserAction) Run(i interface{}) {
	datas := i.(map[string]interface{})
	for {
		if datas["loginChan"] != nil {
			break
		}
		time.Sleep(time.Millisecond * 50)
	}
	loginChan := datas["loginChan"].(*mychan.YibanChan)
	userCount := *(datas["userCount"].(*[]int))

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		count := 0
		// 控制协程数量
		var pool chan struct{}
		// 最大协程数
		pool = make(chan struct{}, config.MaxNum)
		go func() {
			for {
				pool <- struct{}{}
			}
		}()
		for {
			if count >= userCount[0] {
				wg.Done()
				break
			}
			select {
			case loginBrowser := <-loginChan.C:
				<-pool
				go func() {
					wg.Add(1)
					a.ClockWorkflow(loginBrowser)
					wg.Done()
				}()
				count++
			}
		}
	}()
	log.Println("[核心程序加载] [完成]")
	wg.Wait()
	log.Println("[本次打卡结束!]")
}

// EndAction 功能拓展占位
type EndAction struct{}

func (a *EndAction) Run(i interface{}) {
	log.Println("[功能拓展占位]")
}

package action

import (
	"Yiban3/browser/config"
	browser "Yiban3/browser/types"
	"Yiban3/workflow/action/utils"
	"Yiban3/workflow/mychan"
	"context"
	"fmt"
	"github.com/Holmose/go-workflow/workflow"
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
						err := actionfunc.LoginAddVerifyToMysql(&b)
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
type GetLoginBrowserAction struct{}

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
					clockWf(loginBrowser)
					wg.Done()
				}()
				count++
			}
		}
	}()

	log.Println("[核心程序加载] [完成]")
	wg.Wait()
	fmt.Println("[本次打卡结束!]")
}

// 创建工作流进行打卡
func clockWf(loginBrowser interface{}) {
	// 获取一个浏览器对象，发送数据到数据流中
	wf := workflow.NewWorkFlow()
	// 构建节点
	PositionTemplateNode := workflow.NewNode(&PositionTemplateAction{}) // 获取位置模板
	UnClockListNode := workflow.NewNode(&UnClockListAction{})           // 获取未打卡的列表
	CreateFormNode := workflow.NewNode(&CreateFormAction{})             // 获取打卡表单信息
	GetDetailFormNode := workflow.NewNode(&GetDetailFormAction{})       // 获取更为详细的表单信息
	FillFormSubmitNode := workflow.NewNode(&FillFormSubmitAction{})     // 填写打卡表单并提交

	// 构建节点之间的关系
	// 启始节点
	wf.AddStartNode(PositionTemplateNode)
	wf.AddStartNode(UnClockListNode)

	// 中间节点
	wf.AddEdge(UnClockListNode, CreateFormNode)
	wf.AddEdge(CreateFormNode, FillFormSubmitNode)
	wf.AddEdge(UnClockListNode, GetDetailFormNode)
	wf.AddEdge(GetDetailFormNode, FillFormSubmitNode)
	wf.AddEdge(PositionTemplateNode, FillFormSubmitNode)

	// 收尾节点
	wf.ConnectToEnd(FillFormSubmitNode)

	// 数据
	var completedAction map[string]interface{}
	completedAction = make(map[string]interface{})
	completedAction["loginBrowser"] = loginBrowser

	ctx, _ := context.WithCancel(context.Background())
	wf.StartWithContext(ctx, completedAction)
	wf.WaitDone()
}

// EndAction 功能拓展占位
type EndAction struct{}

func (a *EndAction) Run(i interface{}) {
	fmt.Println("[功能拓展占位]")
}

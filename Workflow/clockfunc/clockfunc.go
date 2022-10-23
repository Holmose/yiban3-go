package clockfunc

import (
	browser "Yiban3/Browser/types"
	"Yiban3/Workflow/graphnode"
	"context"
	"github.com/Holmose/go-workflow/workflow"
	"github.com/robfig/cron/v3"
	"log"
	"sync"
	"time"
)

/*
  根据需要选择定时执行或全部执行
*/

// ClockWorkflow 无过滤全部执行打卡
func ClockWorkflow(loginBrowser interface{}) {
	// 获取一个浏览器对象，发送数据到数据流中
	wf := workflow.NewWorkFlow()
	// 构建节点
	PositionTemplateNode := workflow.NewNode(&action.PositionTemplateAction{}) // 获取位置模板
	SendTipsMailNode := workflow.NewNode(&action.SendTipsMailAction{})         // 模板获取失败发送提示邮件
	UnClockListNode := workflow.NewNode(&action.UnClockListAction{})           // 获取未打卡的列表
	CreateFormNode := workflow.NewNode(&action.CreateFormAction{})             // 获取打卡表单信息
	GetDetailFormNode := workflow.NewNode(&action.GetDetailFormAction{})       // 获取更为详细的表单信息
	FillFormSubmitNode := workflow.NewNode(&action.FillFormSubmitAction{})     // 填写打卡表单并提交
	SendMailNode := workflow.NewNode(&action.SendMailAction{})                 // 打卡成功发送邮件

	// 构建节点之间的关系
	// 启始节点
	wf.AddStartNode(PositionTemplateNode) // 获取位置模板
	wf.AddStartNode(UnClockListNode)      // 获取未打卡的列表

	// 中间节点
	wf.AddEdge(UnClockListNode, CreateFormNode)    // 获取未打卡的列表->填写打卡表单并提交
	wf.AddEdge(CreateFormNode, FillFormSubmitNode) // 获取打卡表单信息->填写打卡表单并提交

	wf.AddEdge(PositionTemplateNode, FillFormSubmitNode) // 获取位置模板->填写打卡表单并提交
	wf.AddEdge(PositionTemplateNode, SendTipsMailNode)   // 获取位置模板->模板获取失败发送提示邮件

	wf.AddEdge(UnClockListNode, GetDetailFormNode)    // 获取未打卡的列表->获取更为详细的表单信息
	wf.AddEdge(GetDetailFormNode, FillFormSubmitNode) // 获取更为详细的表单信息->填写打卡表单并提交

	wf.AddEdge(FillFormSubmitNode, SendMailNode) // 填写打卡表单并提交->打卡成功发送邮件

	// 收尾节点
	wf.ConnectToEnd(SendTipsMailNode) // 模板获取失败发送提示邮件
	wf.ConnectToEnd(SendMailNode)     // 打卡成功发送邮件

	// 数据
	var completedAction map[string]interface{}
	completedAction = make(map[string]interface{})
	completedAction["loginBrowser"] = loginBrowser

	ctx, _ := context.WithCancel(context.Background())
	wf.StartWithContext(ctx, completedAction)
	wf.WaitDone()
}

// ClockWorkflowFilter 有过滤 只执行无个人定时
func ClockWorkflowFilter(loginBrowser interface{}) {
	// 判断是否丢弃
	b := loginBrowser.(browser.Browser)
	if b.User.Crontab != "" {
		return
	} else {
		ClockWorkflow(loginBrowser)
	}
}

// ClockWorkflowCronSingle 根据个人定时创建定时任务
func ClockWorkflowCronSingle(loginBrowser interface{}) {
	// 判断是否存在cron配置
	b := loginBrowser.(browser.Browser)
	if b.User.Crontab == "" {
		return
	}
	// 创建定时任务
	var wg sync.WaitGroup
	spec := b.User.Crontab
	c := cron.New(cron.WithChain())
	_, err := c.AddFunc(spec, func() {
		log.Printf("[%v 用户：%v个人定时打卡任务执行]",
			time.Now().Format("2006年01月02日15:04"), b.User.Username)
		ClockWorkflow(loginBrowser)
	})
	if err != nil {
		log.Printf("[用户：%v 个人定时任务创建失败]", b.User.Username)
	} else {
		log.Printf("[用户：%v 个人定时任务创建成功，等待执行中...]", b.User.Username)
		wg.Add(1)
		defer wg.Done()
		c.Start()
		defer c.Stop()
	}
	wg.Wait()
}

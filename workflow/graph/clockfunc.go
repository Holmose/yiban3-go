package graph

import (
	browser "Yiban3/browser/types"
	"Yiban3/workflow/clockaction"
	"context"
	"github.com/Holmose/go-workflow/workflow"
)

// ClockWorkflow 无过滤全部执行打卡
func ClockWorkflow(loginBrowser interface{}) {
	// 获取一个浏览器对象，发送数据到数据流中
	wf := workflow.NewWorkFlow()
	// 构建节点
	PositionTemplateNode := workflow.NewNode(&action.PositionTemplateAction{}) // 获取位置模板
	UnClockListNode := workflow.NewNode(&action.UnClockListAction{})           // 获取未打卡的列表
	CreateFormNode := workflow.NewNode(&action.CreateFormAction{})             // 获取打卡表单信息
	GetDetailFormNode := workflow.NewNode(&action.GetDetailFormAction{})       // 获取更为详细的表单信息
	FillFormSubmitNode := workflow.NewNode(&action.FillFormSubmitAction{})     // 填写打卡表单并提交

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

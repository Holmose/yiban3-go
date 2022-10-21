package Flowcharts

import (
	"Yiban3/Workflow/clockfunc"
	"Yiban3/Workflow/graphnode"
	"Yiban3/Workflow/graphnode/initialize"
	"context"
	"github.com/Holmose/go-workflow/workflow"
)

/*
创建工作流程图
*/

// ClockOnce 执行打卡程序
func ClockOnce() {
	wf := workflow.NewWorkFlow()

	// 构建节点
	LoadSystemConfigNode := workflow.NewNode(&initialize.LoadSystemConfigAction{}) // 加载系统配置
	NewUserChanNode := workflow.NewNode(&initialize.NewUserChanAction{})           // 从数据库获取用户信息，存入通道
	CreateBrowserNode := workflow.NewNode(&initialize.NewBrowserChanAction{})      // 为每个用户创建浏览器对象，存入通道
	LoginNode := workflow.NewNode(&action.LoginAction{})                           // 获取浏览器对象进行登录
	GetLoginBrowserNode := workflow.NewNode(
		&action.GetLoginBrowserAction{ClockWorkflow: clockfunc.ClockWorkflow}) // 获取浏览器对象执行打卡任务
	EndNode := workflow.NewNode(&action.EndAction{}) // 结束占位

	// 构建节点之间的关系
	// 启始节点
	wf.AddStartNode(LoadSystemConfigNode) // 加载系统配置
	wf.AddStartNode(LoginNode)            // 获取浏览器对象进行登录
	wf.AddStartNode(GetLoginBrowserNode)  // 获取浏览器对象执行打卡任务

	// 中间节点
	wf.AddEdge(LoadSystemConfigNode, NewUserChanNode) // 加载系统配置->从数据库获取用户信息，存入通道
	wf.AddEdge(NewUserChanNode, CreateBrowserNode)    // 从数据库获取用户信息，存入通道->为每个用户创建浏览器对象，存入通道
	wf.AddEdge(LoginNode, EndNode)                    // 获取浏览器对象进行登录->结束占位
	wf.AddEdge(GetLoginBrowserNode, EndNode)          // 获取浏览器对象执行打卡任务->结束占位

	// 收尾节点
	wf.ConnectToEnd(CreateBrowserNode) // 为每个用户创建浏览器对象，存入通道
	wf.ConnectToEnd(EndNode)           // 结束占位

	// 数据
	var completedAction map[string]interface{}
	completedAction = make(map[string]interface{})

	// 执行
	ctx, _ := context.WithCancel(context.Background())
	wf.StartWithContext(ctx, completedAction)
	wf.WaitDone()
}

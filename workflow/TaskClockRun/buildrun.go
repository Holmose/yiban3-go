package TaskClockRun

import (
	"Yiban3/workflow/action"
	"context"
	"github.com/Holmose/go-workflow/workflow"
)

/*
创建工作流程图
*/

// TaskClockRun 执行打卡程序
func TaskClockRun() {
	wf := workflow.NewWorkFlow()

	// 构建节点
	LoadSystemConfigNode := workflow.NewNode(&action.LoadSystemConfigAction{}) // 加载系统配置
	GetUserArrayNode := workflow.NewNode(&action.NewUserChanAction{})          // 从数据库获取用户信息数组
	CreateBrowserNode := workflow.NewNode(&action.NewBrowserChanAction{})      // 为每个用户创建浏览器对象
	LoginNode := workflow.NewNode(&action.LoginAction{})                       // 获取浏览器对象进行登录
	GetLoginBrowserNode := workflow.NewNode(&action.GetLoginBrowserAction{})   // 获取浏览器对象执行打卡任务

	EndNode := workflow.NewNode(&action.EndAction{}) // 结束占位

	// 构建节点之间的关系
	// 启始节点
	wf.AddStartNode(LoadSystemConfigNode)
	wf.AddStartNode(LoginNode)
	wf.AddStartNode(GetLoginBrowserNode)

	// 中间节点
	wf.AddEdge(LoadSystemConfigNode, GetUserArrayNode)
	wf.AddEdge(GetUserArrayNode, CreateBrowserNode)
	wf.AddEdge(LoginNode, EndNode)
	wf.AddEdge(GetLoginBrowserNode, EndNode)

	// 收尾节点
	wf.ConnectToEnd(CreateBrowserNode)
	wf.ConnectToEnd(EndNode)

	// 数据
	var completedAction map[string]interface{}
	completedAction = make(map[string]interface{})

	// 执行
	ctx, _ := context.WithCancel(context.Background())
	wf.StartWithContext(ctx, completedAction)
	wf.WaitDone()
}

package Flowcharts

import (
	"Yiban3/Workflow/graphnode/initialize"
	"Yiban3/Workflow/timingaction"
	"context"
	"github.com/Holmose/go-workflow/workflow"
)

/*
创建工作流程图
*/

// ClockTimingSys 定时任务系统
func ClockTimingSys() {
	wf := workflow.NewWorkFlow()

	// 构建节点
	LoadSystemConfigNode := workflow.NewNode(&initialize.LoadSystemConfigAction{})   // 加载系统配置
	CronTaskByConfigNode := workflow.NewNode(&timingaction.CronTaskByConfigAction{}) // 根据配置文件创建定时任务
	DateUpdateNode := workflow.NewNode(&timingaction.DateUpdateAction{})             // 日期更新
	CronTaskBySingleNode := workflow.NewNode(&timingaction.CronTaskBySingleAction{}) // 根据用户cron创建打卡任务

	// 构建节点之间的关系
	// 启始节点
	wf.AddStartNode(LoadSystemConfigNode) // 加载系统配置

	// 中间节点
	wf.AddEdge(LoadSystemConfigNode, CronTaskByConfigNode) // 加载系统配置->根据配置文件创建定时任务
	wf.AddEdge(LoadSystemConfigNode, DateUpdateNode)       // 加载系统配置->日期更新
	wf.AddEdge(LoadSystemConfigNode, CronTaskBySingleNode) // 加载系统配置->根据用户cron创建打卡任务

	// 收尾节点
	wf.ConnectToEnd(DateUpdateNode)       // 日期更新
	wf.ConnectToEnd(CronTaskBySingleNode) // 根据用户cron创建打卡任务

	// 数据
	var completedAction map[string]interface{}
	completedAction = make(map[string]interface{})

	// 执行
	ctx, _ := context.WithCancel(context.Background())
	wf.StartWithContext(ctx, completedAction)
	wf.WaitDone()
}

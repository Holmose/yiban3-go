package timingaction

import (
	"Yiban3/Browser/config"
	"Yiban3/Workflow/clockfunc"
	action "Yiban3/Workflow/graphnode"
	"Yiban3/Workflow/graphnode/initialize"
	"Yiban3/Workflow/utils"
	"context"
	"fmt"
	"github.com/Holmose/go-workflow/workflow"
	"github.com/robfig/cron/v3"
	"log"
	"strconv"
	"strings"
	"time"
)

// 定时打卡任务
func clockTaskCron() (*cron.Cron, error) {
	c := cron.New(cron.WithChain())

	// 解析配置
	var perMinuteStrArr []string
	var perHourStrArr []string

	if config.PerMinute == nil || config.PerHour == nil {
		return c, fmt.Errorf("未设置定时任务配置")
	}
	for _, minute := range config.PerMinute {
		minuteStr := strconv.Itoa(minute)
		perMinuteStrArr = append(perMinuteStrArr, minuteStr)
	}
	perMinuteStr := strings.Join(perMinuteStrArr, ",")
	for _, minute := range config.PerHour {
		minuteStr := strconv.Itoa(minute)
		perHourStrArr = append(perHourStrArr, minuteStr)
	}
	perHourStr := strings.Join(perHourStrArr, ",")

	// 创建任务
	spec := fmt.Sprintf("%v %v * * *", perMinuteStr, perHourStr)
	_, err := c.AddFunc(spec, func() {
		log.Printf("[%v 定时打卡任务执行]\n", time.Now().Format("2006年01月02日15:04"))
		clockExec()
	})
	if err != nil {
		return c, err
	}

	return c, nil
}

// 每天2点 剩余天数减一
func dailyReduceCron() (*cron.Cron, error) {
	c := cron.New(cron.WithChain())

	spec := fmt.Sprintf("0 2 * * *")
	_, err := c.AddFunc(spec, func() {
		log.Println(time.Now().Format("2006年01月02日15:04"), "定时剩余天数减一任务执行")
		utils.DayReduce()
	})

	if err != nil {
		return c, err
	}
	return c, nil
}

// 用户数变化检查
func userCheckCron() (*cron.Cron, error) {
	c := cron.New(cron.WithChain())

	spec := fmt.Sprintf("*/10 9-17 * * *")
	_, err := c.AddFunc(spec, func() {
		log.Println(time.Now().Format("2006年01月02日15:04"), "用户数量心跳检测执行")
		log.Println("心跳检测执行")
	})
	if err != nil {
		return c, err
	}
	return c, nil
}

// clockExec 跳过存在cron的用户进行打卡
func clockExec() {
	wf := workflow.NewWorkFlow()

	// 构建节点
	LoadSystemConfigNode := workflow.NewNode(&initialize.LoadSystemConfigAction{}) // 加载系统配置
	GetUserArrayNode := workflow.NewNode(&initialize.NewUserChanAction{})          // 从数据库获取用户信息数组
	CreateBrowserNode := workflow.NewNode(&initialize.NewBrowserChanAction{})      // 为每个用户创建浏览器对象
	LoginNode := workflow.NewNode(&action.LoginAction{})                           // 获取浏览器对象进行登录
	GetLoginBrowserNode := workflow.NewNode(
		&action.GetLoginBrowserAction{
			ClockWorkflow: clockfunc.ClockWorkflowFilter}) // 获取浏览器对象执行打卡任务，跳过存在cron的用户进行打卡

	// 构建节点之间的关系
	// 启始节点
	wf.AddStartNode(LoadSystemConfigNode)
	wf.AddStartNode(LoginNode)
	wf.AddStartNode(GetLoginBrowserNode)

	// 中间节点
	wf.AddEdge(LoadSystemConfigNode, GetUserArrayNode)
	wf.AddEdge(GetUserArrayNode, CreateBrowserNode)

	// 收尾节点
	wf.ConnectToEnd(CreateBrowserNode)

	// 数据
	var completedAction map[string]interface{}
	completedAction = make(map[string]interface{})

	// 执行
	ctx, _ := context.WithCancel(context.Background())
	wf.StartWithContext(ctx, completedAction)
	wf.WaitDone()
}

// clockFilterExec 根据用户cron创建定时任务
func clockFilterExec() {
	wf := workflow.NewWorkFlow()

	// 构建节点
	LoadSystemConfigNode := workflow.NewNode(&initialize.LoadSystemConfigAction{}) // 加载系统配置
	GetUserArrayNode := workflow.NewNode(&initialize.NewUserChanAction{})          // 从数据库获取用户信息数组
	CreateBrowserNode := workflow.NewNode(&initialize.NewBrowserChanAction{})      // 为每个用户创建浏览器对象
	LoginNode := workflow.NewNode(&action.LoginAction{})                           // 获取浏览器对象进行登录
	GetLoginBrowserNode := workflow.NewNode(
		&action.GetLoginBrowserAction{
			ClockWorkflow: clockfunc.ClockWorkflowCronSingle}) // 获取浏览器对象执行打卡任务，根据用户cron创建定时任务

	// 构建节点之间的关系
	// 启始节点
	wf.AddStartNode(LoadSystemConfigNode)
	wf.AddStartNode(LoginNode)
	wf.AddStartNode(GetLoginBrowserNode)

	// 中间节点
	wf.AddEdge(LoadSystemConfigNode, GetUserArrayNode)
	wf.AddEdge(GetUserArrayNode, CreateBrowserNode)

	// 收尾节点
	wf.ConnectToEnd(CreateBrowserNode)

	// 数据
	var completedAction map[string]interface{}
	completedAction = make(map[string]interface{})

	// 执行
	ctx, _ := context.WithCancel(context.Background())
	wf.StartWithContext(ctx, completedAction)
	wf.WaitDone()
}

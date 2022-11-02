package timingaction

import (
	"Yiban3/Browser/config"
	"Yiban3/MysqlConnect"
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
func clockFilterExec(personCrons utils.PersonalCrons) {
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
	completedAction["personCrons"] = personCrons

	// 执行
	ctx, _ := context.WithCancel(context.Background())
	wf.StartWithContext(ctx, completedAction)
	wf.WaitDone()
}

// 定时监测数据变化
func monitorData(personCrons utils.PersonalCrons) (*cron.Cron, error) {
	//c := cron.New(cron.WithChain())
	//spec := fmt.Sprintf("*/5 * * * *")

	c := cron.New(cron.WithSeconds()) // 支持秒级
	spec := fmt.Sprintf("*/10 * * * * *")

	_, err := c.AddFunc(spec, func() {
		//log.Println(time.Now().Format("2006年01月02日15:04"), "心跳检测")
		check(personCrons)
	})

	if err != nil {
		return c, err
	}
	return c, nil
}

func check(personCrons utils.PersonalCrons) {
	// 判断是否有条目被删除

	// 查询数据库
	sql := "SELECT username,clock_crontab,update_time,day FROM yiban_yiban;"
	rst, err := MysqlConnect.Query(sql)
	if err != nil {
		log.Println("[Database connection failed.]")
	}

	for _, v := range rst {
		value, ok := personCrons.Get(v["username"])

		// 新添加了一个cron
		if !ok && v["clock_crontab"] != "" && v["day"] != "0" {
			log.Printf("[定时任务系统] 用户 %v 检测到定时任务信息，添加中", v["username"])
			rebuild(personCrons, v["username"])
		}

		// 判断更新
		if ok && v["update_time"] != value.UpdateTime {
			log.Printf("[定时任务系统] 用户 %v 的信息更新，移除原始定时任务", v["username"])
			rebuild(personCrons, v["username"])
		}
		// 判断删除
		if ok && v["clock_crontab"] == "" {
			log.Printf("[定时任务系统] 用户 %v 删除定时任务，移除任务", v["username"])
			rebuild(personCrons, v["username"])
		}
	}

	// 判断用户信息删除
	personSlice := []string{}
	dbSlict := []string{}
	for key, _ := range personCrons.GetAll() {
		personSlice = append(personSlice, key)
	}
	for _, v := range rst {
		dbSlict = append(dbSlict, v["username"])
	}
	di := utils.Difference(personSlice, dbSlict)
	if len(di) != 0 {
		for _, s := range di {
			log.Printf("[定时任务系统] 用户 %v 数据被删除了, 定时任务移除", s)
			rebuild(personCrons, s)
		}
	}
}

func rebuild(personCrons utils.PersonalCrons,
	username string) {
	personCrons.Stop()
	// 移除定时任务
	personCrons.Remove(username)
	// 全部重建
	clockFilterExec(personCrons)
	personCrons.Start()
}

package main

import (
	program "Yiban3/workflow"
	"context"
	"fmt"
	"github.com/Holmose/go-workflow/workflow"
	"io"
	"log"
	"os"
)

func init() {
	// 获取日志文件句柄
	// 以 只写入文件|没有时创建|文件尾部追加 的形式打开这个文件
	logFile, err := os.OpenFile(`./日志文件.log`, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Panic(err)
	}
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	// 组合一下即可，os.Stdout代表标准输出流
	multiWriter := io.MultiWriter(os.Stderr, logFile)
	// 设置存储位置
	log.SetOutput(multiWriter)
}

func main() {
	wf := workflow.NewWorkFlow()

	// 构建节点
	LoadSystemConfigNode := workflow.NewNode(&program.LoadSystemConfigAction{}) // 加载系统配置
	GetUserArrayNode := workflow.NewNode(&program.NewUserChanAction{})          // 从数据库获取用户信息数组
	CreateBrowserNode := workflow.NewNode(&program.NewBrowserChanAction{})      // 为每个用户创建浏览器对象
	LoginNode := workflow.NewNode(&program.LoginAction{})                       // 获取浏览器对象进行登录
	GetLoginBrowserNode := workflow.NewNode(&program.GetLoginBrowserAction{})   // 获取浏览器对象执行打卡任务

	EndNode := workflow.NewNode(&program.EndAction{}) // 结束占位

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

	ctx, _ := context.WithCancel(context.Background())
	wf.StartWithContext(ctx, completedAction)
	wf.WaitDone()

	fmt.Println("执行其他逻辑")
}

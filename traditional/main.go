package main

import (
	"Yiban3/browser/config"
	"Yiban3/crontask"
	"Yiban3/schedule"
	"io"
	"log"
	"os"
	"strconv"
	"time"
)

func init() {
	// 获取日志文件句柄
	// 以 只写入文件|没有时创建|文件尾部追加 的形式打开这个文件
	logFile, err := os.OpenFile(`./日志文件.log`, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		// 日志文件打开失败，直接退出
		log.Panic(err)
	}
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	// 组合一下即可，os.Stdout代表标准输出流
	multiWriter := io.MultiWriter(os.Stderr, logFile)
	// 设置存储位置
	log.SetOutput(multiWriter)
}

func main() {
	timeUnix := time.Now().Unix()
	currentVersion := "v2.8 beta"

	log.Println("无敌打卡，在线版本：", currentVersion)

	// 获取到期时间
	expirationTime, err := config.FetchConfig("https://wwt.lanzouw.com/b066jv5zi")
	if err != nil {
		log.Println("获取信息失败")
	}
	parseInt, err := strconv.ParseInt(expirationTime, 10, 64)
	if err != nil {
		log.Panicf("配置信息有误")
	}
	tm := time.Unix(parseInt, 0)
	log.Printf("程序有效期到 %v", tm.Format("2006年01月02日15:04"))

	if timeUnix <= tm.Unix() {
		log.Println("当前有效...")
		go schedule.ChanListRunMysql()
		time.Sleep(time.Second * config.ShowSecond)
	} else {
		log.Println("程序已过期，请获取最新版。")
	}
	// 执行定时任务
	crontask.CronRun()
}

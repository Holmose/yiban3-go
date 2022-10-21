package main

import (
	"Yiban3/execute"
	"io"
	"log"
	"os"
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
	execute.ClockExec()
}

package crontask

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestCrons(t *testing.T) {
	c := cron.New(cron.WithChain())
	perMinute := []int{15, 30}
	perHour := []int{7, 9, 10, 12, 13, 14, 18, 19, 20}

	// 转换为字符串数组
	var perMinuteStrArr []string
	for _, minute := range perMinute {
		minuteStr := strconv.Itoa(minute)
		perMinuteStrArr = append(perMinuteStrArr, minuteStr)
	}
	perMinuteStr := strings.Join(perMinuteStrArr, ",")

	var perHourStrArr []string
	for _, minute := range perHour {
		minuteStr := strconv.Itoa(minute)
		perHourStrArr = append(perHourStrArr, minuteStr)
	}
	perHourStr := strings.Join(perHourStrArr, ",")

	spec := fmt.Sprintf("%v %v * * *", perMinuteStr, perHourStr)

	cronId1, err := c.AddFunc(spec, func() {
		fmt.Println(time.Now(), "定时任务执行。。。")
	})

	fmt.Println("执行定时任务id", cronId1, err)
	c.Start()
	defer c.Stop()
	time.Sleep(3 * time.Minute)
}

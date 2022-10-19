package clock

import (
	browser "Yiban3/browser/types"
	"fmt"
	"time"
)

// TestAction 测试占位
type TestAction struct{}

func (a *TestAction) Run(i interface{}) {
	datas := i.(map[string]interface{})
	for {
		if datas["loginBrowser"] != nil {
			break
		}
		time.Sleep(time.Millisecond * 50)
	}
	loginBrowser := datas["loginBrowser"].(browser.Browser)
	fmt.Printf("[用户：%v 打卡任务执行]\n",
		loginBrowser.User.Username)
}

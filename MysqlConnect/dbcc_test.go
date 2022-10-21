package MysqlConnect

import (
	"fmt"
	"testing"
)

func TestDb(t *testing.T) {
	rst, err := Query("select * from yiban_yiban")
	if err != nil {
		fmt.Println("没有找到数据")
	} else {
		fmt.Println(rst[0])
	}
}

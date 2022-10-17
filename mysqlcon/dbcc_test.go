package mysqlcon

import (
	"fmt"
	"testing"
)

func TestDb(t *testing.T) {
	rst, ok := Query("select * from yiban_yiban")
	if ok {
		fmt.Println(rst[0])
	} else {
		fmt.Println("没有找到数据")
	}
}

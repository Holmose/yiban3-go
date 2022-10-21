package MysqlConnect

import (
	"Yiban3/Browser/config"
	"database/sql"
	"errors"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"time"
)

var DB *sql.DB

func initDB() bool { //连接到MySQL
	//path := "root:go_db123@tcp(10.20.133.60:3306)/django_mysql?charset=utf8"
	//root = 用户名
	//password = 密码
	//mydb = 数据库名称
	retryCount := 0
retry:
	DB, _ = sql.Open("mysql", config.MysqlConStr)
	DB.SetConnMaxLifetime(100)
	DB.SetMaxIdleConns(10)
	//验证连接
	if err := DB.Ping(); err != nil {
		log.Println("[Database connection failed.] [Retry...]")
		time.Sleep(time.Second * 3)
		retryCount++
		if retryCount <= 10 {
			goto retry
		}
		return false

	}
	return true
}

func Exec(SQL string) {
	if initDB() == true {
		ret, err := DB.Exec(SQL) //增、删、改就靠这一条命令就够了，很简单
		if err != nil {
			log.Println(err)
		}
		_, err = ret.LastInsertId()
		if err != nil {
			log.Println(err)
		}
	}
}

func Query(SQL string) ([]map[string]string, error) { //通用查询寒素
	if initDB() != true { //连接数据库
		return nil, errors.New("连接数据库失败")
	}
	rows, err := DB.Query(SQL) //执行SQL语句，比如select * from users
	if err != nil {
		return nil, err
	}
	columns, _ := rows.Columns()            //获取列的信息
	count := len(columns)                   //列的数量
	var values = make([]interface{}, count) //创建一个与列的数量相当的空接口
	for i, _ := range values {
		var ii interface{} //为空接口分配内存
		values[i] = &ii    //取得这些内存的指针，因后继的Scan函数只接受指针
	}
	ret := make([]map[string]string, 0) //创建返回值：不定长的map类型切片
	for rows.Next() {
		err := rows.Scan(values...)  //开始读行，Scan函数只接受指针变量
		m := make(map[string]string) //用于存放1列的 [键/值] 对
		if err != nil {
			return nil, err
		}
		for i, colName := range columns {
			var raw_value = *(values[i].(*interface{})) //读出raw数据，类型为byte
			b, _ := raw_value.([]byte)
			v := string(b) //将raw数据转换成字符串
			m[colName] = v //colName是键，v是值
		}
		ret = append(ret, m) //将单行所有列的键值对附加在总的返回值上（以行为单位）
	}

	defer rows.Close()

	if len(ret) != 0 {
		return ret, nil
	}
	return nil, errors.New("数据解析失败")
}

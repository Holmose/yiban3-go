package models

import (
	"Yiban3/FrontWeb/helper"
	"errors"
	"time"
)

type Task struct {
	ID       int       `gorm:"cloumn:id;primaryKey" json:"id"`
	Name     string    `gorm:"column:name" json:"name"`                 // 任务名称
	Datetime time.Time `gorm:"column:datetime" json:"datetime"`         // 打卡时间
	Address  string    `gorm:"column:position" json:"position"`         // 位置
	Position string    `gorm:"column:position_num" json:"position_num"` // 坐标
	Result   int       `gorm:"column:result" json:"result"`             // 结果（成功或失败）
	Message  string    `gorm:"column:message" json:"message"`           // 返回信息
	// 外键关联到用户表
	User   User `json:"yiban"`
	UserId int  `gorm:"column:yiban_id" json:"yiban_id"`
}

type Tasks []Task

// TableName 自定义表名
func (t *Task) TableName() string {
	return "yiban_task"
}

// FindTaskByUserName 查询单个结果
func (t *Task) FindTaskByUserName(username string) error {
	user, err := CheckUser(username)
	if err != nil {
		helper.LogError(err.Error())
		return errors.New("用户查询失败")
	}
	if err := db.First(t, "yiban_id=?", user.ID).Error; err != nil {
		helper.LogError(err.Error())
		return errors.New("任务结果查询失败")
	}
	return nil
}

// FindTasksByUserName 查询多个结果
func (t *Tasks) FindTasksByUserName(username string, opts ...Option) error {
	// 初始化默认值
	query := PageQuery{
		Page: 1,
		Size: 10,
	}
	// 依次调用opts函数列表中的函数，为结构体成员赋值
	for _, o := range opts {
		o(&query)
	}
	user, err := CheckUser(username)
	if err != nil {
		helper.LogError(err.Error())
		return errors.New("用户查询失败")
	}
	if err := db.Scopes(Paginate(query.Page, query.Size)).
		Find(t, "yiban_id=?", user.ID).Error; err != nil {
		helper.LogError(err.Error())
		return errors.New("任务结果查询失败")
	}
	return nil
}

// FindAllByPage 根据分页查询用户信息
// page 页数 pageSize 每页显示的数量
func (t *Tasks) FindAllByPage(opts ...Option) error {
	// 初始化默认值
	query := PageQuery{
		Page: 1,
		Size: 10,
	}
	// 依次调用opts函数列表中的函数，为结构体成员赋值
	for _, o := range opts {
		o(&query)
	}
	if err := db.Scopes(Paginate(query.Page, query.Size)).Find(t).Error; err != nil {
		helper.LogError(err.Error())
		return errors.New("任务结果查询失败")
	}
	return nil
}

package models

import (
	"Yiban3/FrontWeb/helper"
	"errors"
	"time"
)

type Form struct {
	ID       int       `gorm:"cloumn:id;primaryKey" json:"id"`
	Name     string    `gorm:"column:name" json:"name"`         // 任务名称
	Datetime time.Time `gorm:"column:datetime" json:"datetime"` // 打卡时间
	Address  string    `gorm:"column:address" json:"address"`   // 位置
	Data     string    `gorm:"column:data" json:"data"`         // 加密表单数据
	Template int       `gorm:"column:result" json:"result"`     // 打卡模板
	// 外键关联到用户表
	User   User `json:"yiban"`
	UserId int  `gorm:"column:yiban_id" json:"yiban_id"`
}

type Forms []Form

// TableName 自定义表名
func (f *Form) TableName() string {
	return "yiban_formdata"
}

// FindFormByUserName 查询单个结果
func (f *Form) FindFormByUserName(username string) error {
	user, err := CheckUser(username)
	if err != nil {
		helper.LogError(err.Error())
		return errors.New("用户查询失败")
	}
	if err := db.First(f, "yiban_id=?", user.ID).Error; err != nil {
		helper.LogError(err.Error())
		return errors.New("任务结果查询失败")
	}
	return nil
}

// FindFormsByUserName 查询多个结果
func (f *Forms) FindFormsByUserName(username string, opts ...Option) error {
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
		Find(f, "yiban_id=?", user.ID).Error; err != nil {
		helper.LogError(err.Error())
		return errors.New("任务结果查询失败")
	}
	return nil
}

// FindAllByPage 根据分页查询用户信息
// page 页数 pageSize 每页显示的数量
func (f *Forms) FindAllByPage(page int, pageSize int) error {
	if err := db.Scopes(Paginate(page, pageSize)).Find(f).Error; err != nil {
		helper.LogError(err.Error())
		return errors.New("任务结果查询失败")
	}
	return nil
}

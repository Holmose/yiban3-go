package models

import (
	"Yiban3/FrontWeb/helper"
	"encoding/json"
	"errors"
	"strings"
	"time"
)

type User struct {
	ID       int       `gorm:"cloumn:id;primary_key" json:"id"`
	UserName string    `gorm:"column:username;unique;not null" json:"username"`          // 用户名
	Name     string    `gorm:"column:name" json:"name"`                                  // 备注
	PassWord string    `gorm:"column:password;not null" json:"password"`                 // 用户密码
	Email    string    `gorm:"column:e_mail;not null" json:"e_mail"`                     // 邮件
	Verify   string    `gorm:"column:verify" json:"verify"`                              // 认证token
	Period   int       `gorm:"column:day;default:0;not null" json:"day"`                 // 剩余有效期
	Address  string    `gorm:"column:address;default:云南省曲靖市麒麟区;not null" json:"address"` // 位置
	Template int       `gorm:"column:status;default:0;not null" json:"status"`           // 模式切换
	Crontab  string    `gorm:"column:clock_crontab" json:"clock_crontab"`                // 定时任务
	CreateAt time.Time `gorm:"column:create_time;autoCreateTime" json:"create_time"`     // 创建时间
	UpdateAt time.Time `gorm:"column:update_time;autoUpdateTime" json:"update_time"`     // 更新时间
	//Delete   bool      `gorm:"column:delete" json:"delete"`                              // 删除标记
}

type Users []User

// TableName 自定义表名
func (u *User) TableName() string {
	return "yiban_yiban"
}

// CreateUser 创建
func (u *User) Create() error {
	if err := db.First(u, "username=?", u.UserName).Error; err == nil {
		helper.LogError(errors.New(u.UserName + " user already exists").Error())
		return errors.New("user already exists")
	}
	// 创建
	u.CreateAt = time.Now()
	u.UpdateAt = time.Now()
	// TODO 需要优化字段约束
	if strings.TrimSpace(u.UserName) == "" ||
		strings.TrimSpace(u.PassWord) == "" ||
		strings.TrimSpace(u.Address) == "" ||
		strings.TrimSpace(u.Email) == "" {
		helper.LogError(errors.New("the user information is incomplete").Error())
		return errors.New("the user information is incomplete")
	}
	if err := db.Create(u).Error; err != nil {
		helper.LogError(err.Error())
		return errors.New("user creation failed")
	}
	return nil
}

// UpdateUser 更新
func (u *User) Update() error {
	var cache User
	cache = *u
	if err := db.First(u, "username=?", u.UserName).Error; err != nil {
		helper.LogError(err.Error())
		return errors.New("用户不存在")
	}
	// 更新
	cache.ID = u.ID
	bytes, _ := json.Marshal(cache)
	m := make(map[string]interface{})
	json.Unmarshal(bytes, &m)
	if err := db.Model(&u).Omit(
		"create_time", "update_time").Updates(m).Error; err != nil {
		helper.LogError(err.Error())
		return errors.New("用户更新失败")
	}
	cache.UpdateAt = time.Now()
	db.Save(u)
	return nil
}

// FindUserByUserName 查询
func (u *User) FindUserByUserName(username string) error {
	if err := db.First(u, "username=?", username).Error; err != nil {
		helper.LogError(err.Error())
		return errors.New("用户查询失败")
	}
	return nil
}

// FindUsersByUserName 查询
func (u *Users) FindUsersByUserName(username string, opts ...Option) error {
	// 初始化默认值
	query := PageQuery{
		Page: 1,
		Size: 10,
	}

	// 依次调用opts函数列表中的函数，为结构体成员赋值
	for _, o := range opts {
		o(&query)
	}
	if err := db.Scopes(Paginate(query.Page, query.Size)).
		Find(u, "username=?", username).Error; err != nil {
		helper.LogError(err.Error())
		return errors.New("用户查询失败")
	}
	return nil
}

// FindAllByPage 根据分页查询用户信息
// page 页数 pageSize 每页显示的数量
func (u *Users) FindAllByPage(opts ...Option) error {
	// 初始化默认值
	query := PageQuery{
		Page: 1,
		Size: 10,
	}

	// 依次调用opts函数列表中的函数，为结构体成员赋值
	for _, o := range opts {
		o(&query)
	}
	if err := db.Scopes(Paginate(query.Page, query.Size)).Find(u).Error; err != nil {
		helper.LogError(err.Error())
		return errors.New("用户查询失败")
	}
	return nil
}

// DeleteUser 删除
func (u *User) Delete() error {
	// 查询
	if err := u.FindUserByUserName(u.UserName); err != nil {
		return errors.New("用户不存在")
	}
	// 真删除
	if err := db.Delete(u).Error; err != nil {
		helper.LogError(err.Error())
		return errors.New("用户删除失败")
	}
	// 假删除
	//u.Delete = true
	//db.Save(u)
	return nil
}

// CheckUser 查找用户是否存在
func CheckUser(username string) (User, error) {
	var user User
	if err := db.First(&user, "username=?", username).Error; err != nil {
		return User{}, err
	}
	return user, nil
}

package models

import (
	"Yiban3/FrontWeb/helper"
	"encoding/json"
	"errors"
	"strconv"
)

type Admin struct {
	ID       int    `gorm:"column:id; primary_key" json:"id"`
	Name     string `gorm:"column:name;unique;not null" json:"name"`  // 用户名
	PassWord string `gorm:"column:password;not null" json:"password"` // 密码
	Token    string `gorm:"column:token" json:"token"`                // TOKEN
}

// TableName 自定义表名
func (a *Admin) TableName() string {
	return "admin_user"
}

func (a *Admin) CountUserByName(name string) int {
	var count int64
	db.Model(a).Where("name = ?", name).Count(&count)
	ret, _ := strconv.Atoi(strconv.FormatInt(count, 10))
	return ret
}

func (a *Admin) FindUserByName() error {
	if err := db.First(a, "name = ?", a.Name).Error; err != nil {
		helper.LogError(err.Error())
		return errors.New("查询失败")
	}
	return nil
}

func (a *Admin) FindUserByToken() error {
	if err := db.First(a, "token=?", a.Token).Error; err != nil {
		helper.LogError(err.Error())
		return errors.New("查询失败")
	}
	return nil
}

func (a *Admin) Add() error {
	if err := db.Create(a).Error; err != nil {
		helper.LogError(err.Error())
		return errors.New("添加用户失败")
	}
	return nil
}

func (a *Admin) Update() error {

	var cache Admin
	cache = *a
	if err := db.First(a, "name=?", a.Name).Error; err != nil {
		helper.LogError(err.Error())
		return errors.New("用户不存在")
	}
	// 更新
	cache.ID = a.ID
	bytes, _ := json.Marshal(cache)
	m := make(map[string]interface{})
	json.Unmarshal(bytes, &m)
	if err := db.Model(&a).Omit(
		"create_time", "update_time").Updates(m).Error; err != nil {
		helper.LogError(err.Error())
		return errors.New("用户更新失败")
	}
	return nil
}
func (a *Admin) Delete() error {
	// 查询
	if err := a.FindUserByName(); err != nil {
		return errors.New("用户不存在")
	}
	// 真删除
	if err := db.Delete(a).Error; err != nil {
		helper.LogError(err.Error())
		return errors.New("删除用户失败")
	}
	// 假删除
	//u.Delete = true
	//db.Save(u)
	return nil
}

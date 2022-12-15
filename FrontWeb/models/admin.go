package models

import (
	"Yiban3/FrontWeb/helper"
	"errors"
	"strconv"
)

type Admin struct {
	ID       int    `gorm:"column:id; primary_key" json:"id"`
	Name     string `gorm:"column:name" json:"name"`         // 用户名
	PassWord string `gorm:"column:password" json:"password"` // 密码
	Token    string `gorm:"column:token" json:"token"`       // 密码
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
	if err := db.Save(a).Error; err != nil {
		helper.LogError(err.Error())
		return errors.New("修改用户失败")
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

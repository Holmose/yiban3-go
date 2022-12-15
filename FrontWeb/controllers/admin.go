package controllers

import (
	"Yiban3/FrontWeb/helper"
	"Yiban3/FrontWeb/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AdminCreateRequest struct {
	Name     string `json:"name"`     // 用户名
	PassWord string `json:"password"` // 密码
}

type LoginRequest struct {
	Name     string `json:"name"`     // 用户名
	PassWord string `json:"password"` // 密码
}

type AdminController struct {
}

func (a *AdminController) Login(c *gin.Context) {
	var loginRequest LoginRequest
	// 将前端穿过来的json数据绑定存储在这个实体类中，BindJSON()也能使用
	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		c.JSON(400, Response{Message: "数据有误", Error: "data error"})
		helper.LogError(err.Error())
		return
	}
	var admin models.Admin
	admin.Name = loginRequest.Name
	if err := admin.FindUserByName(); err != nil {
		c.JSON(400, err)
		helper.LogError(err.Error())
		return
	}
	if admin.PassWord != helper.Md5(loginRequest.PassWord) {
		c.JSON(400, Response{Message: "密码错误", Error: "password error"})
		return
	}
	admin.Token = helper.Md5(helper.Rand())
	if err := admin.Update(); err != nil {
		c.JSON(400, Response{Message: "认证错误", Error: "token error"})
		return
	}
	c.JSON(200, ResponseLogin{
		Response: Response{Message: "login success"},
		Token:    admin.Token,
	})
}

func (a *AdminController) Create(c *gin.Context) {
	var adminCreateRequest AdminCreateRequest

	// 将前端穿过来的json数据绑定存储在这个实体类中，BindJSON()也能使用
	if err := c.ShouldBindJSON(&adminCreateRequest); err != nil {
		helper.LogError(err.Error())
		return
	}
	var admin models.Admin
	if num := admin.CountUserByName(adminCreateRequest.Name); num > 0 {
		c.JSON(400, Response{Message: "用户已经存在", Error: "user already exists"})
		return
	}

	// 调用业务层的方法
	admin.Name = adminCreateRequest.Name
	admin.PassWord = helper.Md5(adminCreateRequest.PassWord)
	if err := admin.Add(); err != nil {
		c.JSON(400, err)
		helper.LogError(err.Error())
		return
	}
	c.JSON(200, Response{Message: "success"})
}

func (a *AdminController) Update(c *gin.Context) {
	var adminCreateRequest AdminCreateRequest
	// 将前端穿过来的json数据绑定存储在这个实体类中，BindJSON()也能使用
	if err := c.BindJSON(&adminCreateRequest); err != nil {
		helper.LogError(err.Error())
		return
	}

	var admin models.Admin
	admin.Name = adminCreateRequest.Name
	admin.PassWord = helper.Md5(adminCreateRequest.PassWord)
	if err := admin.Update(); err != nil {
		c.JSON(http.StatusBadRequest, Response{Message: "failed", Error: err.Error()})
	} else {
		c.JSON(http.StatusOK, Response{Message: "success"})
	}
}

func (a *AdminController) Delete(c *gin.Context) {
	var admin models.Admin
	c.BindJSON(&admin)
	if err := admin.Delete(); err != nil {
		c.JSON(http.StatusBadRequest, Response{Message: "failed", Error: err.Error()})
	} else {
		c.JSON(http.StatusOK, Response{Message: "success"})
	}
}

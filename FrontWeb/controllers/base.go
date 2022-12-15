package controllers

import (
	"Yiban3/FrontWeb/models"
	"errors"
	"github.com/gin-gonic/gin"
)

var (
	// UserControl 所有的controller类声明都在这儿
	UserControl  = &UserController{}
	AdminControl = &AdminController{}
)

type Response struct {
	Message string `json:"message"`
	Error   string `json:"error"`
}
type ResponseLogin struct {
	Response
	Token string
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.RequestURI == "/login" ||
			c.Request.RequestURI == "/createAdmin" {
			return
		}
		if _, err := Auth(c); err != nil {
			c.Abort()
			c.JSON(400,
				Response{Message: "403 Forbidden", Error: err.Error()})
		}
	}
}

func Auth(c *gin.Context) (*models.Admin, error) {
	token := c.GetHeader("token")
	if token == "" {
		return nil, errors.New("token is null")
	}
	var admin models.Admin
	admin.Token = token
	if err := admin.FindUserByToken(); err != nil {
		return nil, errors.New("insufficient permissions")
	}
	return &admin, nil
}

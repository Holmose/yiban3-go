package controllers

import (
	"Yiban3/FrontWeb/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserController struct{}

func (uc UserController) Create(c *gin.Context) {
	var user models.User
	c.BindJSON(&user)
	if err := user.Create(); err != nil {
		c.JSON(http.StatusBadRequest,
			Response{Message: "user creation failed", Error: err.Error()})
	} else {
		c.JSON(http.StatusOK, Response{Message: "success"})
	}
}
func (uc UserController) Update(c *gin.Context) {
	var user models.User
	c.BindJSON(&user)
	if err := user.Update(); err != nil {
		c.JSON(http.StatusBadRequest,
			Response{Message: "failed", Error: err.Error()})
	} else {
		c.JSON(http.StatusOK, Response{Message: "success"})
	}
}

func (uc UserController) Delete(c *gin.Context) {
	var user models.User
	c.BindJSON(&user)
	if err := user.Delete(); err != nil {
		c.JSON(http.StatusBadRequest,
			Response{Message: "user delete failed", Error: err.Error()})
	} else {
		c.JSON(http.StatusOK, Response{Message: "success"})
	}
}

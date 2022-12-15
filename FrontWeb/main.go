package main

import (
	"Yiban3/FrontWeb/controllers"
	"Yiban3/FrontWeb/models"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	r := gin.Default()
	r.Use(gin.Recovery())
	dsn := "root:123@tcp(192.168.160.142:3306)/django_mysql?charset=utf8&parseTime=True&loc=Local"
	db, sqlDB := models.InitDB(dsn)
	defer sqlDB.Close()
	// 自动迁移模型
	db.AutoMigrate(
		&models.Admin{},
		&models.Users{},
		&models.Task{},
		&models.Form{})

	// 用户鉴权
	r.Use(controllers.AuthMiddleware())

	// 创建测试

	//for i := 0; i < 150; i++ {
	//	user2 := model.User{
	//		UserName: "1103",
	//		Name:     "110",
	//		Email:    "110",
	//		Verify:   "",
	//		Period:   0,
	//		Address:  "11",
	//		Template: 0,
	//		Crontab:  "",
	//	}
	//	user2.UserName = strconv.Itoa(i)
	//	err := user2.CreateUpdateUser()
	//	fmt.Println(err)
	//}

	//var form model.Users
	//form.FindAllByPage(model.SetPage(5))
	//fmt.Println(form)

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	// 账户管理
	r.POST("/create", controllers.UserControl.Create)
	r.POST("/update", controllers.UserControl.Update)
	r.POST("/delete", controllers.UserControl.Delete)

	// 管理员
	r.POST("/login", controllers.AdminControl.Login)
	r.POST("/creatAdmin", controllers.AdminControl.Create)
	r.POST("/updateAdmin", controllers.AdminControl.Update)
	r.POST("/deleteAdmin", controllers.AdminControl.Delete)

	err := r.Run()
	if err != nil {
		return
	} // listen and serve on 0.0.0.0:8080
}

package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type Todo struct {
	ID     int    `json:"id"`
	Tittle string `json:"title"`
	Status bool   `json:"status"`
}

var (
	DB *gorm.DB
)

func initMySQL() (err error) {
	dsn := "root:weidongsheng1212@tcp(127.0.0.1:3306)/db1?charset=utf8mb4&parseTime=True&loc=Local"
	DB, err = gorm.Open("mysql", dsn)
	if err != nil {
		return
	}
	return DB.DB().Ping()
}

func main() {
	//创建数据库 sql:create database bubble; use bubble;\

	//连接数据库
	err := initMySQL()
	if err != nil {
		panic(err)
	}
	defer DB.Close() //程序退出时关闭数据库连接
	//模型绑定
	DB.AutoMigrate(&Todo{})

	r := gin.Default()
	r.Static("/static", "static")
	r.LoadHTMLGlob("templates/*")
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	v1Group := r.Group("v1")
	{
		//待办事项
		//添加
		v1Group.POST("/todo", func(c *gin.Context) {
			//前端页面填写待办事项 点击提交 发送请求到这里
			//1.从请求中把数据拿出来
			var todo Todo
			c.BindJSON(&todo)
			//2.存入数据库
			err = DB.Create(&todo).Error
			//3.返回响应
			if err != nil {
				c.JSON(http.StatusOK, gin.H{"error": err.Error()})
			} else {
				c.JSON(http.StatusOK, todo)
			}
		})

		//查看待办事项
		v1Group.GET("/todo", func(c *gin.Context) {
			var todolist []Todo
			if err = DB.Find(&todolist).Error; err != nil {
				c.JSON(http.StatusOK, gin.H{"error": err.Error()})
			} else {
				c.JSON(http.StatusOK, todolist)
			}
		})

		//修改某一个待办事项
		v1Group.PUT("/todo/:id", func(c *gin.Context) {
			id, ok := c.Params.Get("id")
			if !ok {
				c.JSON(http.StatusOK, gin.H{"error": "id不存在"})
				return
			}
			var todo Todo
			if err = DB.Where("id=?", id).First(&todo).Error; err != nil {
				c.JSON(http.StatusOK, gin.H{"error": err.Error()})
				return
			}
			c.BindJSON(&todo)
			if err = DB.Save(&todo).Error; err != nil {
				c.JSON(http.StatusOK, gin.H{"error": err.Error()})
			} else {
				c.JSON(http.StatusOK, todo)
			}
		})

		//删除待办事项
		v1Group.DELETE("/todo/:id", func(c *gin.Context) {
			id, ok := c.Params.Get("id")
			if !ok {
				c.JSON(http.StatusOK, gin.H{"error": "无效的id"})
				return
			}
			if err = DB.Where("id=?", id).Delete(Todo{}).Error; err != nil {
				c.JSON(http.StatusOK, gin.H{"error": err.Error()})
			} else {
				c.JSON(http.StatusOK, gin.H{id: "deleted"})
			}
		})
	}

	r.Run(":9090")
}

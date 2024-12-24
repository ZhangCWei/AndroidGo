package main

import (
	"github.com/ZhangCWei/AndroidGo/controllers"
	"github.com/ZhangCWei/AndroidGo/minio"
	"github.com/ZhangCWei/AndroidGo/models"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
)

func main() {
	// 设置本机IP
	URL := "192.168.3.200"

	// 连接到SQLite数据库
	db, err := gorm.Open(sqlite.Open("./database/android.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// 自动迁移数据库
	err = db.AutoMigrate(&models.Attention{}, &models.User{}, &models.QuestionBox{})
	if err != nil {
		return
	}

	// 连接Minio
	endpoint := URL + ":9000"
	accessKeyID := "JEbWQujHW7ET6HefNn8f"
	secretAccessKey := "EGJoHjmJvj0RAOVgItp8p8UExJJiKW7ozYjYaU82"
	bucketName := "questionbox"

	// 创建MinioClient实例
	minioClient, err := minio.NewMClient(endpoint, accessKeyID, secretAccessKey, bucketName, false)

	// 创建Gin引擎
	r := gin.Default()

	// 创建控制器
	userController := controllers.NewUserController(db, minioClient)
	attentionController := controllers.NewAttentionController(db, minioClient)
	questionBoxController := controllers.NewQuestionBoxController(db)

	// 注册user路由和处理器函数
	r.GET("/login", userController.Login)
	r.GET("/register/check", userController.RegisterCheck)
	r.POST("/register/confirm", userController.RegisterConfirm)
	r.POST("/changeName", userController.ChangeName)
	r.POST("/changePassword", userController.ChangePassword)
	r.POST("/upload", userController.UploadAudio)
	r.GET("/getHeader", userController.GetHeader)

	// 注册attention路由和处理器函数
	r.GET("/square/myAttention", attentionController.SquareMyAttention)
	r.GET("/square/myFans", attentionController.SquareMyFans)
	r.POST("/square/add", attentionController.SquareAddAttention)
	r.POST("/square/delete", attentionController.SquareDeleteAttention)

	// 注册question_box路由和处理器函数
	r.GET("/getTarget", questionBoxController.GetTarget)
	r.GET("/getSource", questionBoxController.GetSource)
	r.POST("/deleteItem", questionBoxController.DeleteItem)
	r.POST("/answer", questionBoxController.Answer)
	r.POST("/askQuestion", questionBoxController.AskQuestion)
	r.GET("/getDetail", questionBoxController.GetDetail)

	// 启动HTTP服务器
	err = r.Run(URL + ":8080")
}

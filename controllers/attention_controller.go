package controllers

import (
	"encoding/base64"
	"github.com/ZhangCWei/AndroidGo/minio"
	"github.com/ZhangCWei/AndroidGo/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"net/http"
)

type AttentionController struct {
	DB          *gorm.DB
	MinioClient *minio.MClient
}

type ListOfTarget struct {
	TargetName string `json:"TargetName"`
	Target     string `json:"Target"`
	ImageBytes string `json:"imageBytes"`
}

// NewAttentionController 创建 AttentionController 的实例
func NewAttentionController(db *gorm.DB, minioClient *minio.MClient) *AttentionController {
	return &AttentionController{
		DB:          db,
		MinioClient: minioClient,
	}
}

// SquareMyAttention 处理 /square/myAttention 请求
func (ctrl *AttentionController) SquareMyAttention(c *gin.Context) {
	source := c.Query("myattention")

	var attentions []models.Attention
	ctrl.DB.Where("source_phone = ?", source).Find(&attentions)

	var targetList []ListOfTarget
	for _, attention := range attentions {
		objectName := attention.TargetPhone + ".png"
		imageBytes, err := ctrl.MinioClient.GetImage(objectName)
		if err != nil {
			log.Printf("Failed to get %v and return default image.", attention.TargetPhone+".png")
			log.Printf("Caused by %v", err)
			imageBytes, _ = ctrl.MinioClient.GetImage("0.png")
		}
		imageBase64 := base64.StdEncoding.EncodeToString(imageBytes)
		targetList = append(targetList, ListOfTarget{
			TargetName: attention.TargetName,
			Target:     attention.TargetPhone,
			ImageBytes: imageBase64,
		})
	}

	c.JSON(http.StatusOK, targetList)
}

// SquareMyFans 处理 /square/myfans 请求
func (ctrl *AttentionController) SquareMyFans(c *gin.Context) {
	target := c.Query("myfans")

	var fans []models.Attention
	ctrl.DB.Where("target_phone = ?", target).Find(&fans)

	var sourceList []ListOfTarget
	for _, fan := range fans {
		objectName := fan.SourcePhone + ".png"
		imageBytes, err := ctrl.MinioClient.GetImage(objectName)
		if err != nil {
			log.Printf("Failed to get %v and return default image.", fan.SourcePhone+".png")
			log.Printf("Caused by %v", err)
			imageBytes, _ = ctrl.MinioClient.GetImage("0.png")
		}
		imageBase64 := base64.StdEncoding.EncodeToString(imageBytes)
		sourceList = append(sourceList, ListOfTarget{
			TargetName: fan.SourceName,
			Target:     fan.SourcePhone,
			ImageBytes: imageBase64,
		})
	}

	c.JSON(http.StatusOK, sourceList)
}

// SquareAddAttention 处理 /square/add 请求
func (ctrl *AttentionController) SquareAddAttention(c *gin.Context) {
	source := c.PostForm("source")
	sourceName := c.PostForm("sourceName")
	target := c.PostForm("target")

	var attentions []models.Attention
	ctrl.DB.Where("source_phone = ?", source).Find(&attentions)

	for _, attention := range attentions {
		if attention.TargetPhone == target {
			c.String(http.StatusOK, "repeated")
			return
		}
	}

	var targetUser models.User
	if err := ctrl.DB.Where("phone = ?", target).First(&targetUser).Error; err != nil {
		c.String(http.StatusBadGateway, "user not found")
		return
	}

	newAttention := models.Attention{
		SourcePhone: source,
		SourceName:  sourceName,
		TargetPhone: target,
		TargetName:  targetUser.UserName,
	}
	ctrl.DB.Create(&newAttention)
	c.String(http.StatusOK, "successful")
}

// SquareDeleteAttention 处理 /square/delete 请求
func (ctrl *AttentionController) SquareDeleteAttention(c *gin.Context) {
	source := c.PostForm("source")
	target := c.PostForm("target")

	var attentions []models.Attention
	ctrl.DB.Where("target_phone = ? AND source_phone = ?", target, source).Delete(&attentions)

	c.JSON(http.StatusOK, "deleted")
}

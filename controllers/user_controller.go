package controllers

import (
	"bytes"
	"fmt"
	"github.com/ZhangCWei/AndroidGo/minio"
	"github.com/ZhangCWei/AndroidGo/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"image"
	"image/png"
	"log"
	"net/http"
)

type UserController struct {
	DB          *gorm.DB
	MinioClient *minio.MClient
}

func NewUserController(db *gorm.DB, minioClient *minio.MClient) *UserController {
	return &UserController{
		DB:          db,
		MinioClient: minioClient,
	}
}

func (ctrl *UserController) Login(c *gin.Context) {
	phone := c.Query("phone")

	var user models.User
	if err := ctrl.DB.Where("phone = ?", phone).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (ctrl *UserController) RegisterCheck(c *gin.Context) {
	phoneNumber := c.Query("phonenumber")

	var user models.User
	if err := ctrl.DB.Where("phone = ?", phoneNumber).First(&user).Error; err == nil {
		c.String(http.StatusOK, "registered")
	} else {
		c.String(http.StatusOK, "notRegistered")
	}
}

func (ctrl *UserController) RegisterConfirm(c *gin.Context) {
	var user models.User
	if err := c.Bind(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	user.IsChanged = 0
	if err := ctrl.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		return
	}

	c.String(http.StatusOK, "saved")
}

func (ctrl *UserController) ChangeName(c *gin.Context) {
	phone := c.PostForm("phone")
	name := c.PostForm("name")

	var user models.User
	if err := ctrl.DB.Where("phone = ?", phone).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	user.UserName = name
	ctrl.DB.Save(&user)
	c.Status(http.StatusOK)
}

func (ctrl *UserController) ChangePassword(c *gin.Context) {
	phone := c.PostForm("phone")
	password := c.PostForm("password")

	var user models.User
	if err := ctrl.DB.Where("phone = ?", phone).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	user.Password = password
	ctrl.DB.Save(&user)
	c.Status(http.StatusOK)
}

func (ctrl *UserController) ChangePhone(c *gin.Context) {
	phone := c.PostForm("phone")
	newPhone := c.PostForm("new_phone")

	var user models.User
	if err := ctrl.DB.Where("phone = ?", phone).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	user.Phone = newPhone
	ctrl.DB.Save(&user)
	c.Status(http.StatusOK)
}

func (ctrl *UserController) UploadAudio(c *gin.Context) {
	phone := c.Request.Header.Get("phone")
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}

	objectName := fmt.Sprintf("%s.png", phone)
	err = ctrl.MinioClient.UploadImage(objectName, file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload file"})
		return
	}

	var user models.User
	if err := ctrl.DB.Where("phone = ?", phone).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	user.IsChanged = 1
	ctrl.DB.Save(&user)
	c.Status(http.StatusOK)
}

func (ctrl *UserController) GetHeader(c *gin.Context) {
	phone := c.Query("phone")
	objectName := phone + ".png"

	imageBytes, err := ctrl.MinioClient.GetImage(objectName)
	if err != nil {
		log.Printf("Failed to get %v and return default image.", phone+".png")
		log.Printf("Caused by %v", err)
		imageBytes, _ = ctrl.MinioClient.GetImage("0.png")
	}

	img, _, err := image.Decode(bytes.NewReader(imageBytes))
	if err != nil {
		log.Printf("Failed to decode image. Caused by %v", err)
		imageBytes, _ = ctrl.MinioClient.GetImage("0.png")
		img, _, _ = image.Decode(bytes.NewReader(imageBytes))
	}

	buf := new(bytes.Buffer)
	if err := png.Encode(buf, img); err != nil {
		log.Printf("Failed to encode image: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to encode image",
			"error":   err.Error(),
		})
		return
	}

	c.Data(http.StatusOK, "image/png", buf.Bytes())
}

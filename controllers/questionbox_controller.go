package controllers

import (
	"github.com/ZhangCWei/AndroidGo/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

type QuestionBoxController struct {
	DB *gorm.DB
}

func NewQuestionBoxController(db *gorm.DB) *QuestionBoxController {
	return &QuestionBoxController{DB: db}
}

func (ctrl *QuestionBoxController) GetTarget(c *gin.Context) {
	targetPhone := c.Query("phone")
	state := c.Query("state")

	var boxes []models.QuestionBox
	ctrl.DB.Where("target_phone = ? AND state = ?", targetPhone, state).Find(&boxes)

	c.JSON(http.StatusOK, boxes)
}

func (ctrl *QuestionBoxController) GetSource(c *gin.Context) {
	sourcePhone := c.Query("phone")
	state := c.Query("state")

	var boxes []models.QuestionBox
	ctrl.DB.Where("source_phone = ? AND state = ?", sourcePhone, state).Find(&boxes)

	c.JSON(http.StatusOK, boxes)
}

func (ctrl *QuestionBoxController) DeleteItem(c *gin.Context) {
	idStr := c.PostForm("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	ctrl.DB.Delete(&models.QuestionBox{}, id)
	c.JSON(http.StatusOK, gin.H{"message": "delete!"})
}

func (ctrl *QuestionBoxController) Answer(c *gin.Context) {
	idStr := c.PostForm("id")
	answer := c.PostForm("answer")
	answerTime := c.PostForm("answertime")
	state := "1"

	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var qBox models.QuestionBox
	if err := ctrl.DB.First(&qBox, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "QuestionBox not found"})
		return
	}

	qBox.Answer = answer
	qBox.AnswerTime = answerTime
	qBox.State = state

	ctrl.DB.Save(&qBox)
	c.JSON(http.StatusOK, qBox)
}

func (ctrl *QuestionBoxController) AskQuestion(c *gin.Context) {
	source := c.PostForm("source")
	target := c.PostForm("target")
	question := c.PostForm("question")
	questionTime := c.PostForm("questiontime")
	targetName := c.PostForm("targetName")
	state := "0"

	// 赋值
	newQuestion := models.QuestionBox{
		SourcePhone:  source,
		TargetPhone:  target,
		Question:     question,
		QuestionTime: questionTime,
		State:        state,
		TargetName:   targetName,
	}

	// 保存到数据库
	if err := ctrl.DB.Create(&newQuestion).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

func (ctrl *QuestionBoxController) GetDetail(c *gin.Context) {
	idStr := c.Query("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var qBox models.QuestionBox
	if err := ctrl.DB.First(&qBox, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "QuestionBox not found"})
		return
	}

	c.JSON(http.StatusOK, qBox)
}

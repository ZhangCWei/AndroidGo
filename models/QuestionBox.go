package models

type QuestionBox struct {
	ID           uint   `gorm:"primaryKey" json:"Id"`
	Answer       string `gorm:"column:answer" json:"answer"`
	AnswerTime   string `gorm:"column:answer_time" json:"answerTime"`
	Question     string `gorm:"column:question;not null" json:"question"`
	QuestionTime string `gorm:"column:question_time;not null" json:"questionTime"`
	SourcePhone  string `gorm:"column:source_phone;not null" json:"sourcePhone"`
	State        string `gorm:"column:state;not null" json:"state"`
	TargetName   string `gorm:"column:target_name;not null" json:"targetName"`
	TargetPhone  string `gorm:"column:target_phone;not null" json:"targetPhone"`
}

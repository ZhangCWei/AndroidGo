package models

type User struct {
	ID        uint   `gorm:"primaryKey" json:"Id"`
	IsChanged int    `gorm:"column:is_changed;not null" json:"isChanged"`
	UserName  string `gorm:"column:username;not null" json:"username"`
	Phone     string `gorm:"column:phone;not null" json:"phone"`
	Password  string `gorm:"column:password;not null" json:"password"`
}

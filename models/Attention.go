package models

type Attention struct {
	ID          uint   `gorm:"primaryKey" json:"Id"`
	TargetPhone string `gorm:"column:target_phone;not null" json:"target"`
	SourcePhone string `gorm:"column:source_phone;not null" json:"source"`
	TargetName  string `gorm:"column:target_name;not null" json:"targetName"`
	SourceName  string `gorm:"column:source_name;not null" json:"sourceName"`
}
